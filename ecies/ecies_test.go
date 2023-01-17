package ecies

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"testing"
)

func TestECIESFromEd25519(t *testing.T) {

	var tests = []int{
		1 << 8,
		1 << 9,
		1 << 10,
		1 << 11,
		1 << 12,
		1 << 13,
	}

	for _, length := range tests {
		testname := fmt.Sprintf("%d", length)
		t.Run(testname, func(t *testing.T) {
			pub, priv, err := ed25519.GenerateKey(rand.Reader)
			if err != nil {
				t.Fatal(err)
			}

			testMsg := make([]byte, length)
			_, err = rand.Read(testMsg)
			if err != nil {
				t.Fatal(err)
			}

			pubMontgomery, err := PublicEd25519ToMontgomery(pub)
			if err != nil {
				t.Fatal(err)
			}

			privMontgomery, err := PrivateEd25519ToMontgomery(priv)
			if err != nil {
				t.Fatal(err)
			}

			encryptedTestMsg, err := Encrypt(pubMontgomery, testMsg)
			if err != nil {
				t.Fatal(err)
			}

			decryptedTestMsg, err := Decrypt(privMontgomery, encryptedTestMsg)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(testMsg, decryptedTestMsg) {
				t.Fatal("Decryption mismatch")
			}
		})
	}

}
