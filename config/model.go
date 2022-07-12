package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

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

	path = os.ExpandEnv(path)
	if !fileExists(path) {
		return nil, fmt.Errorf("%s does not exist", path)
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read keyb file: %w", err)
	}

	var b Apps
	if err := yaml.Unmarshal(file, &b); err != nil {
		return nil, fmt.Errorf("failed to unmarshal keyb file: %w", err)
	}
	return b, nil
}

func writeDefaultKeybFile() error {
	path, err := getDefaultKeybFilePath()
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(generateDefaultKeyb(path))
	if err != nil {
		return fmt.Errorf("failed to generate default keyb: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to create keyb file: %w", err)
	}
	return nil
}

func getDefaultKeybFilePath() (string, error) {
	baseDir, err := GetBaseDir()
	if err != nil {
		return "", err
	}

	return path.Join(baseDir, keybDir, keybFileName), nil
}

func generateDefaultKeyb(path string) Apps {
	return Apps{{
		Name: "example",
		Keybinds: []KeyBind{{
			Name: "add your keys in",
			Key:  path,
		}},
	}}
}
