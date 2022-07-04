package table

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
)

type Table struct {
	heading   string
	rows      []string
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
		heading: heading,
		rows:    rows,
		Styles:  style,
	}

	if heading != "" {
		t.LineCount += 1
	}

	if len(rows) > 0 && rows[0] != "" {
		t.LineCount += len(rows)
	}

	t.Render()
	return t
}

func (t *Table) AppendRow(row string) {
	t.rows = append(t.rows, row)
	t.LineCount += 1
	t.Render()
}

func (t *Table) AppendRows(rows ...string) {
	t.rows = append(t.rows, rows...)
	t.LineCount += len(rows)
	t.Render()
}

// func (t *Table) PrependRow(row string) {
// 	t.rows = append([]string{row}, t.rows...)
// }

func (t *Table) Render() {
	if t.heading != "" {
		heading := t.HeadingStyle.Render(t.heading)
		t.Output = append(t.Output, heading)
	}

	for _, row := range t.rows {
		if row != "" {
			line := t.LineStyle.Render(row)
			t.Output = append(t.Output, line)
		}
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

func (t *Table) Empty() bool {
	return t.LineCount <= 0
}

func (t *Table) ParsedOutput() string {
	return strings.Join(t.Output, "\n")
}

func (t *Table) Reset() {
	t.heading = ""
	t.rows = nil
	t.Output = nil
	t.LineCount = 0
}
