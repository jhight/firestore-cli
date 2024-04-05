package client

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"jhight.com/firestore-cli/pkg/api/client/query"
	"strings"
	"time"
)

type firestoreClientManager struct {
	ctx    context.Context
	client *firestore.Client
}

func (f *firestoreClientManager) IsPathToDocument(path string) bool {
	return f.client.Doc(path) != nil
}

func (f *firestoreClientManager) IsPathToCollection(path string) bool {
	return f.client.Collection(path) != nil
}

func (f *firestoreClientManager) Create(path string, fields map[string]any) error {
	return create(f.ctx, f.client, path, processInputValues(fields))
}

func (f *firestoreClientManager) Set(path string, fields map[string]any) error {
	return set(f.ctx, f.client, path, processInputValues(fields))
}

func (f *firestoreClientManager) Update(path string, fields map[string]any) error {
	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	return update(f.ctx, f.client, path, processInputValues(fields))
}

func (f *firestoreClientManager) Delete(path string) error {
	return remove(f.ctx, f.client, path)
}

func (f *firestoreClientManager) DeleteField(path string, field string) error {
	return removeField(f.ctx, f.client, path, field)
}

func (f *firestoreClientManager) Close() error {
	return f.client.Close()
}

func processInputValues(fields map[string]any) map[string]any {
	for k, v := range fields {
		switch (v).(type) {
		case string:
			s := v.(string)
			if strings.HasPrefix(s, query.FunctionTimestamp) {
				val := strings.TrimSuffix(strings.TrimPrefix(v.(string), query.FunctionTimestamp+"("), ")")
				parsed, err := time.Parse(time.RFC3339Nano, val)
				if err != nil {
					fmt.Printf("invalid timestamp format %s; see help for more information on query syntax", val)
					return nil
				} else {
					fields[k] = parsed
				}
			} else if strings.ToLower(s) == query.FunctionNow+"()" {
				fields[k] = time.Now()
			}
		}
	}

	return fields
}
