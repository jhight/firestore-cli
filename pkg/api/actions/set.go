package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"slices"
	"strings"
)

func Set(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:     "set <path> [<json>]",
		Aliases: []string{"import"},
		Short:   "Set (e.g., create or replace) a document",
		Long:    "Set the entire specified Firestore document with specified JSON data. Only the specified fields will exist in the document. If the document does not exist, it will be created.",
		Example: strings.ReplaceAll(`%E set users/1234 '{"name": "John Doe", "age": 30, "height": 5.9, "active": true}'
%E set users/1234/orders/5678 '{"item": "shoes", "quantity": 1, "price": 100.00}'
cat file.json | %E set users/1234`, "%E", os.Args[0]),
		Args:    cobra.MinimumNArgs(1),
		PreRunE: a.initializer.Initialize,
		RunE:    a.runSet,
	}

	a.addHelpFlag()

	return a
}

func (a *action) runSet(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	path := args[0]

	var input string
	if len(args) >= 2 {
		input = args[1]
	} else if a.shouldReadFromStdin() {
		var err error
		input, err = a.readFromStdin()
		if err != nil {
			return err
		}
	}

	if len(input) == 0 {
		return errors.New("one or more fields in JSON format are required")
	}

	var fields map[string]any
	err := json.Unmarshal([]byte(input), &fields)
	if err != nil {
		return err
	}

	// backup before update, if configured
	if slices.Contains(a.initializer.Config().Backup.Commands, "update") {
		before, _ := a.initializer.Firestore().Get(path)
		err = a.initializer.Firestore().Set(path, fields)
		if err != nil {
			return err
		}
		after, _ := a.initializer.Firestore().Get(path)
		a.backup(path, before, after)
	} else {
		err = a.initializer.Firestore().Set(path, fields)
		if err != nil {
			return err
		}
	}

	fmt.Printf("%s successfully set\n", path)
	return nil
}
