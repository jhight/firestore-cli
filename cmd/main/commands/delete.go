package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewDeleteCommand() *cobra.Command {
	get := &cobra.Command{
		Use:     "delete <collection> <document>",
		Aliases: []string{"d"},
		Short:   "Delete a document by ID",
		Long:    "Delete a Firestore document from a collection by its ID.",
		Example: `firestore-cli delete users 1234`,
		Args:    cobra.ExactArgs(2),
		PreRun:  runRootCommand,
		Run: func(cmd *cobra.Command, args []string) {
			runDeleteCommand(cmd, args)
		},
	}

	addHelpFlag(get)

	return get
}

func runDeleteCommand(cmd *cobra.Command, args []string) {
	handleHelpFlag(cmd)

	collection := args[0]
	documentID := args[1]

	err := firestore.Delete(collection, documentID)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fmt.Printf("%s/%s successfully deleted\n", collection, documentID)
}
