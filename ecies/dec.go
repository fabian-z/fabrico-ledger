package ecies

import (
	"encoding/asn1"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/sha3"
)

func Decrypt(priv []byte, msg []byte) ([]byte, error) {
	if len(priv) != curve25519.ScalarSize {
		return nil, errors.New("invalid scalar size")
	}

	ecieMsg := &ECIEMessage{}
	rest, err := asn1.Unmarshal(msg, ecieMsg)
	if err != nil {
		return nil, err
	}

	if len(rest) > 0 {
		return nil, errors.New("trailing data after parsing message")
	}

	if ecieMsg.Version != 1 {
		return nil, errors.New("unknown encrypted message algorithm version")
	}

	sharedSecret, err := curve25519.X25519(priv, ecieMsg.EphemeralPublic)
	if err != nil {
		return nil, err
	}

	sharedSecretHash := sha3.New256()
	sharedSecretHash.Write(sharedSecret)
	encKey := sharedSecretHash.Sum(nil)

	aead, err := chacha20poly1305.NewX(encKey)
	if err != nil {
		return nil, err
	}

	// Decryption
	if len(ecieMsg.Nonce) < aead.NonceSize() {
		panic("nonce too short")
	}

	// Decrypt the message and check it wasn't tampered with.
	return aead.Open(nil, ecieMsg.Nonce, ecieMsg.Ciphertext, nil)
}
