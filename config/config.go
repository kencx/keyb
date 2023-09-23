package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	defaultConfigDir  = "keyb"
	defaultConfigFile = "config.yml"
	defaultKeybFile   = "keyb.yml"
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

// Read configuration and keyb file from flags, default path.
// If config directory and/or files do not exist, create them.
func Parse(flagCPath, flagKPath string) (Apps, *Config, error) {
	xdgConfigDir, err := getXDGConfigDir()
	if err != nil {
		return nil, nil, err
	}

	basePath := filepath.Join(xdgConfigDir, defaultConfigDir)
	err = os.MkdirAll(basePath, 0744)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create config dir: %w", err)
	}

	config, err := UnmarshalConfig(flagCPath, basePath)
	if err != nil {
		return nil, nil, err
	}

	if flagKPath == "" {
		flagKPath = config.KeybPath
	}

	keys, err := UnmarshalKeyb(flagKPath, basePath)
	if err != nil {
		return nil, nil, err
	}
	return keys, config, nil
}

// Read config file and merge with default config
func UnmarshalConfig(configFile, basePath string) (*Config, error) {

	// set default config filepath
	if configFile == "" {
		configFile = filepath.Join(basePath, defaultConfigFile)
	}

	res := newDefaultConfig(basePath)
	configFile = os.ExpandEnv(configFile)

	file, err := os.ReadFile(configFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return res, nil
		} else {
			return nil, fmt.Errorf("failed to read config file \"%s\": %w", configFile, err)
		}
	}

	if err = yaml.Unmarshal(file, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file \"%s\": %w", configFile, err)
	}
	return res, nil
}

func newDefaultConfig(basePath string) *Config {
	res := DefaultConfig
	res.KeybPath = filepath.Join(basePath, defaultKeybFile)
	return res
}

// Read keyb file or create default keyb file not exist
func UnmarshalKeyb(keybFile, basePath string) (Apps, error) {
	if keybFile == "" {
		keybFile = filepath.Join(basePath, defaultKeybFile)
	}

	keybFile = os.ExpandEnv(keybFile)
	file, err := os.ReadFile(keybFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {

			k := newDefaultKeyb(keybFile)
			data, err := yaml.Marshal(k)
			if err != nil {
				return nil, fmt.Errorf("failed to generate default keyb: %w", err)
			}

			if err := os.WriteFile(keybFile, data, 0644); err != nil {
				return nil, fmt.Errorf("failed to create keyb file: %w", err)
			}
			return k, nil

		} else {
			return nil, fmt.Errorf("failed to read keyb file: %w", err)
		}
	}

	var b Apps
	if err := yaml.Unmarshal(file, &b); err != nil {
		return nil, fmt.Errorf("failed to unmarshal keyb file: %w", err)
	}
	return b, nil
}

func newDefaultKeyb(path string) Apps {
	return Apps{{
		Name: "example",
		Keybinds: []KeyBind{{
			Name: "add your keys in",
			Key:  path,
		}},
	}}
}

// get user XDG_CONFIG_HOME directory
func getXDGConfigDir() (string, error) {
	val, ok := os.LookupEnv("XDG_CONFIG_HOME")
	if ok {
		return val, nil
	}

	path, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("user config directory not found: %w", err)
	}
	return path, nil
}
