package ecies

import (
	"crypto/rand"
	"encoding/asn1"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/sha3"
)

func Encrypt(publicKey []byte, msg []byte) ([]byte, error) {
	if len(publicKey) != curve25519.ScalarSize {
		return nil, errors.New("invalid public key")
	}

	ephemeralPrivate := make([]byte, 32)
	_, err := rand.Read(ephemeralPrivate)
	if err != nil {
		return nil, err
	}

	ephemeralPublic, err := curve25519.X25519(ephemeralPrivate, curve25519.Basepoint)
	if err != nil {
		return nil, err
	}

	sharedSecret, err := curve25519.X25519(ephemeralPrivate, publicKey)
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

	// Random nonce
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return asn1.Marshal(ECIEMessage{
		Version:         1,
		EphemeralPublic: ephemeralPublic,
		Nonce:           nonce,
		Ciphertext:      aead.Seal(msg[:0], nonce, msg, nil),
	})
}
