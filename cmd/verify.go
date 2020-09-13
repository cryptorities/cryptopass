package cmd

import (
	"fmt"
	"github.com/cryptorities/cryptopass/pkg/crypto"
	"github.com/cryptorities/cryptopass/pkg/util"
	"github.com/pkg/errors"
)

/**
	Alex Shvid
*/

type verifyCommand struct {
}

func (t *verifyCommand) Desc() string {
	return "verify token with expiration date"
}

func (t *verifyCommand) Run(args []string) error {

	if len(args) < 2 {
		return errors.Errorf("expected two or more arguments: %v", args)
	}

	username := args[0]
	token := args[1]

	valid, expiration, err := crypto.VerifyIssued(username, token, util.PromptPublicKey)
	if err != nil {
		return err
	}

	if valid {
		fmt.Printf("valid till %s\n", expiration)
		return nil
	}

	valid, expiration, err = crypto.VerifyRevoked(username, token, util.PromptPublicKey)
	if err != nil {
		return err
	}

	if valid {
		fmt.Printf("revoked at %s\n", expiration)
		return nil
	}

	println("invalid")
	return nil

}
