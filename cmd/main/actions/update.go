package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"slices"
)

func Update(root RootAction) *Action {
	a := &Action{
		root:      root,
		firestore: root.Firestore(),
		cfg:       root.Config(),
	}

	a.command = &cobra.Command{
		Use:     "update <collection> <document> [<json>]",
		Aliases: []string{"u"},
		Short:   "Update a document",
		Long:    "Update the specified Firestore document with the specified JSON data. Other fields will remain unchanged. If the field does not exist, it will be created. If the specified document does not exist, a new one will not be created.",
		Example: `firestore-cli update users 1234 '{"name": "John Doe", "age": 30, "height": 5.9, "active": true}'
cat file.json | firestore-cli update users 1234`,
		Args:    cobra.MinimumNArgs(2),
		PreRunE: a.Initialize,
		RunE:    a.runUpdate,
	}

	a.addHelpFlag()

	return a
}

func (a *Action) runUpdate(_ *cobra.Command, args []string) error {
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
	if slices.Contains(a.cfg.Backup.Commands, "update") {
		before, _ := a.firestore.Get(collection, documentID)
		err = a.firestore.Update(collection, documentID, fields)
		if err != nil {
			return err
		}
		after, _ := a.firestore.Get(collection, documentID)
		a.backup(collection, documentID, before, after)
	} else {
		err = a.firestore.Update(collection, documentID, fields)
		if err != nil {
			return err
		}
	}

	fmt.Printf("%s/%s successfully updated\n", collection, documentID)
	return nil
}
