package client

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/jhight/firestore-cli/pkg/config"
	"google.golang.org/api/option"
	"os"
	"strings"
)

//go:generate go run go.uber.org/mock/mockgen -typed -package $GOPACKAGE -source $GOFILE -destination $GOFILE.mocks.go
type Store interface {
	Count(collection string) int
	Query(input SelectionInput, filter string) ([]map[string]any, error)
	List(input SelectionInput, fieldPath string) ([]any, error)
	Create(path string, fields map[string]any) error
	Get(path string) (map[string]any, error)
	Set(path string, fields map[string]any) error
	Update(path string, fields map[string]any) error
	Delete(path string) error
	DeleteField(path string, field string) error
	Close() error
}

type SelectionInput struct {
	CollectionPath string
	OrderBy        []OrderBy
	Limit          int
	Offset         int
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
