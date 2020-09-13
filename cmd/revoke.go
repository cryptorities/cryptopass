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

type revokeCommand struct {
}

func (t *revokeCommand) Desc() string {
	return "revoke token with expiration date"
}

func (t *revokeCommand) Run(args []string) error {

	if len(args) < 2 {
		return errors.Errorf("expected two or more arguments: %v", args)
	}

	username := args[0]
	date := args[1]

	token, err := crypto.Sign(username, date, app.RevokeSep, util.PromptPrivateKey)
	if err != nil {
		return err
	}

	println(token)

	return nil
}
