package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	keybDirPath    = "keyb"
	configFileName = "config.yml"
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
	Quit                     string
	Up                       string
	Down                     string
	UpFocus                  string `yaml:"up_focus"`
	DownFocus                string `yaml:"down_focus"`
	HalfUp                   string `yaml:"half_up"`
	HalfDown                 string `yaml:"half_down"`
	FullUp                   string `yaml:"full_up"`
	FullDown                 string `yaml:"full_bottom"`
	GoToFirstLine            string `yaml:"first_line"`
	GoToLastLine             string `yaml:"last_line"`
	GoToTop                  string `yaml:"top"`
	GoToMiddle               string `yaml:"middle"`
	GoToBottom               string `yaml:"bottom"`
	Search                   string
	ClearSearch              string `yaml:"clear_search"`
	Normal                   string
	CursorWordForward        string `yaml:"cursor_word_forward"`
	CursorWordBackward       string `yaml:"cursor_word_backward"`
	CursorDeleteWordBackward string `yaml:"cursor_delete_word_backward"`
	CursorDeleteWordForward  string `yaml:"cursor_delete_word_forward"`
	CursorDeleteAfterCursor  string `yaml:"cursor_delete_after_cursor"`
	CursorDeleteBeforeCursor string `yaml:"cursor_delete_before_cursor"`
	CursorLineStart          string `yaml:"cursor_line_start"`
	CursorLineEnd            string `yaml:"cursor_line_end"`
	CursorPaste              string `yaml:"cursor_paste"`
}

var DefaultConfig = &Config{
	Settings: Settings{
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
		Quit:                     "q, ctrl+c",
		Up:                       "k, up",
		Down:                     "j, down",
		UpFocus:                  "ctrl+k",
		DownFocus:                "ctrl+j",
		HalfUp:                   "ctrl+u",
		HalfDown:                 "ctrl+d",
		FullUp:                   "ctrl+b",
		FullDown:                 "ctrl+f",
		GoToFirstLine:            "g",
		GoToLastLine:             "G",
		GoToTop:                  "H",
		GoToMiddle:               "M",
		GoToBottom:               "L",
		Search:                   "/",
		ClearSearch:              "alt+d",
		Normal:                   "esc",
		CursorWordForward:        "alt+right, alt+f",
		CursorWordBackward:       "alt+left, alt+b",
		CursorDeleteWordBackward: "alt+backspace",
		CursorDeleteWordForward:  "alt+delete",
		CursorDeleteAfterCursor:  "alt+k",
		CursorDeleteBeforeCursor: "alt+u",
		CursorLineStart:          "home, ctrl+a",
		CursorLineEnd:            "end, ctrl+e",
		CursorPaste:              "ctrl+v",
	},
}

func Parse(flagKPath, configPath string) (Apps, *Config, error) {
	var (
		config *Config
		err    error
	)

	switch configPath {
	case "":
		config, err = ReadDefaultConfigFile()
	default:
		config, err = ReadConfigFile(configPath)
	}
	if err != nil {
		return nil, nil, err
	}

	// priority: flag > file
	var kPath string
	if flagKPath != "" {
		kPath = flagKPath
	}

	// If no keyb file present, create a default file and set it as kPath
	if kPath == "" {
		kPath = config.KeybPath
		if !pathExists(kPath) {
			if err := writeDefaultKeybFile(); err != nil {
				return nil, nil, err
			}
		}
	}

	keys, err := ReadKeybFile(kPath)
	if err != nil {
		return nil, nil, err
	}
	return keys, config, nil
}

// Read config file at default path if exist. Otherwise, return default config
func ReadDefaultConfigFile() (*Config, error) {
	baseDir, err := getBaseDir()
	if err != nil {
		return nil, err
	}
	configDir, err := getConfigDir(baseDir)
	if err != nil {
		return nil, err
	}

	var config *Config
	defaultConfigFilePath := filepath.Join(configDir, configFileName)
	if !pathExists(defaultConfigFilePath) {
		config, err = newDefaultConfig()
		if err != nil {
			return nil, err
		}
	} else {
		config, err = ReadConfigFile(defaultConfigFilePath)
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

// Read given config file and merge with default config
func ReadConfigFile(path string) (*Config, error) {
	path = os.ExpandEnv(path)
	if !pathExists(path) {
		return nil, fmt.Errorf("config file \"%s\" does not exist", path)
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file \"%s\": %w", path, err)
	}

	c, err := newDefaultConfig()
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(file, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file \"%s\": %w", path, err)
	}
	return c, nil
}

func newDefaultConfig() (*Config, error) {
	res := DefaultConfig

	baseDir, err := getBaseDir()
	if err != nil {
		return nil, err
	}

	res.KeybPath = filepath.Join(baseDir, keybDirPath, keybFileName)
	return res, nil
}

// Get or create keyb config directory
func getConfigDir(baseDir string) (string, error) {
	if baseDir == "" {
		return "", fmt.Errorf("base config directory not found")
	}
	configPath := filepath.Join(baseDir, keybDirPath)

	if !pathExists(configPath) {
		err := os.MkdirAll(configPath, 0744)
		if err != nil {
			return "", fmt.Errorf("failed to create config dir: %w", err)
		}
	}
	return configPath, nil
}

// Fetch OS-dependent user config directory
func getBaseDir() (string, error) {
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

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
