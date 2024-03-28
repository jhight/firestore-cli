package actions

import (
	"github.com/jhight/firestore-cli/pkg/api/client/query"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

func Collections(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:   "collections [<path>]",
		Short: "List collections (or subcollections) at path",
		Long:  "List all collections (or subcollections) at the specified path. If no path is provided, all root collections will be returned.",
		Example: strings.ReplaceAll(`%E collections users
%E collections`, "%E", os.Args[0]),
		Args:    cobra.MinimumNArgs(0),
		PreRunE: a.initializer.Initialize,
		RunE:    a.runCollections,
	}

	a.addHelpFlag()
	a.command.Flags().Int(flagLimit, 0, "Limit expression")
	a.command.Flags().Int(flagOffset, 0, "Offset expression")
	a.command.Flags().BoolP(flagCount, "c", false, "Count the number of documents returned")

	return a
}

func (a *action) runCollections(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	if len(args) == 0 {
		collections, err := a.initializer.Firestore().Collections(query.Input{})
		if err != nil {
			return err
		}

		if a.command.Flag(flagCount).Value.String() == "true" {
			a.printOutput(map[string]any{"$count": len(collections)})
		} else {
			a.printOutput(collections)
		}

		return nil
	}

	path := args[0]

	input := query.Input{
		Path:    path,
		OrderBy: make([]query.OrderBy, 0),
	}

	if a.command.Flag(flagLimit).Changed {
		limit, err := strconv.Atoi(a.command.Flag(flagLimit).Value.String())
		if err != nil {
			return err
		}
		input.Limit = limit
	}

	if a.command.Flag(flagOffset).Changed {
		offset, err := strconv.Atoi(a.command.Flag(flagOffset).Value.String())
		if err != nil {
			return err
		}
		input.Offset = offset
	}

	collections, err := a.initializer.Firestore().Collections(input)
	if err != nil {
		return err
	}

	if a.command.Flag(flagCount).Value.String() == "true" {
		if a.initializer.Config().Flatten {
			a.printOutput(len(collections))
		} else {
			a.printOutput(map[string]any{"$count": len(collections)})
		}
	} else {
		a.printOutput(collections)
	}

	return nil
}
