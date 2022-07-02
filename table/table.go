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

	Title            string
	CursorForeground string
	CursorBackground string
	Border           string
	BorderColor      string
}

func New(heading string, rows []string, style Styles) *Table {
	return &Table{
		heading:   heading,
		rows:      rows,
		Styles:    style,
		LineCount: len(rows) + 1,
	}
}

func (t *Table) AddRow(row string) {
	t.rows = append(t.rows, row)
	t.LineCount += 1
}

func (t *Table) AddRows(rows ...string) {
	t.rows = append(t.rows, rows...)
	t.LineCount += len(rows)
}

func (t *Table) align() {

	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 20, 8, 10, ' ', 0)

	for _, row := range t.rows {
		fmt.Fprintln(tw, row)
	}

	tw.Flush()
	t.rows = strings.Split(sb.String(), "\n")
}

func (t *Table) Render() {
	heading := t.HeadingStyle.Render(t.heading)
	t.Output = append(t.Output, heading)

	for _, row := range t.rows {
		line := t.LineStyle.Render(row)
		t.Output = append(t.Output, line)
	}
}
