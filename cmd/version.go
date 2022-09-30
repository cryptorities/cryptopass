package cmd

import (
	"fmt"
	"github.com/cryptorities/cryptopass/pkg/app"
)

/**
	Alex Shvid
*/

type versionCommand struct {
}

func (t *versionCommand) Desc() string {
	return "show version"
}

func (t *versionCommand) Usage() string {
	return "version"
}

func (t *versionCommand) Run(args []string) error {

	appInfo := app.GetAppInfo()
	fmt.Printf("%s [Version %s, Build %s]\n", app.ApplicationName, appInfo.Version, appInfo.Build)
	fmt.Printf("%s\n", app.Copyright)
	return nil
}
