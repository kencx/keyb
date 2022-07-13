package table

import (
	"fmt"
	"strings"

	"github.com/juju/ansiterm/tabwriter"
	"github.com/muesli/reflow/truncate"
)

type Model struct {
	Rows      []*Row
	LineCount int
	SepWidth  int
	MaxWidth  int // prevents line wrapping
}

func New(rows []*Row) *Model {
	t := &Model{
		Rows: rows,
	}

	for _, row := range rows {
		if row.String() != "" {
			t.LineCount += 1
		}
	}
	return t
}

func NewEmpty(n int) *Model {
	return &Model{
		Rows:      make([]*Row, 1, max(2, n)),
		LineCount: 0,
	}
}

func (t *Model) Empty() bool {
	return t.LineCount <= 0
}

func (t *Model) AppendRow(row *Row) {
	t.Rows = append(t.Rows, row)
	t.LineCount += 1
}

func (t *Model) AppendRows(rows ...*Row) {
	t.Rows = append(t.Rows, rows...)
	t.LineCount += len(rows)
}

func (t *Model) Join(table *Model) {
	t.Rows = append(t.Rows, table.Rows...)
	t.LineCount += table.LineCount
}

func (t *Model) Reset() {
	t.Rows = nil
	t.LineCount = 0
}

// Align and style rows
func (t *Model) Render() string {
	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 8, 4, t.SepWidth, ' ', 0)

	for _, row := range t.Rows {
		if row != nil && row.String() != "" {
			fmt.Fprintln(tw, row.Render())
		}
	}
	tw.Flush()

	if sb.String() != "" {
		sl := strings.Split(strings.TrimSuffix(sb.String(), "\n"), "\n")
		sb.Reset()

		// Unable to truncate while aligning due to nature of tabwriter
		for _, row := range sl {
			if t.MaxWidth > 0 {
				fmt.Fprintln(&sb, truncate.StringWithTail(row, uint(t.MaxWidth), "..."))
			} else {
				fmt.Fprintln(&sb, row)
			}
		}
	}
	return sb.String()
}

// Helper functions for retrieving specific rows in a table
// Aligned but unstyled rows
func (t *Model) GetAlignedRows() string {
	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 8, 4, t.SepWidth, ' ', 0)

	for _, row := range t.Rows {
		if row != nil && row.String() != "" {
			fmt.Fprintln(tw, row.String())
		}
	}

	tw.Flush()
	return strings.TrimSuffix(sb.String(), "\n")
}

func (t *Model) GetPlainHeadings() []string {
	var res []string
	for _, r := range t.Rows {
		if r.IsHeading {
			res = append(res, r.String())
		}
	}
	return res
}

func (t *Model) GetPlainRowsWithoutHeadings() []string {
	var res []string
	for _, r := range t.Rows {
		if !r.IsHeading {
			res = append(res, r.String())
		}
	}
	return res
}

func (t *Model) GetCopyOfHeadings() []Row {
	var res []Row
	for _, r := range t.Rows {
		if r.IsHeading {
			res = append(res, *r)
		}
	}
	return res
}

func (t *Model) GetCopyOfRowsWithoutHeadings() []Row {
	var res []Row
	for _, r := range t.Rows {
		if !r.IsHeading {
			res = append(res, *r)
		}
	}
	return res
}

func (t *Model) GetAllRowsofHeading(heading string) []*Row {
	var res []*Row
	for _, r := range t.Rows {
		if !r.IsHeading && r.Heading == heading {
			res = append(res, r)
		}
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
