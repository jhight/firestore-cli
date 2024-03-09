package commands

import (
	"fmt"
	"github.com/jhight/firestore-cli/pkg/store"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <collection> [<field>]",
		Aliases: []string{"l"},
		Short:   "List documents in a collection",
		Long:    "List all documents by a specific field (document ID, by default) in the Firestore collection.",
		Example: `firestore-cli list users
firestore-cli list users 'address.city'
firestore-cli list users 'name' --order-by 'created_at desc' --limit 10`,
		Args:   cobra.MinimumNArgs(1),
		PreRun: runRootCommand,
		Run: func(cmd *cobra.Command, args []string) {
			runListCommand(cmd, args)
		},
	}

	addHelpFlag(cmd)
	cmd.Flags().StringP(flagOrderBy, "o", "", "Order by expression, including field and direction (asc or desc)")
	cmd.Flags().Int(flagLimit, 0, "Limit expression")
	cmd.Flags().Int(flagOffset, 0, "Offset expression")
	cmd.Flags().BoolP(flagCount, "c", false, "Count the number of documents returned")

	return cmd
}

func runListCommand(cmd *cobra.Command, args []string) {
	handleHelpFlag(cmd)

	collection := args[0]

	field := ""
	if len(args) > 1 {
		field = args[1]
	}

	input := store.SelectionInput{
		Collection: collection,
		OrderBy:    make([]store.OrderBy, 0),
	}

	if cmd.Flag(flagOrderBy).Changed {
		orderByInput := cmd.Flag(flagOrderBy).Value.String()

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

	if cmd.Flag(flagLimit).Changed {
		limit, err := strconv.Atoi(cmd.Flag(flagLimit).Value.String())
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		input.Limit = limit
	}

	if cmd.Flag(flagOffset).Changed {
		offset, err := strconv.Atoi(cmd.Flag(flagOffset).Value.String())
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
		input.Offset = offset
	}

	documents, err := firestore.List(input, field)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	if cmd.Flag(flagCount).Value.String() == "true" {
		printOutput(map[string]any{"$count": len(documents)})
		return
	}

	printOutput(documents)
}
