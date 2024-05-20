package output

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/kencx/keyb/config"
	"github.com/kencx/keyb/ui"
	"github.com/kencx/keyb/ui/list"
	"github.com/kencx/keyb/ui/table"
)

var (
	testTable  = table.New([]*table.Row{table.NewHeading("foo"), {Text: "bar"}})
	testConfig = &config.Config{}
	testApps   = &config.Apps{
		&config.App{
			Name:   "foo",
			Prefix: "bar",
			Keybinds: []config.KeyBind{
				{
					Name: "key foo",
					Key:  "key bar",
				},
			},
		},
	}
	m = &ui.Model{List: list.New(testTable, testConfig), Apps: testApps}
)

func TestToJson(t *testing.T) {
	tempDir := t.TempDir()
	path := filepath.Join(tempDir, "test.json")

	err := ToFile(m, path)
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	want := []byte(`[{"prefix":"bar","name":"foo","keybinds":[{"name":"key foo","key":"key bar"}]}]`)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %s, want %s", got, want)
	}
}

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
