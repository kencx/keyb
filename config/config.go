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
	Keys     `yaml:"keys"`
}

type Settings struct {
	KeybPath       string `yaml:"keyb_path"`
	Debug          bool
	Reverse        bool
	Mouse          bool
	SearchMode     bool `yaml:"search_mode"`
	SortKeys       bool `yaml:"sort_keys"`
	Title          string
	Prompt         string
	PromptLocation string `yaml:"prompt_location"`
	Placeholder    string
	PrefixSep      string `yaml:"prefix_sep"`
	SepWidth       int    `yaml:"sep_width"`
	Margin         int
	Padding        int
	BorderStyle    string `yaml:"border"`
}

type Color struct {
	PromptColor   string `yaml:"prompt"`
	CursorFg      string `yaml:"cursor_fg"`
	CursorBg      string `yaml:"cursor_bg"`
	FilterFg      string `yaml:"filter_fg"`
	FilterBg      string `yaml:"filter_bg"`
	CounterFg     string `yaml:"counter_fg"`
	CounterBg     string `yaml:"counter_bg"`
	PlaceholderFg string `yaml:"placeholder_fg"`
	PlaceholderBg string `yaml:"placeholder_bg"`
	BorderColor   string `yaml:"border_color"`
}

type Keys struct {
	Quit          string
	Up            string
	Down          string
	HalfUp        string `yaml:"half_up"`
	HalfDown      string `yaml:"half_down"`
	FullUp        string `yaml:"full_up"`
	FullDown      string `yaml:"full_bottom"`
	GoToFirstLine string `yaml:"first_line"`
	GoToLastLine  string `yaml:"last_line"`
	GoToTop       string `yaml:"top"`
	GoToMiddle    string `yaml:"middle"`
	GoToBottom    string `yaml:"bottom"`
	Search        string
	ClearSearch   string `yaml:"clear_search"`
	Normal        string
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
			KeybPath:       keybPath,
			Debug:          false,
			Reverse:        false,
			Mouse:          true,
			SearchMode:     false,
			SortKeys:       false,
			Title:          "",
			Prompt:         "keys > ",
			PromptLocation: "top",
			Placeholder:    "...",
			PrefixSep:      ";",
			SepWidth:       4,
			Margin:         0,
			Padding:        1,
			BorderStyle:    "hidden",
		},
		Color: Color{
			FilterFg: "#FFA066",
		},
		Keys: Keys{
			Quit:          "q, ctrl+c",
			Up:            "k, up",
			Down:          "j, down",
			HalfUp:        "ctrl+u",
			HalfDown:      "ctrl+d",
			FullUp:        "ctrl+b",
			FullDown:      "ctrl+f",
			GoToFirstLine: "g",
			GoToLastLine:  "G",
			GoToTop:       "H",
			GoToMiddle:    "M",
			GoToBottom:    "L",
			Search:        "/",
			ClearSearch:   "alt+d",
			Normal:        "esc",
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
	val, ok := os.LookupEnv("XDG_CONFIG_HOME")
	if ok {
		return val, nil
	}

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
