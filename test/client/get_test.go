package client

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"jhight.com/firestore-cli/pkg/api/actions"
	"jhight.com/firestore-cli/pkg/api/client"
	"jhight.com/firestore-cli/pkg/api/client/query"
	"jhight.com/firestore-cli/pkg/config"
	"testing"
)

func TestGetAction(t *testing.T) {
	t.Skip("TODO")

	gc := gomock.NewController(t)
	mockStore := client.NewMockStore(gc)

	root := actions.Root(actions.DefaultsInitializer(config.Config{}, mockStore))
	root.Add(actions.Get(root))
	root.SetArgs([]string{"get", "example-path"})

	example := []map[string]any{
		{
			"foo": "bar",
		},
	}

	mockStore.EXPECT().IsPathToCollection("example-path").Return(true)
	mockStore.EXPECT().Query(query.Input{Path: "example-path"}).Return(example, nil)

	err := root.Execute()
	assert.Nil(t, err)
}
