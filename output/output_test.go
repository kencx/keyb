package output

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/kencx/keyb/ui"
	"github.com/kencx/keyb/ui/list"
	"github.com/kencx/keyb/ui/table"
)

var (
	testTable = table.New([]*table.Row{table.NewHeading("foo"), {Text: "bar"}})
	m         = &ui.Model{List: list.New("", testTable)}
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

	got, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	want := "foo      \nbar     "
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
