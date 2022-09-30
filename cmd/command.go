package cmd

import (
	"flag"
	"fmt"
	"github.com/cryptorities/cryptopass/pkg/app"
	"os"
)

/**
	Alex Shvid
*/

type commandFace interface {

	Run(args []string) error

	Desc() string

	Usage() string

}

var allCommands = map[string]commandFace{

	"version": &versionCommand{},

	"gen": &genCommand{},

	"issue": &issueCommand{},

	"revoke": &revokeCommand{},

	"verify": &verifyCommand{},

	"enc": &encCommand{},

	"dec": &decCommand{},

	"help": &helpCommand{},
}

func preprocessArgs(args []string) []string {

	if len(args) == 1 && (args[0] == "-v" || args[0] == "-version" || args[0] == "--version") {
		return []string{"version"}
	}

	return args
}

func printUsage() {

	fmt.Printf("Usage: %s [command]\n", app.ExecutableName)

	for name, command := range allCommands {
		fmt.Printf("    %s - %s\tUsage: ./%s %s\n", name, command.Desc(), app.ExecutableName, command.Usage())
	}

	fmt.Println("Flags:")
	flag.PrintDefaults()

}

func Run(args []string) int {

	args = preprocessArgs(args)

	if len(args) >= 1 {

		cmd := args[0]

		if inst, ok := allCommands[cmd]; ok {

			if err := inst.Run(args[1:]); err != nil {
				fmt.Fprintf(os.Stderr,"Error: %v\n", err)
				return 1
			}
			return 0

		} else {
			fmt.Fprintf(os.Stderr,"Invalid command: %s\n", cmd)
			printUsage()
			return 1
		}

	} else {
		printUsage()
		return 0
	}

	return 0
}
