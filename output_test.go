package main

import (
	"io/ioutil"
	"keyb/table"
	"os"
	"testing"
)

var m = &model{
	Table: &table.Table{
		Output: []string{"\x1b[1mthis is a test string\x1b[0m", "followed by the second line"},
	},
}

func TestOutputBodyToStdout(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := m.OutputBodyToStdout(true)
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	got, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	want := "this is a test string\nfollowed by the second line"
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
