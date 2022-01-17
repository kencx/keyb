package main

import (
	"fmt"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type Program struct {
	Prefix   string    `mapstructure:"prefix,omitempty"`
	KeyBinds []KeyBind `mapstructure:"keybinds"`
}

type KeyBind struct {
	Command string
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
