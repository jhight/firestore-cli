package store

import (
	actions2 "github.com/jhight/firestore-cli/pkg/api/actions"
	"github.com/jhight/firestore-cli/pkg/api/store"
	"github.com/jhight/firestore-cli/pkg/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGetAction(t *testing.T) {
	gc := gomock.NewController(t)
	mockStore := store.NewMockStore(gc)

	root := actions2.Root(actions2.DefaultsInitializer(config.Config{}, mockStore))
	root.Add(actions2.Get(root))
	root.SetArgs([]string{"get", "example-collection", "example-document"})

	example := map[string]any{
		"foo": "bar",
	}

	mockStore.EXPECT().Get("example-collection", "example-document").Return(example, nil)

	err := root.Execute()
	assert.Nil(t, err)
}
