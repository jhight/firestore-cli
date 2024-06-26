package actions

import (
	"fmt"
	"github.com/spf13/cobra"
	"jhight.com/firestore-cli/pkg/api/client/query"
	"os"
	"slices"
	"strings"
)

func Delete(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:   "delete <path> [<field>]",
		Short: "Delete a collection, document, or field",
		Long:  "Delete a Firestore collection, document, or field.",
		Example: strings.ReplaceAll(`%E delete users/1234
%E delete users
%E delete users/1234 field_to_remove`, "%E", os.Args[0]),
		Args:    cobra.MinimumNArgs(1),
		PreRunE: a.initializer.Initialize,
		RunE:    a.runDelete,
	}

	a.addHelpFlag()

	return a
}

func (a *action) runDelete(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	path := args[0]

	fields := make([]string, 0)
	if len(args) > 1 {
		fields = strings.Split(args[1], ",")
	}

	// if the path is a collection, confirm the deletion
	components := strings.Split(path, "/")
	if len(components)%2 == 1 {
		fmt.Printf("Delete collection %s? (y/N): ", path)
		var response string
		_, _ = fmt.Scanln(&response)
		if !strings.HasPrefix(strings.TrimSpace(strings.ToUpper(response)), "Y") {
			fmt.Println("Deletion cancelled")
			return nil
		}
	}

	if slices.Contains(a.initializer.Config().Backup.Commands, "delete") {
		before, _ := a.initializer.Firestore().Get(query.Input{Path: path})
		a.backup(path, before, nil)
	}

	if len(fields) == 0 {
		if err := a.initializer.Firestore().Delete(path); err != nil {
			return err

		}
		fmt.Printf("%s successfully deleted\n", path)
		return nil
	}

	for _, field := range fields {
		if err := a.initializer.Firestore().DeleteField(path, field); err != nil {
			return err
		}
		fmt.Printf("%s %s successfully deleted\n", path, field)
	}

	return nil
}
