package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/kencx/keyb/ui"
	"github.com/kencx/keyb/ui/list"
	"github.com/kencx/keyb/ui/table"
)

var (
	testTable = table.New(table.NewHeading("\x1b[1mfoo\x1b[0m"), []table.Row{{Text: "bar"}})
	m         = &ui.Model{List: list.New("", testTable)}
)

func TestOutputBodyToStdout(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := OutputBodyToStdout(m, true)
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	got, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	want := "foo\t \nbar\t"
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
