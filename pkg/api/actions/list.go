package actions

import (
	"fmt"
	"github.com/jhight/firestore-cli/pkg/api/store"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func List(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:     "list [<collection> [<field>]]",
		Aliases: []string{"l"},
		Short:   "List documents or collections",
		Long:    "List all documents in a collection by a specific field (document ID, by default), or all collections if none is specified.",
		Example: `firestore-cli list users
firestore-cli list users 'address.city'
firestore-cli list users 'name' --order-by 'created_at desc' --limit 10
firestore-cli list`,
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
		collections, err := a.initializer.Firestore().List(store.SelectionInput{}, "")
		if err != nil {
			return err
		}
		a.printOutput(collections)
		return nil
	}

	collection := args[0]

	field := ""
	if len(args) > 1 {
		field = args[1]
	}

	input := store.SelectionInput{
		Collection: collection,
		OrderBy:    make([]store.OrderBy, 0),
	}

	if a.command.Flag(flagOrderBy).Changed {
		orderByInput := a.command.Flag(flagOrderBy).Value.String()

		clauses := strings.Split(orderByInput, ",")
		for _, clause := range clauses {
			direction := store.Ascending

			clause = strings.TrimSpace(clause)
			if strings.HasSuffix(clause, fmt.Sprintf(" %s", store.Descending)) {
				direction = store.Descending
				clause = strings.TrimSuffix(clause, fmt.Sprintf(" %s", store.Descending))
			} else if strings.HasSuffix(clause, fmt.Sprintf(" %s", store.Ascending)) {
				direction = store.Ascending
				clause = strings.TrimSuffix(clause, fmt.Sprintf(" %s", store.Ascending))
			}

			orderBy := store.OrderBy{
				Field:     strings.TrimSpace(clause),
				Direction: store.Direction(strings.TrimSpace(string(direction))),
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
