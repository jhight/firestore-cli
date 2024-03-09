package actions

import (
	"github.com/spf13/cobra"
)

func Count(root RootAction) *Action {
	a := &Action{
		root:      root,
		cfg:       root.Config(),
		firestore: root.Firestore(),
	}

	a.command = &cobra.Command{
		Use:     "count <collection>",
		Aliases: []string{"c"},
		Short:   "Returns a count of all documents in a collection",
		Long:    "Returns a count of all documents in a Firestore collection.",
		Example: `firestore-cli count users`,
		Args:    cobra.ExactArgs(1),
		PreRunE: a.Initialize,
		RunE:    a.runCount,
	}

	a.addHelpFlag()

	return a
}

func (a *Action) runCount(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()
	collection := args[0]
	count := a.firestore.Count(collection)
	a.printOutput(map[string]any{"$count": count})
	return nil
}
