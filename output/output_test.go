package output

import (
	"io"
	"os"
	"testing"

	"github.com/kencx/keyb/config"
	"github.com/kencx/keyb/ui"
	"github.com/kencx/keyb/ui/list"
	"github.com/kencx/keyb/ui/table"
)

var (
	testTable  = table.New([]*table.Row{table.NewHeading("foo"), {Text: "bar"}})
	testConfig = &config.Config{}
	m          = &ui.Model{List: list.New(testTable, testConfig)}
)

func TestToStdout(t *testing.T) {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := ToStdout(m)
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	got, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	want := "foo      \nbar     "
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
