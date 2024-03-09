package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/jhight/firestore-cli/pkg/config"
	"google.golang.org/api/option"
	"os"
	"strings"
)

type Store interface {
	Count(collection string) int
	Query(input SelectionInput, filter string) ([]map[string]any, error)
	List(input SelectionInput, path string) ([]any, error)
	Create(collection string, document string, fields map[string]any) error
	Get(collection string, document string) (map[string]any, error)
	Set(collection string, document string, fields map[string]any) error
	Update(collection string, document string, fields map[string]any) error
	Delete(collection string, document string) error
	Close() error
}

type SelectionInput struct {
	Collection string
	OrderBy    []OrderBy
	Limit      int
	Offset     int
}

type OrderBy struct {
	Field     string
	Direction Direction
}

func New(ctx context.Context, cfg config.Config) (Store, error) {
	home := os.Getenv("HOME")
	path := strings.ReplaceAll(cfg.ServiceAccount, "~", home)

	client, err := firestore.NewClient(ctx, cfg.ProjectID, option.WithCredentialsFile(path))
	if err != nil {
		return nil, fmt.Errorf("error creating firestore client, %s", err)
	}

	return &firestoreClientManager{
		ctx:    ctx,
		client: client,
	}, nil
}
