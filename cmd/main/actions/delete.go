package actions

import (
	"fmt"
	"github.com/spf13/cobra"
	"slices"
)

func Delete(root RootAction) *Action {
	a := &Action{
		root:      root,
		firestore: root.Firestore(),
		cfg:       root.Config(),
	}

	a.command = &cobra.Command{
		Use:     "delete <collection> <document>",
		Short:   "Delete a document by ID",
		Long:    "Delete a Firestore document from a collection by its ID.",
		Example: `firestore-cli delete users 1234`,
		Args:    cobra.ExactArgs(2),
		PreRunE: a.Initialize,
		RunE:    a.runDelete,
	}

	a.addHelpFlag()

	return a
}

func (a *Action) runDelete(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	collection := args[0]
	documentID := args[1]

	if slices.Contains(a.cfg.Backup.Commands, "delete") {
		before, _ := a.firestore.Get(collection, documentID)
		a.backup(collection, documentID, before, nil)
	}

	err := a.firestore.Delete(collection, documentID)
	if err != nil {
		return err
	}

	fmt.Printf("%s/%s successfully deleted\n", collection, documentID)
	return nil
}
