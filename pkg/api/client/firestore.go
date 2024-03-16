package client

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

func (f *firestoreClientManager) Create(path string, fields map[string]any) error {
	return create(f.ctx, f.client, path, fields)
}

func (f *firestoreClientManager) Get(path string) (map[string]any, error) {
	var result map[string]any
	err := get(f.ctx, f.client, path, &result)
	return result, err
}

func (f *firestoreClientManager) Set(path string, fields map[string]any) error {
	return set(f.ctx, f.client, path, fields)
}

func (f *firestoreClientManager) Update(path string, fields map[string]any) error {
	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	return update(f.ctx, f.client, path, fields)
}

func (f *firestoreClientManager) Delete(path string) error {
	return remove(f.ctx, f.client, path)
}

func (f *firestoreClientManager) Close() error {
	return f.client.Close()
}
