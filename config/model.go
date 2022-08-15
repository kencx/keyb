package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

type App struct {
	Name     string    `yaml:"name"`
	Prefix   string    `yaml:"prefix,omitempty"`
	Keybinds []KeyBind `yaml:"keybinds"`
}

type Apps []*App

func (a App) String() string {
	return fmt.Sprintf("App{name=%s,prefix=%s,keybinds=%v}", a.Name, a.Prefix, a.Keybinds)
}

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

func AddEntry(path, binding string, kbIgnorePrefix bool) error {

	// load existing struct from filepath
	apps, err := ParseApps(path)
	if err != nil {
		return err
	}

	if binding == "" {
		return fmt.Errorf("binding must be given in format [app; name; keybind]")
	}

	s := strings.Split(binding, ";")
	if len(s) < 3 {
		return fmt.Errorf("binding must be given in format [app; name; keybind]")
	}
	input := struct {
		AppName string
		Name    string
		Key     string
	}{
		AppName: strings.Trim(s[0], " "),
		Name:    strings.Trim(s[1], " "),
		Key:     strings.Trim(s[2], " "),
	}

	apps.addOrUpdate(input.AppName, input.Name, input.Key, kbIgnorePrefix)

	// rewrite file
	data, err := yaml.Marshal(apps)
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}

	path = os.ExpandEnv(path)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write keyb file: %w", err)
	}

	return nil
}

func (apps *Apps) addOrUpdate(appName string, name, key string, ignorePrefix bool) {
	newKeyBind := KeyBind{
		Name:         name,
		Key:          key,
		IgnorePrefix: ignorePrefix,
	}

	if !apps.exist(appName) {
		a := App{
			Name:     appName,
			Keybinds: []KeyBind{newKeyBind},
		}
		*apps = append(*apps, &a)

	} else {
		for _, app := range *apps {
			if appName == app.Name {
				app.Keybinds = append(app.Keybinds, newKeyBind)
			}
		}
	}
	// return apps
}

func (apps Apps) exist(appName string) bool {
	for _, app := range apps {
		if appName == app.Name {
			return true
		}
	}
	return false
}
