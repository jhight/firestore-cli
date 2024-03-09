package main

import (
	"github.com/jhight/firestore-cli/pkg/api/actions"
	"os"
)

func main() {
	root := actions.Root(nil)

	root.Add(
		actions.Get(root),
		actions.Count(root),
		actions.Create(root),
		actions.Update(root),
		actions.Set(root),
		actions.Delete(root),
		actions.Query(root),
		actions.List(root),
	)

	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
