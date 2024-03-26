package client

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"strings"
)

func count(ctx context.Context, client *firestore.Client, collectionPath string) int {
	cr := client.Collection(collectionPath)
	if cr == nil {
		return 0
	}

	iter := cr.Documents(ctx)
	c := 0
	for {
		_, err := iter.Next()
		if err != nil {
			break
		}
		c++
	}
	return c
}

func create[T any](ctx context.Context, client *firestore.Client, documentPath string, data T) error {
	dr := client.Doc(documentPath)
	if dr == nil {
		return fmt.Errorf("invalid document path, %s", documentPath)
	}

	if _, err := dr.Create(ctx, data); err != nil {
		return fmt.Errorf("error creating document, %s", err)
	}

	return nil
}

func set[T any](ctx context.Context, client *firestore.Client, documentPath string, data T) error {
	dr := client.Doc(documentPath)
	if dr == nil {
		return fmt.Errorf("invalid document path, %s", documentPath)
	}

	if _, err := dr.Set(ctx, data); err != nil {
		return fmt.Errorf("error setting document contents, %s", err)
	}

	return nil
}

func update[T any](ctx context.Context, client *firestore.Client, documentPath string, fields map[string]T) error {
	dr := client.Doc(documentPath)
	if dr == nil {
		return fmt.Errorf("invalid document path, %s", documentPath)
	}

	updates := make([]firestore.Update, 0)
	for k, v := range fields {
		updates = append(updates, firestore.Update{Path: k, Value: v})
	}

	if _, err := dr.Update(ctx, updates); err != nil {
		return fmt.Errorf("error updating document, %s", err)
	}

	return nil
}

func remove(ctx context.Context, client *firestore.Client, path string) error {
	cr := client.Collection(path)
	dr := client.Doc(path)

	if cr != nil {
		err := removeCollection(ctx, client, cr)
		if err != nil {
			return err
		}
		return nil
	}

	if dr != nil {
		err := removeDocument(ctx, client, dr)
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("invalid path format, %s", path)
}

func removeField(ctx context.Context, client *firestore.Client, documentPath string, fieldPath string) error {
	dr := client.Doc(documentPath)
	if dr == nil {
		return fmt.Errorf("invalid document path, %s", documentPath)
	}

	updates := []firestore.Update{{Path: fieldPath, Value: firestore.Delete}}
	if _, err := dr.Update(ctx, updates); err != nil {
		return fmt.Errorf("error deleting field, %s", err)
	}

	return nil
}

func removeCollection(ctx context.Context, client *firestore.Client, cr *firestore.CollectionRef) error {
	if cr == nil {
		return fmt.Errorf("invalid collection reference, %s", cr.Path)
	}

	removed := 0

	iter := cr.Documents(ctx)
	for {
		ds, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		if err = removeDocument(ctx, client, ds.Ref); err != nil {
			return fmt.Errorf("error deleting document, %s", err)
		}

		removed++
	}

	if removed == 0 {
		return fmt.Errorf("invalid collection reference, %s", cr.Path)
	}

	return nil
}

func removeDocument(ctx context.Context, client *firestore.Client, dr *firestore.DocumentRef) error {
	if dr == nil {
		return fmt.Errorf("invalid document reference, %s", dr.Path)
	}

	iter := dr.Collections(ctx)
	for {
		cr, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		if err = removeCollection(ctx, client, cr); err != nil {
			return fmt.Errorf("error deleting collection, %s", err)
		}
	}

	if _, err := dr.Delete(ctx); err != nil {
		return fmt.Errorf("error deleting document, %s", err)
	}

	return nil
}

func collections(ctx context.Context, client *firestore.Client, documentPath string) []any {
	var iter *firestore.CollectionIterator

	if len(documentPath) > 0 {
		iter = client.Doc(documentPath).Collections(ctx)
	} else {
		iter = client.Collections(ctx)
	}

	c := make([]any, 0)
	for {
		collection, err := iter.Next()
		if err != nil {
			break
		}
		c = append(c, collection.ID)
	}
	return c
}

func fieldValue(document map[string]any, path string) (any, bool) {
	fields := strings.Split(path, ".")
	return walkDocument(document, fields)
}

func walkDocument(document map[string]any, fields []string) (any, bool) {
	if len(fields) == 0 {
		return nil, false
	}

	value, ok := document[fields[0]]
	if !ok {
		return nil, false
	}

	if len(fields) == 1 {
		return value, true
	}

	if nested, ok := value.(map[string]any); ok {
		return walkDocument(nested, fields[1:])
	}

	return nil, false
}
