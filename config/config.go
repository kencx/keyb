package config

import (
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var ConfigFs string

const (
	keybDir        = "keyb"
	configFileName = "config.yml"
	keybFileName   = "keyb.yml"
)

type Config struct {
	Defaults `yaml:"defaults"`
	Colour   `yaml:"color"`
}

type Defaults struct {
	KeybPath    string `yaml:"keyb_path"`
	Debug       bool
	Reverse     bool
	Mouse       bool
	Title       string
	Prompt      string
	Placeholder string
	PrefixSep   string `yaml:"prefix_sep"`
	SepWidth    int    `yaml:"sep_width"`
	Margin      int
	Padding     int
	BorderStyle string `yaml:"border_style"`
}

type Colour struct {
	PromptColor string `yaml:"prompt"`
	CursorFg    string `yaml:"cursor_fg"`
	CursorBg    string `yaml:"cursor_bg"`
	FilterFg    string `yaml:"filter_fg"`
	FilterBg    string `yaml:"filter_bg"`
	BorderColor string `yaml:"border_color"`
}

// TODO combine set config with defaults
func Parse(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("no config path given")
	}

	file, err := ioutil.ReadFile(os.ExpandEnv(path))
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var c Config
	if err = yaml.Unmarshal(file, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}
	return &c, nil
}

func getBaseDir() (string, error) {
	var err error

	path := os.Getenv("XDG_CONFIG_HOME")
	if path == "" {
		path, err = os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("base config directory not found: %v", err)
		}
	}
	return path, err
}

func createConfigDir() (string, error) {
	baseDir, err := getBaseDir()
	if err != nil {
		return "", err
	}

	if baseDir == "" {
		return "", fmt.Errorf("base config directory not found")
	}
	configPath := filepath.Join(baseDir, keybDir)

	_, err = os.Stat(configPath)
	if err != nil {

		if errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(configPath, 0744)
			if err != nil {
				return "", fmt.Errorf("failed to create config dir: %w", err)
			}
		} else {
			return "", fmt.Errorf("failed to determine file structure: %w", err)
		}
	}
	return configPath, nil
}

func CreateConfigFile() error {
	basePath, err := createConfigDir()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(basePath, configFileName)

	_, err = os.Stat(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.WriteFile(fullPath, []byte(ConfigFs), 0744); err != nil {
				return fmt.Errorf("failed to create config file: %w", err)
			}
		} else {
			return fmt.Errorf("failed to determine file structure: %w", err)
		}
	}
	return nil
}
