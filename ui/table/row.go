package table

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Row struct {
	Text   string
	Key    string
	Prefix string
	// row prefix ignore defaults to true
	// which equates to prefix show defaulting to false
	// otherwise, all rows will show (empty) prefix regardless
	ShowPrefix bool

	MatchedIndex []int
	Styles       RowStyles

	// only used to show row's corresponding heading during filtering
	Heading    string
	IsHeading  bool
	IsSelected bool
	IsFiltered bool
}

type RowStyles struct {
	Normal          lipgloss.Style
	Heading         lipgloss.Style
	Selected        lipgloss.Style
	SelectedHeading lipgloss.Style
	Filtered        lipgloss.Style
}

func DefaultRowStyles() RowStyles {
	return RowStyles{
		Selected: lipgloss.NewStyle().
			Background(lipgloss.Color("#448448")).
			Foreground(lipgloss.Color("#edb")).
			Margin(0, 2),
		SelectedHeading: lipgloss.NewStyle().
			Background(lipgloss.Color("#448448")).
			Foreground(lipgloss.Color("#edb")).
			Margin(0, 1).
			Bold(true),
		Filtered: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA066")),
	}
}

func NewHeading(text string) *Row {
	return &Row{
		Text:      text,
		IsHeading: true,
		Styles:    DefaultRowStyles(),
	}
}

// non-heading row
func NewRow(text, key, prefix, heading string) *Row {
	r := &Row{
		Text:    text,
		Key:     key,
		Prefix:  prefix,
		Styles:  DefaultRowStyles(),
		Heading: heading,
	}
	r.ShowPrefix = r.Prefix != ""
	return r
}

func EmptyRow() *Row {
	return &Row{}
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

// margins and paddings in stylerunes seem to mess it up
func (r *Row) Render() string {
	s := r.Styles

	if r.IsSelected {
		if r.IsFiltered {
			// Inline to remove margins, paddings and borders from styledrunes
			unmatched := s.Selected.Inline(true)
			matched := s.Filtered.Copy().Inherit(unmatched)
			str := lipgloss.StyleRunes(r.String(), r.MatchedIndex, matched, unmatched)

			if r.IsHeading {
				return s.SelectedHeading.Render(str)
			}
			return s.Selected.Render(str)
		}

		if r.IsHeading {
			return s.SelectedHeading.Render(r.String())
		}
		return s.Selected.Render(r.String())
	}

	if r.IsFiltered {
		unmatched := s.Normal.Inline(true)
		matched := s.Filtered.Copy().Inherit(unmatched)
		str := lipgloss.StyleRunes(r.String(), r.MatchedIndex, matched, unmatched)

		if r.IsHeading {
			return s.Heading.Render(str)
		}
		return s.Normal.Render(str)
	}

	if r.IsHeading {
		return s.Heading.Render(r.String())
	}
	return s.Normal.Render(r.String())
}
