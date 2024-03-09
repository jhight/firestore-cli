package store

import (
	"fmt"
)

func (f *firestoreClientManager) List(input SelectionInput, path string) ([]any, error) {
	p := make([]string, 0)
	if len(path) > 0 {
		p = append(p, path)
	}

	query := f.client.Collection(input.Collection).Select(p...).Offset(input.Offset)

	if input.Limit > 0 {
		query = query.Limit(input.Limit)
	}

	for _, o := range input.OrderBy {
		query = query.OrderBy(o.Field, toFirestoreDirection(o.Direction))
	}

	iter := query.Documents(f.ctx)

	result := make([]any, 0)
	for {
		ds, err := iter.Next()
		if err != nil {
			break
		}

		if len(path) == 0 {
			result = append(result, ds.Ref.ID)
			continue
		}

		var value map[string]any
		if err = ds.DataTo(&value); err != nil {
			return nil, fmt.Errorf("error decoding document, %s", err)
		}
		result = append(result, value)
	}

	return result, nil
}
