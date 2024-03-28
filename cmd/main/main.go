package main

import (
	"github.com/jhight/firestore-cli/pkg/api/actions"
	"os"
)

func main() {
	root := actions.Root(nil)

	root.Add(
		actions.Collections(root),
		actions.Get(root),
		actions.Update(root),
		actions.Set(root),
		actions.Create(root),
		actions.Delete(root),
	)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
