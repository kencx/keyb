package table

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Row struct {
	Text   string
	Key    string
	Prefix string
	Style  lipgloss.Style

	IsHeading    bool
	IgnorePrefix bool
}

func (r Row) String() string {
	if r.IsHeading {
		return fmt.Sprintf("%s\t ", r.Text)
	}

	if r.IgnorePrefix {
		return fmt.Sprintf("%s\t%s", r.Text, r.Key)
	}
	return fmt.Sprintf("%s\t%s ; %s", r.Text, r.Prefix, r.Key)
}

func (r *Row) Render() string {
	return r.Style.Render(r.String())
}

func NewHeading(text, key string) Row {
	return Row{
		Text:      text,
		Key:       key,
		IsHeading: true,
	}
}

// New non-heading row
func NewRow(text, key, prefix string) Row {
	r := Row{
		Text:      text,
		Key:       key,
		Prefix:    prefix,
		IsHeading: false,
	}
	r.IgnorePrefix = r.Prefix == ""
	return r
}

func EmptyRow() Row {
	return Row{}
}
