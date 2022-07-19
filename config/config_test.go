package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	parentDir     = "../testdata"
	testConfigDir = "keyb"
)

func TestCreateConfigDir(t *testing.T) {

	want := filepath.Join(parentDir, testConfigDir)
	got, err := GetorCreateConfigDir(parentDir)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(path.Join(parentDir, testConfigDir))
	})

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}

	info, err := os.Stat(got)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if info.Name() != testConfigDir {
		t.Errorf("got %v, want %v", info.Name(), testConfigDir)
	}

	if !info.IsDir() {
		t.Errorf("%v not a directory", info.Name())
	}

	if fmt.Sprintf("%#o", info.Mode().Perm()) != "0744" {
		t.Errorf("got %v, want %v", info.Mode().Perm(), "0744")
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
		Keys: Keys{
			Quit:          "q, ctrl+c",
			Up:            "k, up",
			Down:          "j, down",
			HalfUp:        "ctrl+u",
			HalfDown:      "ctrl+d",
			FullUp:        "ctrl+b",
			FullDown:      "ctrl+f",
			GoToFirstLine: "g",
			GoToLastLine:  "G",
			GoToTop:       "H",
			GoToMiddle:    "M",
			GoToBottom:    "L",
			Search:        "/",
			ClearSearch:   "alt+d",
			Normal:        "esc",
		},
	}
	got, err := ParseConfig(path.Join(parentDir, "testconfig.yml"))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
