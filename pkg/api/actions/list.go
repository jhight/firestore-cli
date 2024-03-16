package actions

import (
	"fmt"
	"github.com/jhight/firestore-cli/pkg/api/client"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

func List(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:     "list [<path> [<field>]]",
		Aliases: []string{"l"},
		Short:   "List collections or documents within a collection",
		Long:    "List all documents in a collection by a specific field (document ID, by default), or all collections if none is specified.",
		Example: strings.ReplaceAll(`%E list users
%E list users 'address.city'
%E list users/1234/orders
%E list users 'name' --order-by 'created_at desc' --limit 10
%E list`, "%E", os.Args[0]),
		Args:    cobra.MinimumNArgs(0),
		PreRunE: a.initializer.Initialize,
		RunE:    a.runList,
	}

	a.addHelpFlag()
	a.command.Flags().StringP(flagOrderBy, "o", "", "Order by expression, including field and direction (asc or desc)")
	a.command.Flags().Int(flagLimit, 0, "Limit expression")
	a.command.Flags().Int(flagOffset, 0, "Offset expression")
	a.command.Flags().BoolP(flagCount, "c", false, "Count the number of documents returned")

	return a
}

func (a *action) runList(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	if len(args) == 0 {
		collections, err := a.initializer.Firestore().List(client.SelectionInput{}, "")
		if err != nil {
			return err
		}
		a.printOutput(collections)
		return nil
	}

	path := args[0]

	field := ""
	if len(args) > 1 {
		field = args[1]
	}

	input := client.SelectionInput{
		CollectionPath: path,
		OrderBy:        make([]client.OrderBy, 0),
	}

	if a.command.Flag(flagOrderBy).Changed {
		orderByInput := a.command.Flag(flagOrderBy).Value.String()

		clauses := strings.Split(orderByInput, ",")
		for _, clause := range clauses {
			direction := client.Ascending

			clause = strings.TrimSpace(clause)
			if strings.HasSuffix(clause, fmt.Sprintf(" %s", client.Descending)) {
				direction = client.Descending
				clause = strings.TrimSuffix(clause, fmt.Sprintf(" %s", client.Descending))
			} else if strings.HasSuffix(clause, fmt.Sprintf(" %s", client.Ascending)) {
				direction = client.Ascending
				clause = strings.TrimSuffix(clause, fmt.Sprintf(" %s", client.Ascending))
			}

			orderBy := client.OrderBy{
				Field:     strings.TrimSpace(clause),
				Direction: client.Direction(strings.TrimSpace(string(direction))),
			}

			input.OrderBy = append(input.OrderBy, orderBy)
		}
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

	documents, err := a.initializer.Firestore().List(input, field)
	if err != nil {
		return err
	}

	if a.command.Flag(flagCount).Value.String() == "true" {
		a.printOutput(map[string]any{"$count": len(documents)})
	} else {
		a.printOutput(documents)
	}

	return nil
}
