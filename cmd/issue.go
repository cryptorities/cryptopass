package cmd

import (
	"github.com/cryptorities/cryptopass/pkg/app"
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
	return "issue token with expiration date"
}

func (t *issueCommand) Run(args []string) error {

	if len(args) < 2 {
		return errors.Errorf("expected two or more arguments: %v", args)
	}

	username := args[0]
	date := args[1]

	token, err := crypto.Sign(username, date, app.IssueSep, util.PromptPrivateKey)
	if err != nil {
		return err
	}

	println(token)

	return nil
}
