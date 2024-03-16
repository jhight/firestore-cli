package client

import (
	"cloud.google.com/go/firestore"
	"encoding/json"
	"fmt"
)

func (f *firestoreClientManager) Query(input SelectionInput, filter string) ([]map[string]any, error) {
	var q firestore.Query

	if len(filter) == 0 {
		// if no filter is provided, return all documents
		q = f.client.Collection(input.CollectionPath).Offset(input.Offset)
	} else {
		// parse the filter and create a query
		var body map[string]any
		err := json.Unmarshal([]byte(filter), &body)
		if err != nil {
			return nil, fmt.Errorf("query parse failure, %s; see help for more information on query syntax", err)
		}

		root, err := createRootExpression(body)
		if err != nil {
			return nil, err
		}

		parse(root)

		q = f.client.
			Collection(input.CollectionPath).
			WhereEntity(root.toEntityFilter())
	}

	if len(input.OrderBy) > 0 {
		for _, o := range input.OrderBy {
			q = q.OrderBy(o.Field, toFirestoreDirection(o.Direction))
		}
	}

	if input.Limit > 0 {
		q = q.Limit(input.Limit)
	}

	iter := q.Documents(f.ctx)
	ds, err := iter.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error querying documents, %s", err)
	}

	documents := make([]map[string]any, 0)
	for _, d := range ds {
		var t map[string]any
		if err = d.DataTo(&t); err != nil {
			return nil, fmt.Errorf("error decoding document, %s", err)
		}
		documents = append(documents, t)
	}

	return documents, nil
}
