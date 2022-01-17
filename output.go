package main

import (
	"fmt"
	"os"
	"strings"
)

func (m *model) outputBodyToFile(path string) error {

	output := strings.Join(m.body, "\n")
	if err := os.WriteFile(path, []byte(output), 0744); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

func (m *model) outputBodyToStdout() error {

	output := strings.Join(m.body, "\n")
	_, err := os.Stdout.Write([]byte(output))
	if err != nil {
		return fmt.Errorf("error writing to stdout: %w", err)
	}
	return nil
}
