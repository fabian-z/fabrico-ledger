package ecies

import (
	"encoding/asn1"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
)

// Version 1:
// Curve25519 / X25519 style Montgomery Public Key
// ChaCha20-Poly1305 with extended nonce (randomly generated)
// SHA3-256 to derive key from X25519 secret

type ECIEMessage struct {
	Version         int    // ASN1 cannot serialize uint
	EphemeralPublic []byte // Curve25519 / Montgomery
	Nonce           []byte
	Ciphertext      []byte
}

var (
	zeroCompare = make([]byte, 32)
)

func MarshalECIEMessage(msg *ECIEMessage) []byte {
	out, err := asn1.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return out
}

func UnmarshalECIEMEssage(in []byte) (*ECIEMessage, error) {
	msg := &ECIEMessage{}
	rest, err := asn1.Unmarshal(in, msg)
	if err != nil {
		return nil, err
	}
	if len(rest) > 0 {
		return nil, errors.New("unexpected trailing data")
	}
	if msg.Version != 1 || len(msg.EphemeralPublic) != curve25519.PointSize || len(msg.Nonce) != chacha20poly1305.NonceSizeX || len(msg.Ciphertext) == 0 {
		return nil, errors.New("invalid message encoding")
	}

	return msg, nil
}
