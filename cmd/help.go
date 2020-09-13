package cmd

/**
	Alex Shvid
*/

type helpCommand struct {

}
func (t *helpCommand) Desc() string {
	return "help command"
}

func (t *helpCommand) Run(args []string) error {
	printUsage()
	return nil
}
