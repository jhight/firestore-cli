package actions

import (
	"fmt"
	"github.com/jhight/firestore-cli/pkg/config"
	"github.com/jhight/firestore-cli/pkg/store"
	"github.com/spf13/cobra"
	"time"
)

type Action struct {
	root      RootAction
	firestore store.Store
	cfg       config.Config
	command   *cobra.Command
}

func (a *Action) nestedField(document map[string]any, fields []string) (any, bool) {
	if len(fields) == 0 {
		return nil, false
	}

	value, ok := document[fields[0]]
	if !ok {
		return nil, false
	}

	if len(fields) == 1 {
		return value, true
	}

	if nested, ok := value.(map[string]any); ok {
		return a.nestedField(nested, fields[1:])
	}

	return nil, false
}

func (a *Action) printOutput(value any) {
	switch value.(type) {
	case map[string]any, []any, []map[string]any:
		json, err := a.toJSON(value)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		fmt.Println(json)
	default:
		fmt.Println(value)
	}
}

func (a *Action) backup(collection string, documentID string, before map[string]any, after map[string]any) {
	bc := a.cfg.Backup.Collection
	if len(bc) == 0 {
		bc = defaultBackupCollection
	}

	b := map[string]any{
		"created_at": time.Now(),
		"collection": collection,
		"document":   documentID,
		"before":     before,
		"after":      after,
	}

	bi := fmt.Sprintf("%d", time.Now().UnixMilli())

	err := a.firestore.Create(bc, bi, b)
	if err != nil {
		fmt.Printf("Failed to create backup: %s\n", err)
	}
}
