package client

import (
	"github.com/jhight/firestore-cli/pkg/api/client/query"
)

func (f *firestoreClientManager) Collections(input query.Input) ([]any, error) {
	return collections(f.ctx, f.client, input.Path), nil
}
