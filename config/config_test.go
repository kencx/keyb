package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

var (
	parentDir = "../testdata"
	testDir   = "keyb"
	fileName  = "config.yml"
)

func TestCreateConfigDir(t *testing.T) {

	want := filepath.Join(parentDir, testDir)
	got, err := GetorCreateConfigDir(parentDir)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(path.Join(parentDir, testDir))
	})

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	info, err := os.Stat(got)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if info.Name() != testDir {
		t.Errorf("got %v, want %v", info.Name(), testDir)
	}

	if !info.IsDir() {
		t.Errorf("%v not a directory", info.Name())
	}

	if fmt.Sprintf("%#o", info.Mode().Perm()) != "0744" {
		t.Errorf("got %v, want %v", info.Mode().Perm(), "0744")
	}
}

func TestCreateConfigFile(t *testing.T) {
	want := path.Join(parentDir, testDir, fileName)
	err := CreateConfigFile(parentDir)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(path.Join(parentDir, testDir))
	})

	data, err := os.ReadFile(want)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	defaultConfig, err := generateDefaultConfig()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	var c Config
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if reflect.DeepEqual(c, defaultConfig) {
		t.Errorf("got %v, want %v", c, defaultConfig)
	}

	info, err := os.Stat(want)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if info.Name() != fileName {
		t.Errorf("got %v, want %v", info.Name(), fileName)
	}

	if info.IsDir() {
		t.Errorf("%v is a directory", info.Name())
	}

	if fmt.Sprintf("%#o", info.Mode().Perm()) != "0644" {
		t.Errorf("got %v, want %v", info.Mode().Perm(), "0644")
	}
}

func TestParse(t *testing.T) {
	want := &Config{
		Settings: Settings{
			KeybPath:    "./custom.yml",
			Debug:       true,
			Reverse:     true,
			Mouse:       false,
			Title:       "",
			Prompt:      "keys > ",
			Placeholder: "...",
			PrefixSep:   ";",
			SepWidth:    4,
			Margin:      1,
			Padding:     1,
			BorderStyle: "normal",
		},
		Color: Color{
			FilterFg: "#FFA066",
		},
	}
	got, err := Parse(path.Join(parentDir, "testconfig.yml"))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
