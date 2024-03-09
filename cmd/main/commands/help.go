package commands

import (
	"github.com/spf13/cobra"
	"os"
)

const flagHelp = "help"

func addHelpFlag(cmd *cobra.Command) {
	cmd.Flags().Bool(flagHelp, false, "Print this help message and exit")
}

func handleHelpFlag(cmd *cobra.Command) {
	if cmd.Flag(flagHelp).Value.String() == "true" {
		_ = cmd.Help()
		os.Exit(0)
	}
}
