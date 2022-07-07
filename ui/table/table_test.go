package table

import (
	"reflect"
	"testing"
)

var (
	testHeading = NewHeading("heading")
	testRows    = []Row{
		{Text: "foo", Key: "bar"},
		{Text: "baz", Key: "foo"},
	}
	want = []string{testHeading.String(), testRows[0].String(), testRows[1].String()}
)

func TestNew(t *testing.T) {
	t.Run("populated", func(t *testing.T) {
		tt := New(testHeading, testRows)

		assertEqual(t, tt.heading.String(), "heading\t ")
		assertEqual(t, tt.LineCount, 3)
		assertSliceEqual(t, tt.Output, want)
		assertSliceEqual(t, tt.StyledOutput, want)
	})

	t.Run("empty", func(t *testing.T) {
		tt := NewEmpty(5)

		assertEqual(t, tt.heading.String(), "")
		assertEqual(t, len(tt.rows), 1)
		assertEqual(t, cap(tt.rows), 5)
		assertEqual(t, tt.LineCount, 0)
		assertSliceEqual(t, tt.Output, []string(nil))
		assertSliceEqual(t, tt.StyledOutput, []string(nil))
	})
}

func TestAppend(t *testing.T) {

	row1 := NewRow("baz", "", "")
	row2 := NewRow("foobar", "", "")

	t.Run("appendRow", func(t *testing.T) {
		tt := New(testHeading, testRows)
		originalCount := tt.LineCount
		want1 := append(want, row1.String())

		tt.AppendRow(row1)
		assertEqual(t, tt.LineCount, originalCount+1)
		assertSliceEqual(t, tt.Output, want1)
	})

	t.Run("appendRows", func(t *testing.T) {
		tt := New(testHeading, testRows)
		originalCount := tt.LineCount
		want2 := append(want, row1.String(), row2.String())

		tt.AppendRows(row1, row2)
		assertEqual(t, tt.LineCount, originalCount+2)
		assertSliceEqual(t, tt.Output, want2)
	})
}

func TestAssemble(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		tt := New(testHeading, testRows)
		assertSliceEqual(t, tt.Output, want)
	})

	t.Run("missing heading", func(t *testing.T) {
		tt := New(EmptyRow(), testRows)
		assertSliceEqual(t, tt.Output, []string{testRows[0].String(), testRows[1].String()})
	})

	t.Run("missing row", func(t *testing.T) {
		tt := New(testHeading, []Row{EmptyRow(), {Text: "foo"}})
		assertSliceEqual(t, tt.Output, []string{testHeading.String(), "foo\t"})
	})

	t.Run("missing rows", func(t *testing.T) {
		tt := New(testHeading, []Row{EmptyRow(), {Text: "bar"}, EmptyRow()})
		assertSliceEqual(t, tt.Output, []string{testHeading.String(), "bar\t"})
	})
}

func TestJoin(t *testing.T) {
	t1 := New(NewHeading("table 1"), testRows)
	t1Count := t1.LineCount
	t2 := New(NewHeading("table 2"), []Row{{Text: "baz"}})
	t2Count := t2.LineCount
	t1.Join(t2)

	assertSliceEqual(t, t1.Output, []string{"table 1\t ", "foo\tbar", "baz\tfoo", "table 2\t ", "baz\t"})
	assertSliceEqual(t, t1.StyledOutput, []string{"table 1\t ", "foo\tbar", "baz\tfoo", "table 2\t ", "baz\t"})
	assertEqual(t, t1.LineCount, t1Count+t2Count)
}

func TestReset(t *testing.T) {
	tt := New(testHeading, testRows)
	tt.Reset()

	assertEqual(t, tt.heading.String(), "")
	assertEqual(t, tt.LineCount, 0)
	assertSliceEqual(t, tt.Output, nil)
	assertSliceEqual(t, tt.StyledOutput, nil)
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func assertSliceEqual[T comparable](t *testing.T, got, want []T) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}
