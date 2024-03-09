package store

import (
	"encoding/json"
	"fmt"
)

func (f *firestoreClientManager) Query(input QueryInput) ([]map[string]any, error) {
	var body map[string]any
	err := json.Unmarshal([]byte(input.Query), &body)
	if err != nil {
		return nil, fmt.Errorf("query parse failure, %s; see help for more information on query syntax", err)
	}

	root, err := createRootExpression(body)
	if err != nil {
		return nil, err
	}

	parse(root)

	q := f.client.
		Collection(input.Collection).
		WhereEntity(root.toEntityFilter())

	if len(input.OrderBy) > 0 {
		for _, o := range input.OrderBy {
			q = q.OrderBy(o.Field, toFirestoreDirection(o.Direction))
		}
	}

	if input.Limit > 0 {
		q = q.Limit(input.Limit)
	}

	if input.Offset > 0 {
		q = q.Offset(input.Offset)
	}

	ds, err := q.Documents(f.ctx).GetAll()

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
