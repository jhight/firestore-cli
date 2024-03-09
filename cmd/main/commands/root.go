package commands

import (
	"context"
	"fmt"
	"github.com/jhight/firestore-cli/pkg/config"
	"github.com/jhight/firestore-cli/pkg/store"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
)

var firestore store.Store
var cfg config.Config

const defaultConfigPath = "~/.firestore-cli.yaml"
const defaultSpacing = 2

const (
	flagConfigFile     = "config"
	flagServiceAccount = "service-account"
	flagProjectID      = "project-id"
	flagPrettyPrint    = "pretty"
	flagRawPrint       = "raw"
	flagSpacing        = "spacing"
)

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:     "firestore-cli",
		Short:   "A Firebase Firestore command line interface",
		Long:    "A command line interface for Firebase Firestore, allowing querying and CRUD operations on collections and documents.",
		Example: `firestore-cli get accounts account-1234`,
		Run: func(cmd *cobra.Command, args []string) {
			runRootCommand(cmd, args)
		},
	}

	addHelpFlag(root)
	root.PersistentFlags().String(flagConfigFile, defaultConfigPath, "The file path to the firestore-cli configuration file")
	root.PersistentFlags().StringP(flagServiceAccount, "s", "", "The file path to the Google Cloud Platform service account JSON file")
	root.PersistentFlags().StringP(flagProjectID, "p", "", "The Google Cloud Platform project ID")
	root.PersistentFlags().Bool(flagPrettyPrint, true, "Pretty print JSON output")
	root.PersistentFlags().Bool(flagRawPrint, false, "Raw print JSON output (disables pretty print)")
	root.PersistentFlags().Int(flagSpacing, defaultSpacing, "The number of spaces to use for pretty printing JSON output")

	root.AddCommand(NewGetCommand())
	root.AddCommand(NewUpdateCommand())
	root.AddCommand(NewSetCommand())
	root.AddCommand(NewDeleteCommand())
	root.AddCommand(NewCreateCommand())
	root.AddCommand(NewQueryCommand())

	return root
}

func runRootCommand(cmd *cobra.Command, _ []string) {
	handleHelpFlag(cmd)

	// otherwise, load config
	var err error
	cfg, err = loadConfig(cmd)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}

	// start firestore client
	firestore, err = store.New(context.Background(), cfg)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func loadConfig(cmd *cobra.Command) (config.Config, error) {
	// try to read config file
	var err error
	path, _ := cmd.Flags().GetString(flagConfigFile)
	err = readConfigFile(path)

	// create new config if unable to read file
	if err != nil {
		cfg = config.Config{}
	}

	// override config file data with command-line flags
	if cmd.Flag(flagServiceAccount).Changed && len(cmd.Flag(flagServiceAccount).Value.String()) > 0 {
		cfg.ServiceAccount = cmd.Flag(flagServiceAccount).Value.String()
	}
	if cmd.Flag(flagProjectID).Changed && len(cmd.Flag(flagProjectID).Value.String()) > 0 {
		cfg.ProjectID = cmd.Flag(flagProjectID).Value.String()
	}
	if cmd.Flag(flagPrettyPrint).Changed && len(cmd.Flag(flagPrettyPrint).Value.String()) > 0 {
		cfg.PrettyPrint = cmd.Flag(flagPrettyPrint).Value.String() == "true"
	}
	if cmd.Flag(flagRawPrint).Changed && len(cmd.Flag(flagRawPrint).Value.String()) > 0 {
		cfg.Raw = cmd.Flag(flagRawPrint).Value.String() == "true"
	}
	if cmd.Flag(flagSpacing).Changed && len(cmd.Flag(flagSpacing).Value.String()) > 0 {
		amt, _ := strconv.Atoi(cmd.Flag(flagSpacing).Value.String())
		cfg.PrettySpacing = amt
	}

	// make sure required fields are set
	if len(cfg.ServiceAccount) == 0 {
		return config.Config{}, fmt.Errorf("service account file path must either be defined in config file (%s) or provided as flag (--%s=...)", path, flagServiceAccount)
	}
	if len(cfg.ProjectID) == 0 {
		return config.Config{}, fmt.Errorf("project ID must either be defined in config file (%s) or provided as flag (--%s=...)", path, flagProjectID)
	}

	return cfg, nil
}

func readConfigFile(path string) error {
	home := os.Getenv("HOME")
	path = strings.ReplaceAll(path, "~", home)

	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading config file, %s", err)
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return fmt.Errorf("error unmarshalling config file, %s", err)
	}

	return nil
}
