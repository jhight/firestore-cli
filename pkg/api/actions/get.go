package actions

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func Get(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:     "get <path> [<field>]",
		Aliases: []string{"g"},
		Short:   "Get a single document, or a field within",
		Long:    "Get either a single Firestore document from a collection or the specified field within the document.",
		Example: strings.ReplaceAll(`%E get users/1234
%E get users/1234 name
%E get users/1234/orders/5678
%E get users/1234 address.city`, "%E", os.Args[0]),
		Args:    cobra.MinimumNArgs(1),
		PreRunE: a.initializer.Initialize,
		RunE:    a.runGet,
	}

	a.addHelpFlag()

	return a
}

func (a *action) runGet(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	path := args[0]
	field := ""

	if len(args) > 1 {
		field = args[1]
	}

	document, err := a.initializer.Firestore().Get(path)
	if err != nil {
		return err
	}

	if len(field) > 0 {
		if strings.Contains(field, ".") {
			fields := strings.Split(field, ".")
			value, ok := a.nestedField(document, fields)
			if !ok {
				return fmt.Errorf("field %s does not exist in %s\n", field, path)
			}
			a.printOutput(value)
			return nil
		}

		value, ok := document[field]
		if !ok {
			return fmt.Errorf("field %s does not exist in %s\n", field, path)
		}
		a.printOutput(value)
		return nil
	}

	a.printOutput(document)
	return nil
}
