package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	keybDir        = "keyb"
	configFileName = "config.yml"
	keybFileName   = "keyb.yml"
)

type Config struct {
	Settings `yaml:"settings"`
	Color    `yaml:"color"`
}

type Settings struct {
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
	BorderStyle string `yaml:"border"`
}

type Color struct {
	PromptColor string `yaml:"prompt"`
	CursorFg    string `yaml:"cursor_fg"`
	CursorBg    string `yaml:"cursor_bg"`
	FilterFg    string `yaml:"filter_fg"`
	FilterBg    string `yaml:"filter_bg"`
	BorderColor string `yaml:"border_color"`
}

func Parse(flagKPath, cfgPath string) (Apps, *Config, error) {
	if err := CreateDefaultConfigFile(); err != nil {
		return nil, nil, fmt.Errorf("no config file found: %w", err)
	}

	cfg, err := ParseConfig(cfgPath)
	if err != nil {
		return nil, nil, err
	}

	// priority: flag > file
	var kPath string
	if flagKPath != "" {
		kPath = flagKPath
	}

	// set default path and create if absent
	if kPath == "" {
		kPath = cfg.KeybPath
		if !fileExists(kPath) {
			if err := writeDefaultKeybFile(); err != nil {
				return nil, nil, err
			}
		}
	}

	keys, err := ParseApps(kPath)
	if err != nil {
		return nil, nil, err
	}
	return keys, cfg, nil
}

// Create default config file if does not exist
func CreateDefaultConfigFile() error {
	baseDir, err := GetBaseDir()
	if err != nil {
		return err
	}

	configDir, err := GetorCreateConfigDir(baseDir)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(configDir, configFileName)
	if !fileExists(fullPath) {
		defaultConfig, err := generateDefaultConfig()
		if err != nil {
			return err
		}
		data, err := yaml.Marshal(defaultConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal default config: %w", err)
		}
		if err := os.WriteFile(fullPath, data, 0644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}
	return nil
}

func ParseConfig(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("no config path given")
	}
	path = os.ExpandEnv(path)

	if !fileExists(path) {
		return nil, fmt.Errorf("%s does not exist", path)
	}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	c, err := generateDefaultConfig()
	if err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(file, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}
	return c, nil
}

func generateDefaultConfig() (*Config, error) {
	keybPath, err := getDefaultKeybFilePath()
	if err != nil {
		return nil, err
	}

	return &Config{
		Settings: Settings{
			KeybPath:    keybPath,
			Debug:       false,
			Reverse:     false,
			Mouse:       true,
			Title:       "",
			Prompt:      "keys > ",
			Placeholder: "...",
			PrefixSep:   ";",
			SepWidth:    4,
			Margin:      0,
			Padding:     1,
			BorderStyle: "hidden",
		},
		Color: Color{
			FilterFg: "#FFA066",
		},
	}, nil
}

func GetorCreateConfigDir(baseDir string) (string, error) {
	if baseDir == "" {
		return "", fmt.Errorf("base config directory not found")
	}
	configPath := filepath.Join(baseDir, keybDir)

	if !fileExists(configPath) {
		err := os.MkdirAll(configPath, 0744)
		if err != nil {
			return "", fmt.Errorf("failed to create config dir: %w", err)
		}
	}
	return configPath, nil
}

func GetBaseDir() (string, error) {
	path, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("base config directory not found: %v", err)
	}
	return path, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
