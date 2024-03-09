package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"slices"
)

func NewSetCommand() *cobra.Command {
	update := &cobra.Command{
		Use:   "set <collection> <document> [<json>]",
		Short: "Set a document",
		Long:  "Set (e.g., create or replace) the entire specified Firestore document with specified JSON data. Only the specified fields will exist in the document. If the document does not exist, it will be created.",
		Example: `firestore-cli set users 1234 '{"name": "John Doe", "age": 30, "height": 5.9, "active": true}'
cat file.json | firestore-cli set users 1234`,
		Args:   cobra.MinimumNArgs(2),
		PreRun: runRootCommand,
		Run: func(cmd *cobra.Command, args []string) {
			runSetCommand(cmd, args)
		},
	}

	addHelpFlag(update)

	return update
}

func runSetCommand(cmd *cobra.Command, args []string) {
	handleHelpFlag(cmd)

	collection := args[0]
	documentID := args[1]

	var input string
	if len(args) == 3 {
		input = args[2]
	} else if shouldReadFromStdin(cmd) {
		var err error
		input, err = readFromStdin(cmd)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	}

	if len(input) == 0 {
		fmt.Println("Error: one or more fields in JSON format are required")
		return
	}

	var fields map[string]any
	err := json.Unmarshal([]byte(input), &fields)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	// backup before update, if configured
	if slices.Contains(cfg.Backup.Commands, "update") {
		before, _ := firestore.Get(collection, documentID)
		err = firestore.Set(collection, documentID, fields)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		after, _ := firestore.Get(collection, documentID)
		backup(collection, documentID, before, after)
	} else {
		err = firestore.Set(collection, documentID, fields)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	}

	fmt.Printf("%s/%s successfully set\n", collection, documentID)
}
