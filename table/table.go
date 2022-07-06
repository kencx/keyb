package table

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	heading      string
	rows         []string
	Output       []string
	StyledOutput []string
	LineCount    int
	Styles
}

type Styles struct {
	BodyStyle    lipgloss.Style
	HeadingStyle lipgloss.Style
	RowStyle     lipgloss.Style
}

func New(heading string, rows []string) *Model {
	t := &Model{
		heading: heading,
		rows:    rows,
	}

	if heading != "" {
		t.LineCount += 1
	}

	if len(rows) > 0 && rows[0] != "" {
		t.LineCount += len(rows)
	}

	// generate outputs
	t.Update()
	return t
}

func NewWithStyle(heading string, rows []string, style Styles) *Model {
	table := New(heading, rows)
	table.Styles = style
	return table
}

// New default empty table with empty heading & maximum of n rows.
func NewEmpty(n int) *Model {
	return &Model{
		rows:      make([]string, 1, max(2, n)),
		Styles:    DefaultStyles(),
		LineCount: 0,
	}
}

func DefaultStyles() Styles {
	return Styles{
		BodyStyle:    lipgloss.NewStyle(),
		HeadingStyle: lipgloss.NewStyle(),
		RowStyle:     lipgloss.NewStyle(),
	}
}

func (t *Model) String() string {
	return strings.Join(t.StyledOutput, "\n")
}

func (t *Model) Empty() bool {
	return t.LineCount <= 0
}

func (t *Model) AppendRow(row string) {
	t.rows = append(t.rows, row)
	t.LineCount += 1
	t.Update()
}

func (t *Model) AppendRows(rows ...string) {
	t.rows = append(t.rows, rows...)
	t.LineCount += len(rows)
	t.Update()
}

func (t *Model) PrependRow(row string) {
	t.rows = append([]string{row}, t.rows...)
}

func (t *Model) Update() {
	t.Assemble()
	t.Render()
}

// Generates unstyled output
func (t *Model) Assemble() {
	var rows []string
	if t.heading != "" && t.heading != "\n" {
		rows = append(rows, t.heading)
	}

	for _, row := range t.rows {
		if row != "" && row != "\n" {
			rows = append(rows, row)
		}
	}
	t.Output = rows
}

// Generates styled output
func (t *Model) Render() {
	var rows []string
	if t.heading != "" {
		heading := t.HeadingStyle.Render(t.heading)
		rows = append(rows, heading)
	}

	for _, row := range t.rows {
		if row != "" {
			line := t.RowStyle.Render(row)
			rows = append(rows, line)
		}
	}
	t.StyledOutput = rows
}

func (t *Model) Join(table *Model) {
	t.rows = append(t.rows, table.rows...)
	t.Output = append(t.Output, table.Output...)
	t.StyledOutput = append(t.StyledOutput, table.StyledOutput...)
	t.LineCount += table.LineCount
}

func (t *Model) Reset() {
	t.heading = ""
	t.rows = nil
	t.Output = nil
	t.LineCount = 0
}

// Only works properly after Render
// Calling this before Join will not align different sub-tables to each other
func (t *Model) Align() {
	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 20, 8, 10, ' ', 0)

	for _, row := range t.StyledOutput {
		fmt.Fprintln(tw, row)
	}
	tw.Flush()
	t.StyledOutput = strings.Split(sb.String(), "\n")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
