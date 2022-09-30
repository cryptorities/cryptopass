package cmd

import (
	"crypto/rand"
	"github.com/cryptorities/cryptopass/pkg/app"
	"crypto/ed25519"
)

func generateKeyPair() (string, string, error) {

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)

	if err != nil {
		return "", "", err
	}

	publicKeyPem := app.Encoding.EncodeToString(publicKey)
	privateKeyPem := app.Encoding.EncodeToString(privateKey)

	return publicKeyPem, privateKeyPem, nil
}
