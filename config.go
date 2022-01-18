package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

const (
	parentDir  = "keyb"
	configFile = "config"

	// TODO reformat templ
	configTempl = `TITLE = "All your keybindings in one place"
VIM = true
KEYB_PATH = "$HOME/.config/keyb/keyb.yml"
` +
		"CURSOR_FG = `#edb`\n" +
		"CURSOR_BG = `#448448`\n" +
		`BORDER = true
BORDER_COLOR = ""
`
)

type Config struct {
	Title    string
	Vim      bool
	KeybPath string

	Cursor_fg    string
	Cursor_bg    string
	Border       bool
	Border_color string
}

type Program struct {
	Prefix   string
	KeyBinds []KeyBind
}

type KeyBind struct {
	Desc    string
	Key     string
	Comment string
}

func GetConfig(configFile string) (*Config, error) {

	options := ini.LoadOptions{
		SkipUnrecognizableLines: true,
		AllowBooleanKeys:        true,
	}
	config, err := ini.LoadSources(options, os.ExpandEnv(configFile))
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	cfgSection := config.Section("")

	// parse booleans
	vim, err := cfgSection.Key("VIM").Bool()
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	border, err := cfgSection.Key("BORDER").Bool()
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	cfg := &Config{
		Title:    cfgSection.Key("TITLE").String(),
		Vim:      vim,
		KeybPath: cfgSection.Key("KEYB_PATH").String(),

		Cursor_fg:    cfgSection.Key("CURSOR_FG").String(),
		Cursor_bg:    cfgSection.Key("CURSOR_BG").String(),
		Border:       border,
		Border_color: cfgSection.Key("BORDER_COLOR").String(),
	}
	return cfg, nil
}

func GetPrograms(keybFile string) (map[string]Program, error) {

	var programs map[string]Program

	file, err := ioutil.ReadFile(os.ExpandEnv(keybFile))
	if err != nil {
		return nil, fmt.Errorf("error reading keyb.yaml: %w", err)
	}

	if err := yaml.Unmarshal(file, &programs); err != nil {
		return nil, fmt.Errorf("error unmarshalling keyb.yaml: %w", err)
	}
	return programs, nil
}

func getBaseDir() (string, error) {
	var err error
	var path string

	switch runtime.GOOS {
	case "windows":
		path = os.Getenv("APPDATA")
	case "linux":
		path = os.Getenv("XDG_CONFIG_HOME")
		if path == "" {
			path = filepath.Join(os.Getenv("HOME"), ".config")
		}
	default:
		err = fmt.Errorf("error: unsupported platform")
	}

	if path == "" {
		return "", fmt.Errorf("error: base directory not found")
	}
	return path, err
}

func createConfigDir() (string, error) {

	baseDir, err := getBaseDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(baseDir, parentDir)

	_, err = os.Stat(configPath)
	if err != nil {

		if errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(configPath, 0755)
			if err != nil {
				return "", fmt.Errorf("error creating config dir: %w", err)
			}
		} else {
			return "", fmt.Errorf("error determining file structure: %w", err)
		}
	}
	return configPath, nil
}

func createConfigFile() error {

	basePath, err := createConfigDir()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(basePath, configFile)

	_, err = os.Stat(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.WriteFile(fullPath, []byte(configTempl), 0755); err != nil {
				return fmt.Errorf("error writing config file: %w", err)
			}
		} else {
			return fmt.Errorf("error determining file structure: %w", err)
		}
	}

	return nil
}
