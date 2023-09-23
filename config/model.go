package config

import (
	"fmt"
	"os"
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

func AddEntry(path, binding string, kbIgnorePrefix bool) error {
	xdgConfigDir, err := getXDGConfigDir()
	if err != nil {
		return err
	}

	// load existing struct from filepath
	apps, err := UnmarshalKeyb(path, xdgConfigDir)
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
