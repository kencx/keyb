package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

//go:embed examples/config
var configTempl string

const (
	parentDir  = "keyb"
	configFile = "config"
)

type Config struct {
	Title    string
	Vim      bool
	KeybPath string

	Cursor_fg    string
	Cursor_bg    string
	Border       string
	Border_color string
}

func GetConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, fmt.Errorf("no config path given")
	}

	options := ini.LoadOptions{
		SkipUnrecognizableLines: true,
		AllowBooleanKeys:        true,
	}
	config, err := ini.LoadSources(options, os.ExpandEnv(configPath))
	if err != nil {
		return nil, fmt.Errorf("failed to load config path: %w", err)
	}

	cfgSection := config.Section("")

	// parse boolean
	vim, err := cfgSection.Key("VIM").Bool()
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	cfg := &Config{
		Title:    cfgSection.Key("TITLE").String(),
		Vim:      vim,
		KeybPath: cfgSection.Key("KEYB_PATH").String(),

		Cursor_fg:    cfgSection.Key("CURSOR_FG").String(),
		Cursor_bg:    cfgSection.Key("CURSOR_BG").String(),
		Border:       cfgSection.Key("BORDER").String(),
		Border_color: cfgSection.Key("BORDER_COLOR").String(),
	}
	return cfg, nil
}

func GetBindings(keybPath string) (Bindings, error) {

	var b Bindings

	file, err := ioutil.ReadFile(os.ExpandEnv(keybPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read keyb file: %w", err)
	}

	if err := yaml.Unmarshal(file, &b); err != nil {
		return nil, fmt.Errorf("failed to unmarshall keyb file: %w", err)
	}
	return b, nil
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
		err = fmt.Errorf("unsupported platform")
	}

	if path == "" {
		return "", fmt.Errorf("base config directory not found")
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
			err := os.MkdirAll(configPath, 0664)
			if err != nil {
				return "", fmt.Errorf("failed to create config dir: %w", err)
			}
		} else {
			return "", fmt.Errorf("failed to determine file structure: %w", err)
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
			if err := os.WriteFile(fullPath, []byte(configTempl), 0664); err != nil {
				return fmt.Errorf("failed to write config file: %w", err)
			}
		} else {
			return fmt.Errorf("failed to determine file structure: %w", err)
		}
	}
	return nil
}
