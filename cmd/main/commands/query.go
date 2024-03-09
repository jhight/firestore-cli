package commands

import (
	"fmt"
	"github.com/jhight/firestore-cli/pkg/store"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

const (
	flagOrderBy = "order-by"
	flagLimit   = "limit"
	flagOffset  = "offset"
	flagCount   = "count"
)

func NewQueryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "query <collection> [<query>]",
		Aliases: []string{"q"},
		Short:   "Execute a query",
		Long:    "Execute a query against a Firestore collection. See examples below for more information about query JSON syntax.",
		Args:    cobra.MinimumNArgs(1),
		Example: `- gets all users with id 1234
    firestore-cli query users '{"id":{"==":1234}}'

- shorter, also gets all users with id 1234 ("==" is the default field operator)
    firestore-cli query users '{"id":1234}'

- same as above, ordered by age descending then name ascending, and limited to 10
    firestore-cli query users '{"id":1234}' --order-by "age desc, name asc" --limit 10

- gets all users with id 1234 and age > 30
    firestore-cli query users '{"$and":{"id":1234,"age":{">":30}}}'

- shorter, also gets all users with id 1234 and age > 30 ("$and" is the default composite operator)
    firestore-cli query users '{"id":1234,"age":{">":30}}'

- complex filter: a = "abc" and b > 30 and (c = true or (d <= 25 and e != "def"))
    firestore-cli query users '{"$and":{"a":"abc","b":{">":30},"$or":{"c":true,"$and":{"d":{"<=":25},"e":{"!=":"def"}}}}'

- shorter version of the above, without explicit outer $and composite operator
    firestore-cli query users '{"a":"abc","b":{">":30},"$or":{"c":true,"$and":{"d":{"<=":25},"e":{"!=":"def"}}}}'

- get all users where address city is one of: "New York", "Los Angeles", or "Chicago"
	firestore-cli query users '{"address.city":{"$in":["New York","Los Angeles","Chicago"]}}'

- get the count of all users with address.city of "New York"
    firestore-cli query users '{"address.city":"New York"}' --count

- execute query from stdin
    cat query.json | firestore-cli query users`,
		PreRun: runRootCommand,
		Run: func(cmd *cobra.Command, args []string) {
			runQueryCommand(cmd, args)
		},
	}

	addHelpFlag(cmd)
	cmd.Flags().StringP(flagOrderBy, "o", "", "Order by expression, including field and direction (asc or desc)")
	cmd.Flags().Int(flagLimit, 0, "Limit expression")
	cmd.Flags().Int(flagOffset, 0, "Offset expression")
	cmd.Flags().BoolP(flagCount, "c", false, "Count the number of documents returned")

	return cmd
}

func runQueryCommand(cmd *cobra.Command, args []string) {
	handleHelpFlag(cmd)

	collection := args[0]

	var query string
	if len(args) > 1 {
		query = args[1]
	} else if shouldReadFromStdin(cmd) {
		var err error
		query, err = readFromStdin(cmd)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}
	} else {
		fmt.Println("Error: query is required")
		return
	}

	input := store.QueryInput{
		Collection: collection,
		Query:      query,
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

	documents, err := firestore.Query(input)
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
