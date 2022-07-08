package ui

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type App struct {
	Name     string    `yaml:"name"`
	Prefix   string    `yaml:"prefix,omitempty"`
	Keybinds []KeyBind `yaml:"keybinds"`
}

type Apps []App

type KeyBind struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`

	// ignore prefix defaults to false
	// so user can choose to ignore prefix for a specific kb
	IgnorePrefix bool `yaml:"ignore_prefix,omitempty"`
}

func ParseApps(path string) (Apps, error) {
	if path == "" {
		return nil, fmt.Errorf("no keyb path given")
	}

	file, err := ioutil.ReadFile(os.ExpandEnv(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read keyb file: %w", err)
	}

	var b Apps
	if err := yaml.Unmarshal(file, &b); err != nil {
		return nil, fmt.Errorf("failed to unmarshal keyb file: %w", err)
	}
	return b, nil
}
