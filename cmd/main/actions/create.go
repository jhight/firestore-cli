package actions

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
)

func Create(root RootAction) *Action {
	a := &Action{
		root:      root,
		firestore: root.Firestore(),
		cfg:       root.Config(),
	}

	a.command = &cobra.Command{
		Use:   "create <collection> <document> [<json>]",
		Short: "Create a document",
		Long:  "Set (replace or create) an entire Firestore document with the specified ID using the specified field(s). If a document exists with the same ID, it will be replaced.",
		Example: `firestore-cli set users 1234 '{"name": "John Doe", "age": 30, "height": 5.9, "active": true}'
cat file.json | firestore-cli create users 1234`,
		Args:    cobra.MinimumNArgs(2),
		PreRunE: a.Initialize,
		RunE:    a.runCreate,
	}

	a.addHelpFlag()

	return a
}

func (a *Action) runCreate(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	collection := args[0]
	documentID := args[1]

	var jsonValue string
	if len(args) >= 3 {
		jsonValue = args[2]
	} else if a.shouldReadFromStdin() {
		var err error
		jsonValue, err = a.readFromStdin()
		if err != nil {
			return err
		}
	}

	var u any
	err := json.Unmarshal([]byte(jsonValue), &u)
	if err != nil {
		return err
	}

	switch u.(type) {
	case map[string]any:
		err = a.firestore.Create(collection, documentID, u.(map[string]any))
		if err != nil {
			return err
		}
	default:
		return errors.New("invalid JSON format")
	}

	fmt.Printf("%s/%s successfully created\n", collection, documentID)
	return nil
}
