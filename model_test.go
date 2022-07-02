package main

import (
	"reflect"
	"testing"
)

var testData = Bindings{
	"tmux": {Prefix: "ctrl + b",
		Keybinds: []KeyBind{
			{Comment: "close window", Key: "shift + x"},
		}},
	"vim": {Keybinds: []KeyBind{
		{Comment: "focus left", Key: "ctrl + h"},
		{Comment: "swap left", Key: "ctrl + shift + h"},
	}},
	"firefox": {Prefix: "test",
		Keybinds: []KeyBind{
			{Comment: "incognito", Key: "ctrl + shift + p", IgnorePrefix: true},
			{Comment: "new tab", Key: "ctrl + shift + t", IgnorePrefix: true},
			{Comment: "bookmarks bar", Key: "ctrl + b", IgnorePrefix: true},
		}},
}

func TestSortKeys(t *testing.T) {
	got := testData.sortedKeys()
	want := []string{"firefox", "tmux", "vim"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
