package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

func NewGetCommand() *cobra.Command {
	get := &cobra.Command{
		Use:     "get <collection> <document> [<field>]",
		Aliases: []string{"g", "r"},
		Short:   "Get a document by ID or a field within",
		Long:    "Get either an entire Firestore document from a collection by its ID or the specified field.",
		Example: `firestore-cli get users 1234
firestore-cli get users 1234 name
firestore-cli get users 1234 address.city`,
		Args:   cobra.MinimumNArgs(2),
		PreRun: runRootCommand,
		Run: func(cmd *cobra.Command, args []string) {
			runGetCommand(cmd, args)
		},
	}

	addHelpFlag(get)

	return get
}

func runGetCommand(cmd *cobra.Command, args []string) {
	handleHelpFlag(cmd)

	collection := args[0]
	documentID := args[1]
	field := ""

	if len(args) > 2 {
		field = args[2]
	}

	document, err := firestore.Get(collection, documentID)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	if len(field) > 0 {
		// the firestore sdk makes it difficult to select nested fields when getting a document by id
		if strings.Contains(field, ".") {
			fields := strings.Split(field, ".")
			value, ok := nestedField(document, fields)
			if !ok {
				fmt.Printf("Error: field %s does not exist in %s/%s\n", field, collection, documentID)
				return
			}
			printOutput(value)
			return
		}

		value, ok := document[field]
		if !ok {
			fmt.Printf("Error: field %s does not exist in %s/%s\n", field, collection, documentID)
			return
		}
		printOutput(value)
		return
	}

	printOutput(document)
}
