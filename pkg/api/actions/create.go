package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func Create(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:     "create <path> [<json>]",
		Aliases: []string{"insert"},
		Short:   "Create a document",
		Long:    "Create a Firestore document with the specified ID using the specified field(s). If a document exists with the same ID, it will be replaced.",
		Example: strings.ReplaceAll(`%E create users/1234 '{"name": "John Doe", "age": 30, "height": 5.9, "active": true}'
%E create users/1234/orders/5678 '{"item": "shoes", "quantity": 1, "price": 100.00}'
cat file.json | %E create users 1234`, "%E", os.Args[0]),
		Args:    cobra.MinimumNArgs(1),
		PreRunE: a.initializer.Initialize,
		RunE:    a.runCreate,
	}

	a.addHelpFlag()

	return a
}

func (a *action) runCreate(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	path := args[0]

	var jsonValue string
	if len(args) >= 2 {
		jsonValue = args[1]
	} else if a.shouldReadFromStdin() {
		var err error
		jsonValue, err = a.readFromStdin()
		if err != nil {
			return err
		}
	}

	var u any
	if err := json.Unmarshal([]byte(jsonValue), &u); err != nil {
		return err
	}

	switch u.(type) {
	case map[string]any:
		if err := a.initializer.Firestore().Create(path, u.(map[string]any)); err != nil {
			return err
		}
	default:
		return errors.New("invalid JSON format")
	}

	fmt.Printf("%s successfully created\n", path)
	return nil
}
