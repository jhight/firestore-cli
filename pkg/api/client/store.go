package client

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/jhight/firestore-cli/pkg/api/client/query"
	"github.com/jhight/firestore-cli/pkg/config"
	"google.golang.org/api/option"
	"os"
	"strings"
)

//go:generate go run go.uber.org/mock/mockgen -typed -package $GOPACKAGE -source $GOFILE -destination $GOFILE.mocks.go
type Store interface {
	IsPathToDocument(path string) bool
	IsPathToCollection(path string) bool
	Get(input query.Input) (map[string]any, error)
	Query(input query.Input) ([]map[string]any, error)
	Collections(input query.Input) ([]any, error)
	Create(path string, fields map[string]any) error
	Set(path string, fields map[string]any) error
	Update(path string, fields map[string]any) error
	Delete(path string) error
	DeleteField(path string, field string) error
	Close() error
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
