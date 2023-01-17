package ecies

import (
	"crypto/ed25519"
	"crypto/sha512"
	"crypto/subtle"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"

	"filippo.io/edwards25519"
)

func ReadCertificateFromFile(path string) (*x509.Certificate, error) {
	crt, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(crt)
	if block == nil {
		return nil, errors.New("failed to parse certificate PEM")
	}

	return x509.ParseCertificate(block.Bytes)
}

func ReadPrivateKeyFromFile(path string) (ed25519.PrivateKey, error) {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("failed to parse private key from PEM")
	}

	privParsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	privKey, ok := privParsed.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("invalid private key type")
	}
	return privKey, nil
}

func PublicEd25519FromCertificate(cert *x509.Certificate) (ed25519.PublicKey, error) {
	if cert.PublicKeyAlgorithm != x509.Ed25519 {
		return nil, errors.New("invalid public key algorithm")
	}

	pub, ok := cert.PublicKey.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("invalid public key algorithm")
	}

	return pub, nil
}

// Functions for using Ed25519 keys in X25519 key exchange

func PublicEd25519ToMontgomery(pub ed25519.PublicKey) ([]byte, error) {
	// Convert Ed25519 PublicKey from Edwards to Montogomery form
	nodePublicEdwards, err := edwards25519.NewIdentityPoint().SetBytes(pub)
	if err != nil {
		return nil, err
	}

	nodePublicMontogomery := nodePublicEdwards.BytesMontgomery()
	if subtle.ConstantTimeCompare(nodePublicMontogomery, zeroCompare) == 1 {
		return nil, errors.New("invalid input public key")
	}
	return nodePublicMontogomery, nil
}

func PrivateEd25519ToMontgomery(priv ed25519.PrivateKey) ([]byte, error) {
	// https://words.filippo.io/using-ed25519-keys-for-encryption/
	// To decrypt, we derive the secret scalar according to the Ed25519 spec,
	// and simply use it as an X25519 private key in Ephemeral-Static Diffie-Hellman.
	// The two peers might end up with different v coordinates, if they were to calculate
	// them, but in X25519 the shared secret is just the u coordinate, so no one will notice.

	tmpKey := sha512.Sum512(priv.Seed())
	secretScalar := tmpKey[:32]

	// Clamp scalar in case our CA signed an unsafe CSR
	// See https://www.jcraige.com/an-explainer-on-ed25519-clamping for in-depth explanation
	secretScalar[0] &= 248
	secretScalar[31] &= 127
	secretScalar[31] |= 64

	return secretScalar, nil
}
