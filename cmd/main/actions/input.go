package actions

import (
	"io"
	"os"
)

func (a *Action) shouldReadFromStdin() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func (a *Action) readFromStdin() (string, error) {
	input, err := io.ReadAll(a.command.InOrStdin())
	return string(input), err
}
