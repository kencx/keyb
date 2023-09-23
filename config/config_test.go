package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

const testBasePath = "../testdata"

func TestUnmarshalConfig(t *testing.T) {
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
		got, err := UnmarshalConfig(filepath.Join(testBasePath, "testConfig.yml"), testBasePath)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("minimal config", func(t *testing.T) {
		want := newDefaultConfig(testBasePath)

		got, err := UnmarshalConfig(filepath.Join(testBasePath, "testConfigMinimal.yml"), testBasePath)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

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

	t.Run("config file absent", func(t *testing.T) {
		want := newDefaultConfig(testBasePath)

		got, err := UnmarshalConfig(filepath.Join(testBasePath, "testConfigAbsent.yml"), testBasePath)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestUnmarshalKeyb(t *testing.T) {
	t.Run("file present", func(t *testing.T) {
		want := Apps{{
			Name: "test",
			Keybinds: []KeyBind{{
				Name: "foo",
				Key:  "bar",
			}},
		}}

		got, err := UnmarshalKeyb(filepath.Join(testBasePath, "testkeyb.yml"), testBasePath)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

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
