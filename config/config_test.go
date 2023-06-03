package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const (
	testDirPath     = "../testdata"
	testConfigDir = keybDirPath
)

func TestGetConfigDir(t *testing.T) {
	want := filepath.Join(testDirPath, testConfigDir)

	got, err := getConfigDir(testDirPath)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(filepath.Join(testDirPath, testConfigDir))
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

func TestReadConfigFile(t *testing.T) {
	t.Run("full config", func(t *testing.T) {
		want := &Config{
			Settings: Settings{
				KeybPath:       "./custom.yml",
				Debug:          true,
				Reverse:        true,
				Mouse:          false,
				SearchMode:     false,
				SortKeys:       true,
				Title:          "",
				Prompt:         "keys > ",
				PromptLocation: "bottom",
				Placeholder:    "...",
				PrefixSep:      ";",
				SepWidth:       4,
				Margin:         1,
				Padding:        1,
				BorderStyle:    "normal",
			},
			Color: Color{
				FilterFg: "#FFA066",
			},
			Keys: Keys{
				Quit:          "q, ctrl+c",
				Up:            "k, up",
				Down:          "j, down",
				UpFocus:       "alt+k",
				DownFocus:     "alt+j",
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
		got, err := ReadConfigFile(filepath.Join(testDirPath, "testConfig.yml"))
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("minimal config", func(t *testing.T) {
		want, err := newDefaultConfig()
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		got, err := ReadConfigFile(filepath.Join(testDirPath, "testConfigMinimal.yml"))
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestGetBaseDirWithSetEnvVar(t *testing.T) {
	want := "/foo/bar/.config"
	t.Setenv("XDG_CONFIG_HOME", want)

	got, err := getBaseDir()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
