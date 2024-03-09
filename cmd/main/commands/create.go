package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	create := &cobra.Command{
		Use:   "create <collection> <document> [<json>]",
		Short: "Create a document",
		Long:  "Set (replace or create) an entire Firestore document with the specified ID using the specified field(s). If a document exists with the same ID, it will be replaced.",
		Example: `firestore-cli set users 1234 '{"name": "John Doe", "age": 30, "height": 5.9, "active": true}'
cat file.json | firestore-cli create users 1234`,
		Args:   cobra.MinimumNArgs(2),
		PreRun: runRootCommand,
		Run: func(cmd *cobra.Command, args []string) {
			runCreateCommand(cmd, args)
		},
	}

	addHelpFlag(create)

	return create
}

func runCreateCommand(cmd *cobra.Command, args []string) {
	handleHelpFlag(cmd)

	collection := args[0]
	documentID := args[1]

	var jsonValue string
	if len(args) >= 3 {
		jsonValue = args[2]
	} else if shouldReadFromStdin(cmd) {
		var err error
		jsonValue, err = readFromStdin(cmd)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	}

	var u any
	err := json.Unmarshal([]byte(jsonValue), &u)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	switch u.(type) {
	case map[string]any:
		err = firestore.Create(collection, documentID, u.(map[string]any))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	default:
		fmt.Println("Error: invalid JSON value")
	}

	fmt.Printf("%s/%s successfully created\n", collection, documentID)
}
