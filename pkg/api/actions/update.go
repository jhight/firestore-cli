package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"slices"
)

func Update(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:     "update <collection> <document> [<json>]",
		Aliases: []string{"u"},
		Short:   "Update a document",
		Long:    "Update the specified Firestore document with the specified JSON data. Other fields will remain unchanged. If the field does not exist, it will be created. If the specified document does not exist, a new one will not be created.",
		Example: `firestore-cli update users 1234 '{"name": "John Doe", "age": 30, "height": 5.9, "active": true}'
cat file.json | firestore-cli update users 1234`,
		Args:    cobra.MinimumNArgs(2),
		PreRunE: a.initializer.Initialize,
		RunE:    a.runUpdate,
	}

	a.addHelpFlag()

	return a
}

func (a *action) runUpdate(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	collection := args[0]
	documentID := args[1]

	var input string
	if len(args) == 3 {
		input = args[2]
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
		before, _ := a.initializer.Firestore().Get(collection, documentID)
		err = a.initializer.Firestore().Update(collection, documentID, fields)
		if err != nil {
			return err
		}
		after, _ := a.initializer.Firestore().Get(collection, documentID)
		a.backup(collection, documentID, before, after)
	} else {
		err = a.initializer.Firestore().Update(collection, documentID, fields)
		if err != nil {
			return err
		}
	}

	fmt.Printf("%s/%s successfully updated\n", collection, documentID)
	return nil
}
