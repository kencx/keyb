package output

import (
	"fmt"
	"os"

	"github.com/kencx/keyb/ui"
)

func ToFile(m *ui.Model, path string) error {
	output := m.List.UnstyledString()

	path = os.ExpandEnv(path)
	if err := os.WriteFile(path, []byte(output), 0664); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func ToStdout(m *ui.Model) error {
	output := m.List.UnstyledString()

	_, err := os.Stdout.Write([]byte(output))
	if err != nil {
		return fmt.Errorf("failed to write to stdout: %w", err)
	}
	return nil
}
