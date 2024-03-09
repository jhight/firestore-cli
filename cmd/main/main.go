package main

import (
	"github.com/jhight/firestore-cli/cmd/main/commands"
	"os"
)

func main() {
	root := commands.NewRootCommand()

	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
