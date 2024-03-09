package actions

import (
	"os"
)

const flagHelp = "help"

func (a *Action) addHelpFlag() {
	a.command.Flags().Bool(flagHelp, false, "Print this help message and exit")
}

func (a *Action) handleHelpFlag() {
	if a.command.Flag(flagHelp).Value.String() == "true" {
		_ = a.command.Help()
		os.Exit(0)
	}
}
