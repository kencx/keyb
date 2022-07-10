package table

import (
	"testing"
)

var (
	testRows = []*Row{
		NewHeading("heading"),
		{Text: "foo", Key: "bar"},
		{Text: "baz", Key: "foo"},
	}
)

func TestNew(t *testing.T) {
	t.Run("populated", func(t *testing.T) {
		tt := New(testRows)

		assertEqual(t, tt.LineCount, 3)
	})

	t.Run("empty", func(t *testing.T) {
		tt := NewEmpty(5)

		assertEqual(t, len(tt.Rows), 1)
		assertEqual(t, cap(tt.Rows), 5)
		assertEqual(t, tt.LineCount, 0)
		assertEqual(t, tt.Render(), "")
	})
}

func TestAppend(t *testing.T) {

	row1 := NewRow("baz", "", "", "")
	row2 := NewRow("foobar", "", "", "")

	t.Run("appendRow", func(t *testing.T) {
		tt := New(testRows)
		originalCount := tt.LineCount

		tt.AppendRow(row1)
		assertEqual(t, tt.LineCount, originalCount+1)
	})

	t.Run("appendRows", func(t *testing.T) {
		tt := New(testRows)
		originalCount := tt.LineCount

		tt.AppendRows(row1, row2)
		assertEqual(t, tt.LineCount, originalCount+2)
	})
}

func TestJoin(t *testing.T) {
	t1 := New(testRows)
	t1Count := t1.LineCount
	t2 := New([]*Row{{Text: "baz"}})
	t2Count := t2.LineCount
	t1.Join(t2)

	assertEqual(t, t1.LineCount, t1Count+t2Count)
}

func TestReset(t *testing.T) {
	tt := New(testRows)
	tt.Reset()

	assertEqual(t, tt.LineCount, 0)
	assertEqual(t, tt.Render(), "")
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}
