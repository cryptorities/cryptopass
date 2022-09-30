package cmd

import (
	"github.com/cryptorities/cryptopass/pkg/crypto"
	"github.com/cryptorities/cryptopass/pkg/util"
	"github.com/pkg/errors"
	"strings"
)

/**
	Alex Shvid
*/

type decCommand struct {
}

func (t *decCommand) Desc() string {
	return "decrypt file"
}

func (t *decCommand) Usage() string {
	return "dec input_file output_file"
}

func (t *decCommand) Run(args []string) error {

	if len(args) < 1 {
		return errors.New("Usage: ./cryptopass dec input_file [output_file]")
	}

	inputFile := args[0]
	args = args[1:]

	var outputFile string
	if len(args) > 0 {
		outputFile = args[0]
	} else if strings.HasSuffix(inputFile, ".cp") {
		outputFile = inputFile[:len(inputFile)-3]
	} else {
		return errors.New("Usage: ./cryptopass dec input_file output_file")
	}

	err := crypto.DecryptFile(inputFile, outputFile, util.PromptPrivateKey)
	if err != nil {
		return err
	}

	println("Done.")
	return nil
}
