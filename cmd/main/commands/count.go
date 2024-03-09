package commands

import (
	"github.com/spf13/cobra"
)

func NewCountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "count <collection>",
		Aliases: []string{"c"},
		Short:   "Returns a count of all documents in a collection",
		Long:    "Returns a count of all documents in a Firestore collection.",
		Example: `firestore-cli count users`,
		Args:    cobra.ExactArgs(1),
		PreRun:  runRootCommand,
		Run: func(cmd *cobra.Command, args []string) {
			runCountCommand(cmd, args)
		},
	}

	addHelpFlag(cmd)

	return cmd
}

func runCountCommand(cmd *cobra.Command, args []string) {
	handleHelpFlag(cmd)
	collection := args[0]
	count := firestore.Count(collection)
	printOutput(map[string]any{"$count": count})
}
