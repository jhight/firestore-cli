package actions

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

type RootAction interface {
	Add(...*Action)
	Execute() error
	Initialize(cmd *cobra.Command, args []string) error
	Firestore() store.Store
	Config() config.Config
}

func Root() RootAction {
	root := &Action{}

	root.command = &cobra.Command{
		Use:     "firestore-cli",
		Short:   "A Firebase Firestore command line interface",
		Long:    "A command line interface for Firebase Firestore, allowing querying and CRUD operations on collections and documents.",
		Example: `firestore-cli get accounts account-1234`,
		PreRunE: root.Initialize,
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

func (a *Action) Add(actions ...*Action) {
	for _, action := range actions {
		a.command.AddCommand(action.command)
	}
}

func (a *Action) Execute() error {
	return a.command.Execute()
}

func (a *Action) Firestore() store.Store {
	return a.firestore
}

func (a *Action) Config() config.Config {
	return a.cfg
}

func (a *Action) run(_ *cobra.Command, _ []string) error {
	a.handleHelpFlag()
	return nil
}

func (a *Action) Initialize(cmd *cobra.Command, _ []string) error {
	var err error
	a.cfg, err = a.loadConfig(cmd)
	if err != nil {
		return err
	}

	a.firestore, err = store.New(context.Background(), a.cfg)
	if err != nil {
		return err
	}

	return nil
}

func (a *Action) loadConfig(cmd *cobra.Command) (config.Config, error) {
	// try to read config file
	var err error
	path, _ := cmd.Flags().GetString(flagConfigFile)
	err = a.readConfigFile(path)

	// create new config if unable to read file
	if err != nil {
		a.cfg = config.Config{}
	}

	// override config file data with command-line flags
	if cmd.Flag(flagServiceAccount).Changed && len(cmd.Flag(flagServiceAccount).Value.String()) > 0 {
		a.cfg.ServiceAccount = cmd.Flag(flagServiceAccount).Value.String()
	}
	if cmd.Flag(flagProjectID).Changed && len(cmd.Flag(flagProjectID).Value.String()) > 0 {
		a.cfg.ProjectID = cmd.Flag(flagProjectID).Value.String()
	}
	if cmd.Flag(flagPrettyPrint).Changed && len(cmd.Flag(flagPrettyPrint).Value.String()) > 0 {
		a.cfg.PrettyPrint = cmd.Flag(flagPrettyPrint).Value.String() == "true"
	}
	if cmd.Flag(flagRawPrint).Changed && len(cmd.Flag(flagRawPrint).Value.String()) > 0 {
		a.cfg.Raw = cmd.Flag(flagRawPrint).Value.String() == "true"
	}
	if cmd.Flag(flagSpacing).Changed && len(cmd.Flag(flagSpacing).Value.String()) > 0 {
		amt, _ := strconv.Atoi(cmd.Flag(flagSpacing).Value.String())
		a.cfg.PrettySpacing = amt
	}

	// make sure required fields are set
	if len(a.cfg.ServiceAccount) == 0 {
		return config.Config{}, fmt.Errorf("service account file path must either be defined in config file (%s) or provided as flag (--%s=...)", path, flagServiceAccount)
	}
	if len(a.cfg.ProjectID) == 0 {
		return config.Config{}, fmt.Errorf("project ID must either be defined in config file (%s) or provided as flag (--%s=...)", path, flagProjectID)
	}

	return a.cfg, nil
}

func (a *Action) readConfigFile(path string) error {
	home := os.Getenv("HOME")
	path = strings.ReplaceAll(path, "~", home)

	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading config file, %s", err)
	}

	err = yaml.Unmarshal(file, &a.cfg)
	if err != nil {
		return fmt.Errorf("error unmarshalling config file, %s", err)
	}

	return nil
}
