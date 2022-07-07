package config

import (
	"sort"
)

type Bindings map[string]App

func (b Bindings) SortedKeys() []string {
	keys := make([]string, 0, len(b))
	for k := range b {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// TODO replace bindings with []App
// name of App should be a field instead
type App struct {
	Prefix   string    `yaml:"prefix,omitempty"`
	Keybinds []KeyBind `yaml:"keybinds"`
}

type KeyBind struct {
	Comment string `yaml:"desc"`
	Key     string `yaml:"key"`

	// config prefix ignore defaults to false
	// so user can choose to ignore prefix for a specific kb
	IgnorePrefix bool `yaml:"ignore_prefix,omitempty"`
}
