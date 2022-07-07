package table

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Row struct {
	Text       string
	Key        string
	Prefix     string
	ShowPrefix bool
	Styles     RowStyles

	IsHeading  bool
	IsSelected bool
	IsFiltered bool
}

type RowStyles struct {
	Normal   lipgloss.Style
	Selected lipgloss.Style
	Filtered lipgloss.Style
}

func NewHeading(text string) Row {
	return Row{
		Text:      text,
		IsHeading: true,
	}
}

// non-heading row
func NewRow(text, key, prefix string) Row {
	r := Row{
		Text:   text,
		Key:    key,
		Prefix: prefix,
	}
	r.ShowPrefix = r.Prefix != ""
	return r
}

func EmptyRow() Row {
	return Row{}
}

func (r *Row) String() string {
	if r.Text == "" && r.Key == "" {
		return ""
	}

	if r.IsHeading {
		return fmt.Sprintf("%s\t ", r.Text)
	}

	if !r.ShowPrefix {
		return fmt.Sprintf("%s\t%s", r.Text, r.Key)
	}
	return fmt.Sprintf("%s\t%s ; %s", r.Text, r.Prefix, r.Key)
}

func (r *Row) Render(m *Model) string {
	return r.Styles.Normal.Render(r.String())
}
