package actions

import (
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
)

func Root(i Initializer) Action {
	if i == nil {
		i = &initializer{}
	}

	root := &action{
		initializer: i,
	}

	root.command = &cobra.Command{
		Use:     os.Args[0],
		Short:   "A Firebase Firestore command line utility",
		Long:    "A command line utility for Firebase Firestore, allowing querying and CRUD operations on collections and documents.",
		Example: strings.ReplaceAll(`%E get accounts account-1234`, "%E", os.Args[0]),
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
