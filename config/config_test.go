package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const testBasePath = "../testdata"

func TestUnmarshalConfig(t *testing.T) {
	testConfig := &Config{
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
			Quit:                     "q, ctrl+c",
			Up:                       "k, up",
			Down:                     "j, down",
			UpFocus:                  "alt+k",
			DownFocus:                "alt+j",
			HalfUp:                   "ctrl+u",
			HalfDown:                 "ctrl+d",
			FullUp:                   "ctrl+b",
			FullDown:                 "ctrl+f",
			GoToFirstLine:            "g",
			GoToLastLine:             "G",
			GoToTop:                  "H",
			GoToMiddle:               "M",
			GoToBottom:               "L",
			Search:                   "/",
			ClearSearch:              "alt+d",
			Normal:                   "esc",
			CursorWordForward:        "alt+right, alt+f",
			CursorWordBackward:       "alt+left, alt+b",
			CursorDeleteWordBackward: "alt+backspace",
			CursorDeleteWordForward:  "alt+delete",
			CursorDeleteAfterCursor:  "alt+k",
			CursorDeleteBeforeCursor: "alt+u",
			CursorLineStart:          "home, ctrl+a",
			CursorLineEnd:            "end, ctrl+e",
			CursorPaste:              "ctrl+v",
		},
	}

	configFileTests := []struct {
		name string
		file string
		want *Config
	}{
		{"full config yaml", "testConfig.yml", testConfig},
		{"full config json", "testConfig.json", testConfig},
		{"minimal config yaml", "testConfigMinimal.yml", newDefaultConfig(testBasePath)},
		{"minimal config json", "testConfigMinimal.json", newDefaultConfig(testBasePath)},
		{"config file absent", "testConfigAbsent.yml", newDefaultConfig(testBasePath)},
	}

	for _, tt := range configFileTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalConfig(filepath.Join(testBasePath, tt.file), testBasePath)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("empty config file path", func(t *testing.T) {
		want := newDefaultConfig(testBasePath)
		got, err := UnmarshalConfig("", testBasePath)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestUnmarshalKeyb(t *testing.T) {
	apps := Apps{{
		Name: "test",
		Keybinds: []KeyBind{{
			Name: "foo",
			Key:  "bar",
		}},
	}}

	keybFileTests := []struct {
		name string
		file string
		want Apps
	}{
		{"keyb file yaml", "testkeyb.yml", apps},
		{"keyb file json", "testkeyb.json", apps},
	}

	for _, tt := range keybFileTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalKeyb(filepath.Join(testBasePath, tt.file), testBasePath)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("file absent", func(t *testing.T) {
		_, err := UnmarshalKeyb(filepath.Join(testBasePath, "temp.yml"), testBasePath)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		t.Cleanup(func() {
			err := os.Remove(filepath.Join(testBasePath, "temp.yml"))
			if err != nil {
				t.Fatal(err)
			}
		})
	})

	t.Run("empty filepath", func(t *testing.T) {
		_, err := UnmarshalKeyb("", testBasePath)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		t.Cleanup(func() {
			err := os.Remove(filepath.Join(testBasePath, defaultKeybFile))
			if err != nil {
				t.Fatal(err)
			}
		})
	})
}
