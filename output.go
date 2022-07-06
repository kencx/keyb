package main

import (
	"fmt"
	"os"
)

func (m *model) OutputBodyToFile(path string, strip bool) error {

	output := m.list.String()
	if strip {
		output = stripANSI(output)
	}

	path = os.ExpandEnv(path)
	if err := os.WriteFile(path, []byte(output), 0664); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func (m *model) OutputBodyToStdout(strip bool) error {

	output := m.list.String()
	if strip {
		output = stripANSI(output)
	}
	_, err := os.Stdout.Write([]byte(output))
	if err != nil {
		return fmt.Errorf("failed to write to stdout: %w", err)
	}
	return nil
}
