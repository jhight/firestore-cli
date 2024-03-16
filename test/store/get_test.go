package store

import (
	"github.com/jhight/firestore-cli/pkg/api/actions"
	"github.com/jhight/firestore-cli/pkg/api/client"
	"github.com/jhight/firestore-cli/pkg/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGetAction(t *testing.T) {
	gc := gomock.NewController(t)
	mockStore := client.NewMockStore(gc)

	root := actions.Root(actions.DefaultsInitializer(config.Config{}, mockStore))
	root.Add(actions.Get(root))
	root.SetArgs([]string{"get", "example-path"})

	example := map[string]any{
		"foo": "bar",
	}

	mockStore.EXPECT().Get("example-path").Return(example, nil)

	err := root.Execute()
	assert.Nil(t, err)
}
