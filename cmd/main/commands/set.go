package commands

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

func NewSetCommand() *cobra.Command {
	update := &cobra.Command{
		Use:     "set <collection> <document> [<field>]",
		Aliases: []string{"s"},
		Short:   "Set a document",
		Long:    "Set (replace or create) an entire Firestore document with the specified ID using the specified field(s). If a document exists with the same ID, it will be replaced.",
		Example: `firestore-cli set users 1234 '{"name": "John Doe", "age": 30, "height": 5.9, "active": true}'
cat file.json | firestore-cli set users 1234`,
		Args:   cobra.MinimumNArgs(2),
		PreRun: runRootCommand,
		Run: func(cmd *cobra.Command, args []string) {
			runSetCommand(cmd, args)
		},
	}

	addHelpFlag(update)
	update.Flags().String(flagJSON, "", "A JSON value as a string (optional, default)")
	update.Flags().Int(flagInt, 0, "An integer value")
	update.Flags().Float64(flagFloat, 0.0, "A float value")
	update.Flags().String(flagString, "", "A string value")
	update.Flags().Bool(flagBool, false, "A boolean value")

	return update
}

func runSetCommand(cmd *cobra.Command, args []string) {
	handleHelpFlag(cmd)

	collection := args[0]
	documentID := args[1]

	var field string
	if len(args) > 3 {
		field = args[2]
	}

	var jsonValue string
	if len(args) == 4 {
		jsonValue = args[3]
	} else if len(args) == 3 {
		jsonValue = args[2]
	} else if cmd.Flag(flagJSON).Changed {
		jsonValue = cmd.Flag(flagJSON).Value.String()
	} else if shouldReadFromStdin(cmd) {
		var err error
		jsonValue, err = readFromStdin(cmd)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	}

	if len(jsonValue) > 0 {
		var u any
		err := json.Unmarshal([]byte(jsonValue), &u)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		switch u.(type) {
		case map[string]any:
			fields := make(map[string]any)
			if len(field) > 0 {
				fields[field] = u.(map[string]any)
			} else {
				fields = u.(map[string]any)
			}
			err = firestore.Set(collection, documentID, fields)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return
			}
		case []any:
			fields := make(map[string]any)
			fields[field] = u
			err = firestore.Set(collection, documentID, fields)
		default:
			fmt.Println("Error: invalid JSON value")
		}
		return
	}

	if len(field) == 0 {
		fmt.Printf("Error: field argument is required when not using --%s\n", flagJSON)
		return
	}

	if cmd.Flag(flagInt).Changed {
		u, _ := cmd.Flags().GetInt(flagInt)
		setFields(collection, documentID, field, u)
	} else if cmd.Flag(flagFloat).Changed {
		u, _ := cmd.Flags().GetFloat64(flagFloat)
		setFields(collection, documentID, field, u)
	} else if cmd.Flag(flagString).Changed {
		u := cmd.Flag(flagString).Value.String()
		setFields(collection, documentID, field, u)
	} else if cmd.Flag(flagBool).Changed {
		u, _ := cmd.Flags().GetBool(flagBool)
		setFields(collection, documentID, field, u)
	} else {
		fmt.Printf("Error: no value provided for field %s\n", field)
		return
	}

	fmt.Printf("%s/%s successfully set\n", collection, documentID)
}

func setFields(collection string, documentID string, field string, value any) {
	err := firestore.Set(collection, documentID, map[string]any{field: value})
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
}
