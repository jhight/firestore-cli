package actions

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"slices"
	"strings"
)

func Delete(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:   "delete <path>",
		Short: "Delete a collection or document by ID",
		Long:  "Delete a Firestore collection or document in a collection by its ID.",
		Example: strings.ReplaceAll(`%E delete users/1234
%E delete users`, "%E", os.Args[0]),
		Args:    cobra.ExactArgs(1),
		PreRunE: a.initializer.Initialize,
		RunE:    a.runDelete,
	}

	a.addHelpFlag()

	return a
}

func (a *action) runDelete(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	path := args[0]

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
		before, _ := a.initializer.Firestore().Get(path)
		a.backup(path, before, nil)
	}

	err := a.initializer.Firestore().Delete(path)
	if err != nil {
		return err
	}

	fmt.Printf("%s successfully deleted\n", path)
	return nil
}
