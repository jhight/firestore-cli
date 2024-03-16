package client

import (
	"fmt"
)

func (f *firestoreClientManager) List(input SelectionInput, fieldPath string) ([]any, error) {
	p := make([]string, 0)
	if len(fieldPath) > 0 {
		p = append(p, fieldPath)
	}

	if len(input.CollectionPath) == 0 {
		return collections(f.ctx, f.client), nil
	}

	cr := f.client.Collection(input.CollectionPath)
	if cr == nil {
		return nil, fmt.Errorf("invalid collection, %s", input.CollectionPath)
	}

	query := cr.Select(p...).Offset(input.Offset)

	if input.Limit > 0 {
		query = query.Limit(input.Limit)
	}

	for _, o := range input.OrderBy {
		query = query.OrderBy(o.Field, o.Direction.toFirestoreDirection())
	}

	iter := query.Documents(f.ctx)

	result := make([]any, 0)
	for {
		ds, err := iter.Next()
		if err != nil {
			break
		}

		if len(fieldPath) == 0 {
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
