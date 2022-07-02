package table

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
)

type Table struct {
	Heading   string
	Rows      []string
	Output    []string
	LineCount int
	Styles
}

type Styles struct {
	BodyStyle    lipgloss.Style
	HeadingStyle lipgloss.Style
	LineStyle    lipgloss.Style
}

func New(heading string, rows []string, style Styles) *Table {
	t := &Table{
		Heading: heading,
		Rows:    rows,
		Styles:  style,
	}
	t.Render()
	return t
}

// func (t *Table) AppendRow(row string) {
// 	t.Rows = append(t.Rows, row)
// }
//
// func (t *Table) AppendRows(rows ...string) {
// 	t.Rows = append(t.Rows, rows...)
// }
//
// func (t *Table) PrependRow(row string) {
// 	t.Rows = append([]string{row}, t.Rows...)
// }

func (t *Table) Render() {
	heading := t.HeadingStyle.Render(t.Heading)
	t.Output = append(t.Output, heading)

	for _, row := range t.Rows {
		line := t.LineStyle.Render(row)
		t.Output = append(t.Output, line)
	}
}

func (t *Table) Align() {
	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 20, 8, 10, ' ', 0)

	for _, row := range t.Output {
		fmt.Fprintln(tw, row)
	}
	tw.Flush()
	t.Output = strings.Split(sb.String(), "\n")
}

func (t *Table) Join(table *Table) {
	t.Output = append(t.Output, table.Output...)
	t.LineCount = len(t.Output)
}
