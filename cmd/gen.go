package cmd

import (
	"fmt"
)

/**
	Alex Shvid
*/

type genCommand struct {
}

func (t *genCommand) Desc() string {
	return "generate key pair"
}

func (t *genCommand) Usage() string {
	return "gen"
}

func (t *genCommand) Run(args []string) error {

	publicKey, privateKey, err := generateKeyPair()
	if err != nil {
		return err
	}

	fmt.Printf("Public Key: %s\n", publicKey)
	fmt.Printf("Private Key: %s\n", privateKey)

	println(`
You can save your credentials in system variables to avoid future prompts, which is optional and not safe.
CRYPTOPASS_PRIVATE_KEY
CRYPTOPASS_PUBLIC_KEY
`)

	return nil
}
