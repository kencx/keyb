package main

import (
	"reflect"
	"testing"
)

var testData = map[string]App{
	"tmux": {prefix: "ctrl + b",
		keybinds: []KeyBind{
			{comment: "close window", key: "shift + x"},
		}},
	"vim": {keybinds: []KeyBind{
		{comment: "focus left", key: "ctrl + h"},
		{comment: "swap left", key: "ctrl + shift + h"},
	}},
	"firefox": {prefix: "test",
		keybinds: []KeyBind{
			{comment: "incognito", key: "ctrl + shift + p", ignorePrefix: true},
			{comment: "new tab", key: "ctrl + shift + t", ignorePrefix: true},
			{comment: "bookmarks bar", key: "ctrl + b", ignorePrefix: true},
		}},
}

func TestSortKeys(t *testing.T) {
	got := sortKeys(testData)
	want := []string{"firefox", "tmux", "vim"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSplitHeadingsAndKeys(t *testing.T) {

	m := &model{
		categories: testData,
		headings:   []string{"firefox", "tmux", "vim"},
		maxWidth:   88,
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
		categories: testData,
		headings:   []string{"firefox", "tmux", "vim"},
		maxWidth:   88,
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
