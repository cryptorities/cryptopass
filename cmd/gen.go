package cmd

import (
	"crypto/rand"
	"fmt"
	"crypto/ed25519"
	"github.com/cryptorities/cryptopass/pkg/app"
)

/**
	Alex Shvid
*/

type genCommand struct {
}

func (t *genCommand) Desc() string {
	return "generate key pair"
}

func (t *genCommand) Run(args []string) error {

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)

	if err != nil {
		return err
	}

	publicKeyPem := app.Encoding.EncodeToString(publicKey)
	privateKeyPem := app.Encoding.EncodeToString(privateKey)

	fmt.Printf("Public Key: %s\n", publicKeyPem)
	fmt.Printf("Private Key: %s\n", privateKeyPem)

	return nil
}
