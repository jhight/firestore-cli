package client

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"jhight.com/firestore-cli/pkg/api/client/query"
)

func (f *firestoreClientManager) Query(input query.Input) ([]map[string]any, error) {
	var q firestore.Query

	if len(input.Filter) == 0 {
		q = f.client.Collection(input.Path).Offset(input.Offset)
	} else {
		root, err := query.CreateExpression(input.Filter)
		if err != nil {
			return nil, err
		}

		q = f.client.
			Collection(input.Path).
			WhereEntity(root.FirestoreFilter())
	}

	if len(input.OrderBy) > 0 {
		for _, o := range input.OrderBy {
			q = q.OrderBy(o.Field, o.Direction.FirestoreDirection())
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
		if len(input.Fields) == 0 {
			documents = append(documents, d.Data())
		} else {
			projection := f.projection(d, input.Fields)
			documents = append(documents, projection)
		}
	}

	return documents, nil
}
