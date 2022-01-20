package main

import (
	"reflect"
	"testing"
)

var testData = map[string]Program{
	"tmux": {Prefix: "ctrl + b",
		KeyBinds: []KeyBind{
			{Desc: "close window", Key: "shift + x"},
		}},
	"vim": {KeyBinds: []KeyBind{
		{Desc: "focus left", Key: "ctrl + h"},
		{Desc: "swap left", Key: "ctrl + shift + h"},
	}},
	"firefox": {Prefix: "test",
		KeyBinds: []KeyBind{
			{Desc: "incognito", Key: "ctrl + shift + p", Ignore_Prefix: true},
			{Desc: "new tab", Key: "ctrl + shift + t", Ignore_Prefix: true},
			{Desc: "bookmarks bar", Key: "ctrl + b", Ignore_Prefix: true},
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
		t.Errorf("\ngot %q\nwant %q", got, want)
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
