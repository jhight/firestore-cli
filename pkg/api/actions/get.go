package actions

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func Get(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:     "get <collection> <document> [<field>]",
		Aliases: []string{"g", "r"},
		Short:   "Get a document by ID or a field within",
		Long:    "Get either an entire Firestore document from a collection by its ID or the specified field.",
		Example: `firestore-cli get users 1234
		firestore-cli get users 1234 name
		firestore-cli get users 1234 address.city`,
		Args:    cobra.MinimumNArgs(2),
		PreRunE: a.initializer.Initialize,
		RunE:    a.runGet,
	}

	a.addHelpFlag()

	return a
}

func (a *action) runGet(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	collection := args[0]
	documentID := args[1]
	field := ""

	if len(args) > 2 {
		field = args[2]
	}

	document, err := a.initializer.Firestore().Get(collection, documentID)
	if err != nil {
		return err
	}

	if len(field) > 0 {
		if strings.Contains(field, ".") {
			fields := strings.Split(field, ".")
			value, ok := a.nestedField(document, fields)
			if !ok {
				return fmt.Errorf("field %s does not exist in %s/%s\n", field, collection, documentID)
			}
			a.printOutput(value)
			return nil
		}

		value, ok := document[field]
		if !ok {
			return fmt.Errorf("field %s does not exist in %s/%s\n", field, collection, documentID)
		}
		a.printOutput(value)
		return nil
	}

	a.printOutput(document)
	return nil
}
