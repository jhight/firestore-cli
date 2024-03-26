package actions

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const defaultConfigPath = "~/.firestore-cli.yaml"
const defaultSpacing = 2
const defaultBackupCollection = "backup"

const (
	flagConfigFile     = "config"
	flagServiceAccount = "service-account"
	flagProjectID      = "project-id"
	flagPrettyPrint    = "pretty"
	flagRawPrint       = "raw"
	flagSpacing        = "spacing"
	flagOrderBy        = "order"
	flagLimit          = "limit"
	flagOffset         = "offset"
	flagCount          = "count"
)

func Root(i Initializer) Action {
	if i == nil {
		i = &initializer{}
	}

	root := &action{
		initializer: i,
	}

	root.command = &cobra.Command{
		Use:   os.Args[0],
		Short: "A Firebase Firestore command line utility",
		Long:  fmt.Sprintf("A command line utility for Firebase Firestore, allowing querying and CRUD operations on collections and documents. For more information, use help on a specific command. For example, %s get --help", os.Args[0]),
		Example: strings.ReplaceAll(`- list collections
	%E collections

- get a single document by ID
    %E get users/user-1234

- list document subcollections
	%E collections users/user-1234

- get specific fields
	%E get users/user-1234 name,age

- get nested fields
	%E get users/user-1234 address.city

- query for data using a filter expression (see %E get --help for information on query syntax)
	%E get users --filter '{"$and":{"name":"John", "age":{">":30}}}'

- create a new document
	%E create users '{"id":1234,"name":"John","age":30}'

- update a document (only specified fields are updated)
	%E update users/user-1234 '{"age":30}'

- set a document (replaces existing document)
	%E set users/user-1234 '{"id":1234,"name":"John","age":30}'

- delete document
	%E delete users/user-1234`, "%E", os.Args[0]),
		PreRunE: i.Initialize,
		RunE:    root.run,
	}

	root.addHelpFlag()
	root.command.PersistentFlags().String(flagConfigFile, defaultConfigPath, "The file path to the firestore-cli configuration file")
	root.command.PersistentFlags().StringP(flagServiceAccount, "s", "", "The file path to the Google Cloud Platform service account JSON file")
	root.command.PersistentFlags().StringP(flagProjectID, "p", "", "The Google Cloud Platform project ID")
	root.command.PersistentFlags().Bool(flagPrettyPrint, true, "Pretty print JSON output")
	root.command.PersistentFlags().Bool(flagRawPrint, false, "Raw print JSON output (disables pretty print)")
	root.command.PersistentFlags().Int(flagSpacing, defaultSpacing, "The number of spaces to use for pretty printing JSON output")

	return root
}

func (a *action) run(_ *cobra.Command, _ []string) error {
	a.handleHelpFlag()
	return nil
}
