package cmd

import (
	"github.com/cryptorities/cryptopass/pkg/crypto"
	"github.com/cryptorities/cryptopass/pkg/util"
	"github.com/pkg/errors"
)

/**
	Alex Shvid
*/

type issueCommand struct {
}

func (t *issueCommand) Desc() string {
	return "issue token with expiration date in format YYYY-mm-dd"
}

func (t *issueCommand) Usage() string {
	return "issue username expirationDate"
}

func (t *issueCommand) Run(args []string) error {

	if len(args) < 2 {
		return errors.Errorf("expected two or more arguments: %v", args)
	}

	username := args[0]
	expirationDate := args[1]

	token, err := crypto.Issue(username, expirationDate, util.PromptPrivateKey)
	if err != nil {
		return err
	}

	println(token)

	return nil
}
