package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kencx/keyb/ui"
	"gopkg.in/yaml.v2"
)

func ToFile(m *ui.Model, path string) error {
	var (
		output []byte
		err    error
	)

	path = os.ExpandEnv(path)
	ext := filepath.Ext(path)

	switch ext {
	case ".json":
		output, err = json.Marshal(m.Apps)
		if err != nil {
			return fmt.Errorf("failed to marshal to json: %w", err)
		}
	case ".yml", ".yaml":
		output, err = yaml.Marshal(m.Apps)
		if err != nil {
			return fmt.Errorf("failed to marshal to yaml: %w", err)
		}
	default:
		output = []byte(m.List.UnstyledString())
	}
	if err := os.WriteFile(path, output, 0664); err != nil {
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
