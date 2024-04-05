package client

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"jhight.com/firestore-cli/pkg/api/client/query"
)

func (f *firestoreClientManager) Get(input query.Input) (map[string]any, error) {
	var ds *firestore.DocumentSnapshot
	var err error

	d := f.client.Doc(input.Path)
	if d == nil {
		return nil, fmt.Errorf("invalid document path, %s", input.Path)
	}

	ds, err = d.Get(f.ctx)
	if len(input.Fields) == 0 {
		return ds.Data(), err
	}

	projection := f.projection(ds, input.Fields)
	return projection, nil
}

func (f *firestoreClientManager) projection(ds *firestore.DocumentSnapshot, fields []string) map[string]any {
	document := ds.Data()
	projection := make(map[string]any)
	for _, field := range fields {
		if field == query.SelectionDocumentID {
			projection[field] = ds.Ref.ID
			continue
		}
		if field == query.SelectionDocumentPath {
			projection[field] = ds.Ref.Path
			continue
		}

		value, ok := fieldValue(document, field)
		if !ok {
			continue
		}
		projection[field] = value
	}

	return projection
}
