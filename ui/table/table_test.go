package table

import (
	"reflect"
	"testing"
)

var testRows = []string{"foo", "bar"}

func TestNew(t *testing.T) {
	t.Run("populated", func(t *testing.T) {
		tt := New("foo", testRows)

		assertEqual(t, tt.heading, "foo")
		assertEqual(t, tt.LineCount, 3)
		assertSliceEqual(t, tt.Output, []string{"foo", "foo", "bar"})
		assertSliceEqual(t, tt.StyledOutput, []string{"foo", "foo", "bar"})
	})

	t.Run("empty", func(t *testing.T) {
		tt := NewEmpty(5)

		assertEqual(t, tt.heading, "")
		assertEqual(t, len(tt.rows), 1)
		assertEqual(t, cap(tt.rows), 5)
		assertEqual(t, tt.LineCount, 0)
		assertSliceEqual(t, tt.Output, []string(nil))
		assertSliceEqual(t, tt.StyledOutput, []string(nil))
	})
}

func TestAppend(t *testing.T) {

	t.Run("appendRow", func(t *testing.T) {
		tt := New("foo", testRows)
		originalCount := tt.LineCount

		tt.AppendRow("baz")
		assertEqual(t, tt.LineCount, originalCount+1)
		assertSliceEqual(t, tt.Output, []string{"foo", "foo", "bar", "baz"})
	})

	t.Run("appendRows", func(t *testing.T) {
		tt := New("foo", testRows)
		originalCount := tt.LineCount

		tt.AppendRows("baz", "foobar")
		assertEqual(t, tt.LineCount, originalCount+2)
		assertSliceEqual(t, tt.Output, []string{"foo", "foo", "bar", "baz", "foobar"})
	})
}

func TestAssemble(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		tt := New("foo", testRows)
		assertSliceEqual(t, tt.Output, []string{"foo", "foo", "bar"})
	})

	t.Run("missing heading", func(t *testing.T) {
		tt := New("", testRows)
		assertSliceEqual(t, tt.Output, []string{"foo", "bar"})
	})

	t.Run("missing row", func(t *testing.T) {
		tt := New("foo", []string{"", "bar"})
		assertSliceEqual(t, tt.Output, []string{"foo", "bar"})
	})

	t.Run("missing rows", func(t *testing.T) {
		tt := New("foo", []string{"", "bar", "\n"})
		assertSliceEqual(t, tt.Output, []string{"foo", "bar"})
	})
}

func TestJoin(t *testing.T) {
	t1 := New("table 1", testRows)
	t1Count := t1.LineCount
	t2 := New("table 2", []string{"baz"})
	t2Count := t2.LineCount
	t1.Join(t2)

	assertSliceEqual(t, t1.rows, []string{"foo", "bar", "baz"})
	assertSliceEqual(t, t1.Output, []string{"table 1", "foo", "bar", "table 2", "baz"})
	assertSliceEqual(t, t1.StyledOutput, []string{"table 1", "foo", "bar", "table 2", "baz"})
	assertEqual(t, t1.LineCount, t1Count+t2Count)
}

func TestReset(t *testing.T) {
	tt := New("foo", testRows)
	tt.Reset()

	assertEqual(t, tt.heading, "")
	assertEqual(t, tt.LineCount, 0)
	assertSliceEqual(t, tt.rows, nil)
	assertSliceEqual(t, tt.Output, nil)
	assertSliceEqual(t, tt.StyledOutput, nil)
}

func assertEqual[T comparable](t *testing.T, got, want T) {
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertSliceEqual[T comparable](t *testing.T, got, want []T) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
