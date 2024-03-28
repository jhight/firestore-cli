package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jhight/firestore-cli/pkg/api/client/query"
	"github.com/spf13/cobra"
	"os"
	"slices"
	"strings"
)

func Update(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:     "update <path> [<json>]",
		Aliases: []string{"u"},
		Short:   "Update specific properties in a document",
		Long:    "Update the specified Firestore document with the specified JSON data. Other fields will remain unchanged. If the field does not exist, it will be created. If the specified document does not exist, a new one will not be created.",
		Example: strings.ReplaceAll(`%E update users/1234 '{"name": "John Doe", "age": 30, "height": 5.9, "active": true}'
%E update users/1234/orders/5678 '{"item": "shoes"}'
cat file.json | %E update users 1234`, "%E", os.Args[0]),
		Args:    cobra.MinimumNArgs(1),
		PreRunE: a.initializer.Initialize,
		RunE:    a.runUpdate,
	}

	a.addHelpFlag()

	return a
}

func (a *action) runUpdate(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	path := args[0]

	var input string
	if len(args) >= 2 {
		input = args[1]
	} else if a.shouldReadFromStdin() {
		var err error
		if input, err = a.readFromStdin(); err != nil {
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
		before, _ := a.initializer.Firestore().Get(query.Input{Path: path})
		if err = a.initializer.Firestore().Update(path, fields); err != nil {
			return err
		}
		after, _ := a.initializer.Firestore().Get(query.Input{Path: path})
		a.backup(path, before, after)
	} else {
		if err = a.initializer.Firestore().Update(path, fields); err != nil {
			return err
		}
	}

	fmt.Printf("%s successfully updated\n", path)
	return nil
}
