package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
)

func create[T any](ctx context.Context, client *firestore.Client, collection string, document string, data T) error {
	if _, err := client.Collection(collection).Doc(document).Create(ctx, data); err != nil {
		return fmt.Errorf("error creating document, %s", err)
	}

	return nil
}

func get[T any](ctx context.Context, client *firestore.Client, collection string, document string, value *T) error {
	ds, err := client.Collection(collection).Doc(document).Get(ctx)
	if err != nil {
		return fmt.Errorf("error getting document, %s", err)
	}

	if err = ds.DataTo(value); err != nil {
		return fmt.Errorf("error decoding document, %s", err)
	}

	return nil
}

func query[T any](ctx context.Context, client *firestore.Client, collection string, path string, op FieldOperator, value any, orderByProperty string, direction firestore.Direction, decoded *[]T) error {
	filter := firestore.PropertyFilter{
		Path:     path,
		Operator: string(op),
		Value:    value,
	}

	ds, err := client.Collection(collection).WhereEntity(filter).OrderBy(orderByProperty, direction).Documents(ctx).GetAll()
	if err != nil {
		return fmt.Errorf("error querying documents, %s", err)
	}

	for _, d := range ds {
		var t T
		if err = d.DataTo(&t); err != nil {
			// should I return an error here? or just log it?
			return fmt.Errorf("error decoding document, %s", err)
		}
		*decoded = append(*decoded, t)
	}

	return nil
}

func set[T any](ctx context.Context, client *firestore.Client, collection string, document string, data T) error {
	if _, err := client.Collection(collection).Doc(document).Set(ctx, data); err != nil {
		return fmt.Errorf("error setting document contents, %s", err)
	}

	return nil
}

func update[T any](ctx context.Context, client *firestore.Client, collection string, document string, fields map[string]T) error {
	updates := make([]firestore.Update, 0)
	for k, v := range fields {
		updates = append(updates, firestore.Update{Path: k, Value: v})
	}

	if _, err := client.Collection(collection).Doc(document).Update(ctx, updates); err != nil {
		return fmt.Errorf("error updating document, %s", err)
	}

	return nil
}

func remove(ctx context.Context, client *firestore.Client, collection string, document string) error {
	if _, err := client.Collection(collection).Doc(document).Delete(ctx); err != nil {
		return fmt.Errorf("error deleting document, %s", err)
	}

	return nil
}
