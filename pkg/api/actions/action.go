package actions

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

//go:generate go run go.uber.org/mock/mockgen -typed -package $GOPACKAGE -source $GOFILE -destination $GOFILE.mocks.go
type Action interface {
	SetArgs(args []string)
	Command() *cobra.Command
	Add(actions ...Action)
	Execute() error
	Initializer() Initializer
}

type action struct {
	initializer Initializer
	command     *cobra.Command
}

func (a *action) SetArgs(args []string) {
	a.command.SetArgs(args)
}

func (a *action) Add(actions ...Action) {
	for _, an := range actions {
		a.command.AddCommand(an.Command())
	}
}

func (a *action) Execute() error {
	return a.command.Execute()
}

func (a *action) Initializer() Initializer {
	return a.initializer
}

func (a *action) Command() *cobra.Command {
	return a.command
}

func (a *action) printOutput(value any) {
	switch value.(type) {
	case []any:
		values := make([]any, 0)
		for _, v := range value.([]any) {
			if v == nil {
				continue
			}
			values = append(values, v)
		}
		json, err := a.toJSON(values)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		if len(json) == 0 || json == "[]" {
			return
		}
		fmt.Println(json)
	case []map[string]any:
		values := make([]map[string]any, 0)
		for _, v := range value.([]map[string]any) {
			if len(v) == 0 {
				continue
			}
			values = append(values, v)
		}
		json, err := a.toJSON(values)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		if len(json) == 0 || json == "[]" {
			return
		}
		fmt.Println(json)
	case map[string]any:
		json, err := a.toJSON(value)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		if len(json) == 0 || json == "{}" {
			return
		}
		fmt.Println(json)
	default:
		fmt.Println(value)
	}
}

func (a *action) backup(path string, before map[string]any, after map[string]any) {
	bc := a.initializer.Config().Backup.Collection
	if len(bc) == 0 {
		bc = defaultBackupCollection
	}

	b := map[string]any{
		"created_at": time.Now(),
		"path":       path,
		"before":     before,
		"after":      after,
	}

	bi := fmt.Sprintf("%d", time.Now().UnixMilli())

	err := a.initializer.Firestore().Create(fmt.Sprintf("%s/%s", bc, bi), b)
	if err != nil {
		fmt.Printf("Failed to create backup: %s\n", err)
	}
}
