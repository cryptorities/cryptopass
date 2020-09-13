package main

import (
	"github.com/cryptorities/cryptopass/cmd"
	"github.com/cryptorities/cryptopass/pkg/app"
	"math/rand"
	"os"
	"time"
)

/**
	Alex Shvid
 */

var (
	Version   string
	Build     string
)

func main() {

	app.SetAppInfo(Version, Build)

	rand.Seed(time.Now().UnixNano())
	os.Exit(cmd.Run(os.Args[1:]))

}
