package cmd

import (
	"github.com/cryptorities/cryptopass/pkg/crypto"
	"github.com/cryptorities/cryptopass/pkg/util"
	"github.com/pkg/errors"
)

/**
	Alex Shvid
*/

type revokeCommand struct {
}

func (t *revokeCommand) Desc() string {
	return "revoke token with expiration date in format YYYY-mm-dd"
}

func (t *revokeCommand) Usage() string {
	return "revoke username expirationDate"
}

func (t *revokeCommand) Run(args []string) error {

	if len(args) < 2 {
		return errors.Errorf("expected two or more arguments: %v", args)
	}

	username := args[0]
	expirationDate := args[1]

	token, err := crypto.Revoke(username, expirationDate, util.PromptPrivateKey)
	if err != nil {
		return err
	}

	println(token)

	return nil
}
