package main

import (
	"fmt"
	"os"
	"strings"
)

func (m *model) OutputBodyToFile(path string, strip bool) error {

	output := strings.Join(m.body, "\n")
	if strip {
		output = stripANSI(output)
	}
	if err := os.WriteFile(path, []byte(output), 0744); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

func (m *model) OutputBodyToStdout(strip bool) error {

	output := strings.Join(m.body, "\n")
	if strip {
		output = stripANSI(output)
	}
	_, err := os.Stdout.Write([]byte(output))
	if err != nil {
		return fmt.Errorf("error writing to stdout: %w", err)
	}
	return nil
}
