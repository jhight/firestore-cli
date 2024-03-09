package actions

import (
	"fmt"
	"github.com/spf13/cobra"
	"slices"
)

func Delete(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:     "delete <collection> <document>",
		Short:   "Delete a document by ID",
		Long:    "Delete a Firestore document from a collection by its ID.",
		Example: `firestore-cli delete users 1234`,
		Args:    cobra.ExactArgs(2),
		PreRunE: a.initializer.Initialize,
		RunE:    a.runDelete,
	}

	a.addHelpFlag()

	return a
}

func (a *action) runDelete(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	collection := args[0]
	documentID := args[1]

	if slices.Contains(a.initializer.Config().Backup.Commands, "delete") {
		before, _ := a.initializer.Firestore().Get(collection, documentID)
		a.backup(collection, documentID, before, nil)
	}

	err := a.initializer.Firestore().Delete(collection, documentID)
	if err != nil {
		return err
	}

	fmt.Printf("%s/%s successfully deleted\n", collection, documentID)
	return nil
}
