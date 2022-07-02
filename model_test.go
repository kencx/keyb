package main

import (
	"reflect"
	"testing"
)

var testData = map[string]App{
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
	got := sortedKeys(testData)
	want := []string{"firefox", "tmux", "vim"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSplitHeadingsAndKeys(t *testing.T) {

	m := &model{
		BindStructure: BindStructure{
			binds:    testData,
			headings: []string{"firefox", "tmux", "vim"},
		},
		maxWidth: 88,
	}
	gotHeadings, gotLines := m.splitHeadingsAndKeys()

	wantHeadings := map[int]string{
		0: "firefox",
		4: "tmux",
		6: "vim",
	}
	wantLines := map[int]string{
		1: "incognito\tctrl + shift + p",
		2: "new tab\tctrl + shift + t",
		3: "bookmarks bar\tctrl + b",
		5: "close window\tctrl + b ; shift + x",
		7: "focus left\tctrl + h",
		8: "swap left\tctrl + shift + h",
	}
	wantLineCount := 9

	assertMapEqual(t, gotHeadings, wantHeadings)
	assertMapEqual(t, gotLines, wantLines)
	if m.lineCount != wantLineCount {
		t.Errorf("got %d, want %d", m.lineCount, wantLineCount)
	}
}

func BenchmarkSplitHeadingsAndKeys(b *testing.B) {
	m := &model{
		BindStructure: BindStructure{
			binds:    testData,
			headings: []string{"firefox", "tmux", "vim"},
		},
		maxWidth: 88,
	}
	for i := 0; i < b.N; i++ {
		m.splitHeadingsAndKeys()
	}
}

func assertMapEqual(t testing.TB, got, want map[int]string) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot %q\nwant %q", got, want)
	}
}
