package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

const parentdir = "keyb"
const configTempl = `TITLE="All your keybindings in one place"

CURSORFB="#edb"
CURSORBG="#448448"

KEYBRC_DIR="$HOME/.config/keyrc"
`

type Program struct {
	Prefix   string    `mapstructure:"prefix,omitempty"`
	KeyBinds []KeyBind `mapstructure:"keybinds"`
}

type KeyBind struct {
	Desc    string
	Key     string
	Comment string `mapstructure:"comment,omitempty"`
}

func GetConfig(keybrc string) (map[string]string, error) {

	var config map[string]string
	var k = koanf.New(".")

	if err := k.Load(file.Provider(keybrc), dotenv.Parser()); err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	k.Unmarshal("", &config)
	return config, nil
}

func GetPrograms(keyFile string) (map[string]Program, error) {

	var programs map[string]Program
	var k = koanf.New(".")

	if err := k.Load(file.Provider(keyFile), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading keymap: %w", err)
	}

	k.Unmarshal("", &programs)
	return programs, nil
}

func GetBaseDir() (string, error) {
	var err error
	var path string

	switch runtime.GOOS {
	case "windows":
		path = os.Getenv("APPDATA")
	case "linux":
		// check config for custom dir
		path = os.Getenv("XDG_DATA_HOME")
		if path == "" {
			homepath := os.Getenv("HOME")
			path = filepath.Join(homepath, ".config")
		}
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return path, err
}

func CreateConfigDir() (string, error) {

	basedir, err := GetBaseDir()
	if err != nil {
		return "", err
	}
	dirPath := filepath.Join(basedir, parentdir)

	_, err = os.Stat(dirPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(dirPath, 0755)
			if err != nil {
				return "", fmt.Errorf("error creating parent dir: %v", err)
			}
		} else {
			return "", fmt.Errorf("error determining file structure: %v", err)
		}
	}
	return dirPath, nil
}

// add default config file
func CreateConfigFile() error {

	basePath, err := CreateConfigDir()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(basePath, ".keybrc")
	if err := os.WriteFile(fullPath, []byte(configTempl), 0755); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}
	return nil
}
