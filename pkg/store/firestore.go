package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
)

type firestoreClientManager struct {
	ctx    context.Context
	client *firestore.Client
}

func (f *firestoreClientManager) Count(collection string) int {
	return count(f.ctx, f.client, collection)
}

func (f *firestoreClientManager) Create(collection, document string, fields map[string]any) error {
	return create(f.ctx, f.client, collection, document, fields)
}

func (f *firestoreClientManager) Get(collection, document string) (map[string]any, error) {
	var result map[string]any
	err := get(f.ctx, f.client, collection, document, &result)
	return result, err
}

func (f *firestoreClientManager) Set(collection, document string, fields map[string]any) error {
	return set(f.ctx, f.client, collection, document, fields)
}

func (f *firestoreClientManager) Update(collection, ID string, fields map[string]any) error {
	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	return update(f.ctx, f.client, collection, ID, fields)
}

func (f *firestoreClientManager) Delete(collection, document string) error {
	return remove(f.ctx, f.client, collection, document)

}

func (f *firestoreClientManager) Close() error {
	return f.client.Close()
}
