package main

import (
	"fmt"
	"os"

	"github.com/kencx/keyb/ui"
)

func OutputBodyToFile(m *ui.Model, path string, strip bool) error {

	output := m.List.UnstyledString()
	if strip {
		output = stripANSI(output)
	}

	path = os.ExpandEnv(path)
	if err := os.WriteFile(path, []byte(output), 0664); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func OutputBodyToStdout(m *ui.Model, strip bool) error {

	output := m.List.UnstyledString()
	if strip {
		output = stripANSI(output)
	}
	_, err := os.Stdout.Write([]byte(output))
	if err != nil {
		return fmt.Errorf("failed to write to stdout: %w", err)
	}
	return nil
}
