package actions

import (
	"fmt"
	"github.com/jhight/firestore-cli/pkg/api/client"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

const (
	flagOrderBy = "order-by"
	flagLimit   = "limit"
	flagOffset  = "offset"
	flagCount   = "count"
)

func Query(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:     "query <collection> [<query>]",
		Aliases: []string{"q"},
		Short:   "Execute a query",
		Long:    "Execute a query against a Firestore collection. See examples below for more information about query JSON syntax. If no query is provided, all documents in the collection will be returned.",
		Args:    cobra.MinimumNArgs(1),
		Example: strings.ReplaceAll(`- gets all users with id 1234
    %E query users '{"id":{"==":1234}}'

- shorter, also gets all users with id 1234 ("==" is the default field operator)
    %E query users '{"id":1234}'

- same as above, ordered by age descending then name ascending, and limited to 10
    %E query users '{"id":1234}' --order-by "age desc, name asc" --limit 10

- get all orders by user 1234 over $100 (orders is a subcollection of users)
	%E query users/1234/orders '{"price":{">":100}}'

- gets all users with id 1234 and age > 30
    %E query users '{"$and":{"id":1234,"age":{">":30}}}'

- shorter, also gets all users with id 1234 and age > 30 ("$and" is the default composite operator)
    %E query users '{"id":1234,"age":{">":30}}'

- complex filter: a = "abc" and b > 30 and (c = true or (d <= 25 and e != "def"))
    %E query users '{"$and":{"a":"abc","b":{">":30},"$or":{"c":true,"$and":{"d":{"<=":25},"e":{"!=":"def"}}}}'

- shorter version of the above, without explicit outer $and composite operator
    %E query users '{"a":"abc","b":{">":30},"$or":{"c":true,"$and":{"d":{"<=":25},"e":{"!=":"def"}}}}'

- get all users where address city is one of: "New York", "Los Angeles", or "Chicago"
	%E query users '{"address.city":{"$in":["New York","Los Angeles","Chicago"]}}'

- get the count of all users with address.city of "New York"
    %E query users '{"address.city":"New York"}' --count

- execute query from stdin
    cat query.json | %E query users

- get all the id of all users, ordered by name and limited to 10
    %E query users --order-by "name asc" --limit 10`, "%E", os.Args[0]),
		PreRunE: a.initializer.Initialize,
		RunE:    a.runQuery,
	}

	a.addHelpFlag()
	a.command.Flags().StringP(flagOrderBy, "o", "", "Order by expression, including field and direction (asc or desc)")
	a.command.Flags().Int(flagLimit, 0, "Limit expression")
	a.command.Flags().Int(flagOffset, 0, "Offset expression")
	a.command.Flags().BoolP(flagCount, "c", false, "Count the number of documents returned")

	return a
}

func (a *action) runQuery(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	path := args[0]

	var query string
	if len(args) > 1 {
		query = args[1]
	} else if a.shouldReadFromStdin() {
		var err error
		query, err = a.readFromStdin()
		if err != nil {
			return err
		}
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

	documents, err := a.initializer.Firestore().Query(input, query)
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
