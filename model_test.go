package main

import (
	"reflect"
	"testing"
)

var testData = map[string]Program{
	"tmux": {KeyBinds: []KeyBind{
		{Desc: "close window", Key: "shift + x"},
	}},
	"vim": {KeyBinds: []KeyBind{
		{Desc: "focus left", Key: "ctrl + h"},
		{Desc: "swap left", Key: "ctrl + shift + h"},
	}},
	"firefox": {KeyBinds: []KeyBind{
		{Desc: "incognito", Key: "ctrl + shift + p"},
		{Desc: "new tab", Key: "ctrl + shift + t"},
		{Desc: "bookmarks bar", Key: "ctrl + b"},
	}},
}

func TestSortKeys(t *testing.T) {
	got := sortKeys(testData)
	want := []string{"firefox", "tmux", "vim"}
	assertSliceEqual(t, got, want)
}

func TestSplitHeadingsAndKeys(t *testing.T) {

	m := &model{
		categories: testData,
		headings:   []string{"firefox", "tmux", "vim"},
		maxWidth:   88,
	}
	gotHeadings, gotLines, gotLineCount := m.splitHeadingsAndKeys()

	wantHeadings := map[int]string{
		0: "\x1b[1mfirefox\x1b[0m",
		4: "\x1b[1mtmux\x1b[0m",
		6: "\x1b[1mvim\x1b[0m",
	}
	wantLines := map[int]string{
		1: " incognito\tctrl + shift + p ",
		2: " new tab\tctrl + shift + t ",
		3: " bookmarks bar\tctrl + b ",
		5: " close window\tshift + x ",
		7: " focus left\tctrl + h ",
		8: " swap left\tctrl + shift + h ",
	}
	wantLineCount := 9

	assertMapEqual(t, gotHeadings, wantHeadings)
	assertMapEqual(t, gotLines, wantLines)
	if gotLineCount != wantLineCount {
		t.Errorf("got %d, want %d", gotLineCount, wantLineCount)
	}
}

func assertSliceEqual(t testing.TB, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("length not equal: got (%d), want (%d)", len(got), len(want))
	}
	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}

func assertMapEqual(t testing.TB, got, want map[int]string) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
