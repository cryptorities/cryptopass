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

type encCommand struct {
}

func (t *encCommand) Desc() string {
	return "encrypt file"
}

func (t *encCommand) Usage() string {
	return "enc input_file [output_file]"
}

func (t *encCommand) Run(args []string) error {

	if len(args) < 1 {
		return errors.New("Usage: ./cryptopass enc input_file [output_file]")
	}

	inputFile := args[0]
	args = args[1:]

	var outputFile string
	if len(args) > 0 {
		outputFile = args[0]
	} else {
		outputFile = fmt.Sprintf("%s.cp", inputFile)
	}

	n, err := crypto.EncryptFile(inputFile, outputFile, util.PromptRecipientPublicKey)
	if err != nil {
		return err
	}

	fmt.Printf("Written %d bytes.", n)
	return nil
}
