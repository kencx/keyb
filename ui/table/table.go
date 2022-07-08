package table

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/juju/ansiterm/tabwriter"
)

type Model struct {
	Rows      []*Row
	Output    []string
	LineCount int
	SepWidth  int
	Styles
}

type Styles struct {
	HeadingStyle lipgloss.Style
	RowStyle     lipgloss.Style
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

	t.Render()
	return t
}

func NewWithStyle(rows []*Row, style Styles) *Model {
	table := New(rows)
	table.Styles = style
	return table
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
	t.Render()
}

func (t *Model) AppendRows(rows ...*Row) {
	t.Rows = append(t.Rows, rows...)
	t.LineCount += len(rows)
	t.Render()
}

func (t *Model) Join(table *Model) {
	t.Rows = append(t.Rows, table.Rows...)
	t.Output = append(t.Output, table.Output...)
	t.LineCount += table.LineCount
}

func (t *Model) Reset() {
	t.Rows = nil
	t.Output = nil
	t.LineCount = 0
}

// Align and style rows
func (t *Model) Render() {
	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 8, 4, t.SepWidth, ' ', 0)

	// don't use Align here to not have 2 for loops
	for _, row := range t.Rows {
		if row != nil && row.String() != "" {

			if row.IsHeading {
				row.Styles.Heading = t.HeadingStyle
				fmt.Fprintln(tw, row.Render())
			} else {
				row.Styles.Normal = t.RowStyle
				fmt.Fprintln(tw, row.Render())
			}
		}
	}

	tw.Flush()
	s := strings.TrimSuffix(sb.String(), "\n")
	t.Output = strings.Split(s, "\n")
}

func (t *Model) String() string {
	return strings.Join(t.Output, "\n")
}

// The next 3 functions need better names
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

func (t *Model) GetPlainRowsWithoutHeadings() []string {
	var res []string
	for _, r := range t.Rows {
		if !r.IsHeading {
			res = append(res, r.String())
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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
