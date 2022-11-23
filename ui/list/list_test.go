package list

import (
	"testing"

	"github.com/kencx/keyb/config"
	"github.com/kencx/keyb/ui/table"
)

var (
	testRows = []*table.Row{
		table.NewHeading("fooTable"),
		{Text: "foo"},
		{Text: "bar"},
		{Text: "baz"},
	}
	testTable  = table.New(testRows)
	testConfig = &config.Config{
		Settings: config.Settings{
			Title:       "foo",
			Debug:       true,
			Reverse:     true,
			Mouse:       true,
			Prompt:      "prompt",
			Placeholder: "placeholder",
			PrefixSep:   "$",
			SepWidth:    4,
		}}
)

func TestNew(t *testing.T) {
	t.Run("populated", func(t *testing.T) {
		tm := New(testTable, testConfig)

		assertEqual(t, tm.table.LineCount, 4)
		assertEqual(t, tm.filterState, unfiltered)
		assertEqual(t, tm.filteredTable.LineCount, 0)

		assertEqual(t, tm.title, "foo")
		assertEqual(t, tm.debug, true)
		assertEqual(t, tm.table.Rows[0].Reversed, true)
		assertEqual(t, tm.viewport.MouseWheelEnabled, true)
		assertEqual(t, tm.searchBar.Prompt, "prompt")
		assertEqual(t, tm.searchBar.Placeholder, "placeholder")
		assertEqual(t, tm.table.Rows[0].PrefixSep, "$")
		assertEqual(t, tm.table.SepWidth, 4)
	})

	t.Run("empty", func(t *testing.T) {
		tm := New(table.New([]*table.Row{table.EmptyRow(), table.EmptyRow()}), testConfig)

		assertEqual(t, tm.title, "foo")
		assertEqual(t, tm.table.LineCount, 0)
		assertEqual(t, tm.filterState, unfiltered)
		assertEqual(t, tm.filteredTable.LineCount, 0)
	})

}

func TestReset(t *testing.T) {

	tm := New(testTable, testConfig)
	tm.filterState = filtering
	tm.searchBar.SetValue("searching...")
	tm.cursor = 10

	tm.Reset()

	assertEqual(t, tm.filteredTable.Render(), "")
	assertEqual(t, tm.searchBar.Value(), "searching...")
	assertEqual(t, tm.filterState, unfiltered)
	assertEqual(t, tm.cursor, 0)
	assertEqual(t, tm.maxRows, tm.table.LineCount)
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}
