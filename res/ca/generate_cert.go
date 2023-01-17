// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

var (
	cn         = flag.String("cn", "", "Common Name to generate certificate for")
	host       = flag.String("host", "", "Comma-separated hostnames and IPs to generate a certificate for")
	caCertPath = flag.String("ca-cert", "ca.crt", "CA certificate path to use when generating node key, PEM encoded")
	caKeyPath  = flag.String("ca-key", "ca.key", "CA key path to use when generating node key, PEM encoded")

	validFor   = flag.Duration("duration", 365*24*time.Hour, "Duration that certificate is valid for")
	isCA       = flag.Bool("ca", false, "whether this cert should be its own Certificate Authority")
	ed25519Key = flag.Bool("ed25519", false, "Generate an Ed25519 key")
)

func main() {
	flag.Parse()

	if (len(*host) == 0 && !(*isCA)) || (len(*host) != 0 && (*isCA)) {
		log.Fatalf("Missing or invalid host / CA selection")
	}

	_, priv, err := ed25519.GenerateKey(rand.Reader)

	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// ECDSA, ED25519 and RSA subject keys should have the DigitalSignature
	// KeyUsage bits set in the x509.Certificate template
	// Only RSA subject keys should have the KeyEncipherment KeyUsage bits set. In
	// the context of TLS this KeyUsage is particular to RSA key exchange and
	// authentication.
	keyUsage := x509.KeyUsageDigitalSignature

	notBefore := time.Now()
	notAfter := notBefore.Add(*validFor)

	// TODO check existing serial numbers against generated
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	var cn string = *cn
	if *isCA {
		cn = "Root CA"
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,

		Subject: pkix.Name{
			CommonName:   cn,
			Organization: []string{"3D Chain"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		BasicConstraintsValid: true,
	}

	if ip := net.ParseIP(*host); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, *host)
	}

	var signerCert *x509.Certificate
	var signerKey ed25519.PrivateKey

	if *isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
		signerCert = &template
		signerKey = priv
	} else {
		// Generate P2P Certificates
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}

		// Load CA from ca.pem
		caCert, caKey, err := LoadX509KeyPair("ca.crt", "ca.key")
		if err != nil {
			log.Fatal(err)
		}
		signerCert = caCert
		signerKey = caKey
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, signerCert, priv.Public().(ed25519.PublicKey), signerKey)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	certOut, err := os.Create("cert.pem")
	if err != nil {
		log.Fatalf("Failed to open cert.pem for writing: %v", err)
	}

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("Failed to write data to cert.pem: %v", err)
	}

	if err := certOut.Close(); err != nil {
		log.Fatalf("Error closing cert.pem: %v", err)
	}

	log.Print("wrote cert.pem\n")

	keyOut, err := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Failed to open key.pem for writing: %v", err)
		return
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		log.Fatalf("Failed to write data to key.pem: %v", err)
	}
	if err := keyOut.Close(); err != nil {
		log.Fatalf("Error closing key.pem: %v", err)
	}
	log.Print("wrote key.pem\n")
}

func LoadX509KeyPair(certFile, keyFile string) (*x509.Certificate, ed25519.PrivateKey, error) {
	certPEM, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error loading cert: %w", err)
	}

	keyPEM, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error loading key: %w", err)
	}
	certData, _ := pem.Decode(certPEM)
	keyData, _ := pem.Decode(keyPEM)

	crt, err := x509.ParseCertificate(certData.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing x509 cert: %w", err)
	}

	key, err := x509.ParsePKCS8PrivateKey(keyData.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing pkcs8 private key: %w", err)
	}

	edKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, nil, fmt.Errorf("private key is not ed25519 type")
	}

	return crt, edKey, nil
}
