package commands

import (
	"github.com/spf13/cobra"
	"io"
	"os"
)

func shouldReadFromStdin(cmd *cobra.Command) bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func readFromStdin(cmd *cobra.Command) (string, error) {
	input, err := io.ReadAll(cmd.InOrStdin())
	return string(input), err
}
