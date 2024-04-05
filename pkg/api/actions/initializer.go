package actions

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"jhight.com/firestore-cli/pkg/api/client"
	"jhight.com/firestore-cli/pkg/config"
	"os"
	"strconv"
	"strings"
)

type Initializer interface {
	Initialize(cmd *cobra.Command, _ []string) error
	Firestore() client.Store
	Config() config.Config
}

type initializer struct {
	initialized bool
	cfg         config.Config
	firestore   client.Store
}

func DefaultsInitializer(cfg config.Config, firestore client.Store) Initializer {
	return &initializer{
		initialized: true,
		cfg:         cfg,
		firestore:   firestore,
	}
}

func (i *initializer) Initialize(cmd *cobra.Command, _ []string) error {
	if i.initialized {
		return nil
	}

	var err error

	if i.cfg, err = i.loadConfig(cmd); err != nil {
		return err
	}

	if i.firestore, err = client.New(context.Background(), i.cfg); err != nil {
		return err
	}

	i.initialized = true

	return nil
}

func (i *initializer) Firestore() client.Store {
	return i.firestore
}

func (i *initializer) Config() config.Config {
	return i.cfg
}

func (i *initializer) loadConfig(cmd *cobra.Command) (config.Config, error) {
	// try to read config file
	var err error
	path, _ := cmd.Flags().GetString(flagConfigFile)
	err = i.readConfigFile(path)

	// create new config if unable to read file
	if err != nil {
		i.cfg = config.Config{}
	}

	// override config file data with command-line flags
	if cmd.Flag(flagServiceAccount).Changed && len(cmd.Flag(flagServiceAccount).Value.String()) > 0 {
		i.cfg.ServiceAccount = cmd.Flag(flagServiceAccount).Value.String()
	}
	if cmd.Flag(flagProjectID).Changed && len(cmd.Flag(flagProjectID).Value.String()) > 0 {
		i.cfg.ProjectID = cmd.Flag(flagProjectID).Value.String()
	}
	if cmd.Flag(flagPrettyPrint).Changed && len(cmd.Flag(flagPrettyPrint).Value.String()) > 0 {
		i.cfg.PrettyPrint = cmd.Flag(flagPrettyPrint).Value.String() == "true"
	}
	if cmd.Flag(flagRawPrint).Changed && len(cmd.Flag(flagRawPrint).Value.String()) > 0 {
		i.cfg.RawPrint = cmd.Flag(flagRawPrint).Value.String() == "true"
	}
	if cmd.Flag(flagSpacing).Changed && len(cmd.Flag(flagSpacing).Value.String()) > 0 {
		amt, _ := strconv.Atoi(cmd.Flag(flagSpacing).Value.String())
		i.cfg.PrettySpacing = amt
	}
	if cmd.Flag(flagFlatten).Changed && len(cmd.Flag(flagFlatten).Value.String()) > 0 {
		i.cfg.Flatten = cmd.Flag(flagFlatten).Value.String() == "true"
	}

	// make sure required fields are set
	if len(i.cfg.ServiceAccount) == 0 {
		return config.Config{}, fmt.Errorf("service account file path must either be defined in config file (%s) or provided as flag (--%s=...)", path, flagServiceAccount)
	}
	if len(i.cfg.ProjectID) == 0 {
		return config.Config{}, fmt.Errorf("project ID must either be defined in config file (%s) or provided as flag (--%s=...)", path, flagProjectID)
	}

	return i.cfg, nil
}

func (i *initializer) readConfigFile(path string) error {
	home := os.Getenv("HOME")
	path = strings.ReplaceAll(path, "~", home)

	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading config file, %s", err)
	}

	err = yaml.Unmarshal(file, &i.cfg)
	if err != nil {
		return fmt.Errorf("error unmarshalling config file, %s", err)
	}

	return nil
}
