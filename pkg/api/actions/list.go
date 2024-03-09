package actions

import (
	"fmt"
	store2 "github.com/jhight/firestore-cli/pkg/api/store"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func List(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:     "list <collection> [<field>]",
		Aliases: []string{"l"},
		Short:   "List documents in a collection",
		Long:    "List all documents by a specific field (document ID, by default) in the Firestore collection.",
		Example: `firestore-cli list users
firestore-cli list users 'address.city'
firestore-cli list users 'name' --order-by 'created_at desc' --limit 10`,
		Args:    cobra.MinimumNArgs(1),
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

	collection := args[0]

	field := ""
	if len(args) > 1 {
		field = args[1]
	}

	input := store2.SelectionInput{
		Collection: collection,
		OrderBy:    make([]store2.OrderBy, 0),
	}

	if a.command.Flag(flagOrderBy).Changed {
		orderByInput := a.command.Flag(flagOrderBy).Value.String()

		clauses := strings.Split(orderByInput, ",")
		for _, clause := range clauses {
			direction := store2.Ascending

			clause = strings.TrimSpace(clause)
			if strings.HasSuffix(clause, fmt.Sprintf(" %s", store2.Descending)) {
				direction = store2.Descending
				clause = strings.TrimSuffix(clause, fmt.Sprintf(" %s", store2.Descending))
			} else if strings.HasSuffix(clause, fmt.Sprintf(" %s", store2.Ascending)) {
				direction = store2.Ascending
				clause = strings.TrimSuffix(clause, fmt.Sprintf(" %s", store2.Ascending))
			}

			orderBy := store2.OrderBy{
				Field:     strings.TrimSpace(clause),
				Direction: store2.Direction(strings.TrimSpace(string(direction))),
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
