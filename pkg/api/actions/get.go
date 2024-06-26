package actions

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"jhight.com/firestore-cli/pkg/api/client/query"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	flagFilter  = "filter"
	flagWhere   = "where"
	flagOrderBy = "order"
	flagLimit   = "limit"
	flagOffset  = "offset"
	flagCount   = "count"
)

func Get(root Action) Action {
	a := &action{
		initializer: root.Initializer(),
	}

	a.command = &cobra.Command{
		Use:     "get <path> [<fields>]",
		Aliases: []string{"g"},
		Short:   "Get data from a collection or document",
		Long:    "Get data from a Firestore collection or document by ID or by applying a filter. See examples below (or README) for more information about query JSON syntax.",
		Args:    cobra.MinimumNArgs(1),
		Example: strings.ReplaceAll(`- get a document by ID
	%E get users/user-1234

- get specific fields from a document
	%E get users/user-1234 name,age

- query for data using a filter expression
	%E get users --filter '{"id":{"==":1234}}'

- shorter, also gets all users with id 1234 ("==" is the default field operator)
	%E get users --filter '{"id":1234}'

- same as above, ordered by age descending then name ascending, and limited to 10
	%E get users --filter '{"id":1234}' --order age:desc,name:asc --limit 10

- get all orders by user 1234 over $100 (orders is a subcollection of users)
	%E get users/1234/orders --filter '{"price":{">":100}}'

- get all users with id 1234 and age > 30
	%E get users --filter '{"$and":{"id":1234,"age":{">":30}}}'

- shorter, also gets all users with id 1234 and age > 30 ("$and" is the default composite operator)
	%E get users --filter '{"id":1234,"age":{">":30}}'

- complex filter: a = "abc" and b > 30 and (c = true or (d <= 25 and e != "def"))
	%E get users --filter '{"$and":{"a":"abc","b":{">":30},"$or":{"c":true,"$and":{"d":{"<=":25},"e":{"!=":"def"}}}}'

- shorter version of the above, without explicit outer $and composite operator
	%E get users --filter '{"a":"abc","b":{">":30},"$or":{"c":true,"$and":{"d":{"<=":25},"e":{"!=":"def"}}}}'

- get all users where address city is one of: "New York", "Los Angeles", or "Chicago"
	%E get users --filter '{"address.city":{"$in":["New York","Los Angeles","Chicago"]}}'

- get the count of all users with address.city of "New York"
	%E get users --filter '{"address.city":"New York"}' --count

- get all the id of all users, ordered by name and limited to 10
	%E get users --order name --limit 10`, "%E", os.Args[0]),
		PreRunE: a.initializer.Initialize,
		RunE:    a.runGet,
	}

	a.addHelpFlag()
	a.command.Flags().StringP(flagFilter, "f", "", "Filter by expression. See examples above.")
	a.command.Flags().StringP(flagWhere, "w", "", fmt.Sprintf("Alias for filter by expression (--%s). See examples above.", flagFilter))
	a.command.Flags().StringP(flagOrderBy, "o", "", "Order by expression, including field and direction (asc or desc). See examples above.")
	a.command.Flags().IntP(flagLimit, "l", 0, "Limit integer value.")
	a.command.Flags().Int(flagOffset, 0, "Offset integer value.")
	a.command.Flags().Bool(flagCount, false, "Return only the count of documents matching query.")

	return a
}

func (a *action) runGet(_ *cobra.Command, args []string) error {
	a.handleHelpFlag()

	path := args[0]
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	fields := make([]string, 0)
	if len(args) > 1 {
		fields = strings.Split(args[1], ",")
	}

	input := query.Input{
		Path:   path,
		Fields: fields,
	}

	filterString := ""
	if a.command.Flag(flagFilter).Changed {
		filterString = a.command.Flag(flagFilter).Value.String()
	} else if a.command.Flag(flagWhere).Changed {
		filterString = a.command.Flag(flagWhere).Value.String()
	}

	if len(filterString) > 0 {
		err := json.Unmarshal([]byte(filterString), &input.Filter)
		if err != nil {
			return fmt.Errorf("query parse failure, %s; see help for more information on query syntax", err)
		}
	}

	if a.command.Flag(flagOrderBy).Changed {
		orderByInput := a.command.Flag(flagOrderBy).Value.String()

		clauses := strings.Split(orderByInput, ",")
		for _, clause := range clauses {
			direction := query.Ascending
			field := clause
			ascSuffix := fmt.Sprintf(":%s", query.Ascending)
			descSuffix := fmt.Sprintf(":%s", query.Descending)

			clause = strings.TrimSpace(clause)
			if strings.HasSuffix(clause, descSuffix) {
				direction = query.Descending
				field = strings.TrimSuffix(clause, descSuffix)
			} else if strings.HasSuffix(clause, ascSuffix) {
				direction = query.Ascending
				field = strings.TrimSuffix(clause, ascSuffix)
			}

			orderBy := query.OrderBy{
				Field:     strings.TrimSpace(field),
				Direction: query.Direction(strings.TrimSpace(string(direction))),
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

	if a.command.Flag(flagCount).Changed {
		input.Count = true
	}

	var err error
	if a.initializer.Firestore().IsPathToCollection(path) {
		var docs []map[string]any
		docs, err = a.initializer.Firestore().Query(input)
		a.handleOutput(false, docs, fields)
	} else if a.initializer.Firestore().IsPathToDocument(path) {
		var doc map[string]any
		doc, err = a.initializer.Firestore().Get(input)
		a.handleOutput(true, []map[string]any{doc}, fields)
	}

	if err != nil && strings.Contains(fmt.Sprintf("%s", err), "The query requires an index. You can create it here") {
		re := regexp.MustCompile("^.*?(https://.*?)$")
		matches := re.FindStringSubmatch(fmt.Sprintf("%s", err))
		if len(matches) > 1 {
			a.printOutput("Requires index to be created: " + matches[1])
			return nil
		}
		return err
	}

	return err
}

func (a *action) handleOutput(isDocument bool, docs []map[string]any, fields []string) {
	if a.command.Flag(flagCount).Value.String() == "true" {
		if a.initializer.Config().Flatten {
			a.printOutput(len(docs))
		} else {
			a.printOutput(map[string]any{"$count": len(docs)})
		}
	} else if isDocument && len(docs) == 1 {
		a.printOutput(docs[0])
	} else if a.initializer.Config().Flatten && len(fields) == 1 {
		flattened := make([]any, 0)
		for _, doc := range docs {
			for _, v := range doc {
				flattened = append(flattened, v)
				break
			}
		}

		if len(flattened) == 1 {
			a.printOutput(flattened[0])
		} else {
			a.printOutput(flattened)
		}
	} else {
		a.printOutput(docs)
	}
}
