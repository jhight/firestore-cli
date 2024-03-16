package actions

import (
	"github.com/jhight/firestore-cli/pkg/api/actions"
	"github.com/jhight/firestore-cli/pkg/api/client"
	"github.com/jhight/firestore-cli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGetAction(t *testing.T) {
	gc := gomock.NewController(t)
	mockStore := client.NewMockStore(gc)
	mockGetAction := actions.NewMockAction(gc)

	mockGetAction.EXPECT().Command().Return(&cobra.Command{Use: "get", Short: "Get a document from a collection."})

	root := actions.Root(actions.DefaultsInitializer(config.Config{}, mockStore))
	root.Add(mockGetAction)
	root.SetArgs([]string{"get", "example-path"})

	err := root.Execute()
	assert.Nil(t, err)
}
