package actions

import (
	"os"
)

const flagHelp = "help"

func (a *action) addHelpFlag() {
	a.command.Flags().Bool(flagHelp, false, "Print this help message and exit")
}

func (a *action) handleHelpFlag() {
	if a.command.Flag(flagHelp).Value.String() == "true" {
		_ = a.command.Help()
		os.Exit(0)
	}
}
