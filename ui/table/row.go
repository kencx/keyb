package table

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Row struct {
	Text      string
	Key       string
	Prefix    string
	PrefixSep string

	// default false unless prefix defined
	ShowPrefix bool
	// only used to show row's corresponding heading during filtering
	Heading string

	MatchedIndex []int
	Styles       RowStyles

	IsHeading  bool
	IsSelected bool
	IsFiltered bool
	Reversed   bool
}

type RowStyles struct {
	Normal          lipgloss.Style
	Heading         lipgloss.Style
	Selected        lipgloss.Style
	SelectedHeading lipgloss.Style
	Filtered        lipgloss.Style
}

func NewHeading(text string) *Row {
	return &Row{
		Text:      text,
		IsHeading: true,
	}
}

func NewRow(text, key, prefix, heading string) *Row {
	r := &Row{
		Text:    text,
		Key:     key,
		Prefix:  prefix,
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

	if r.Reversed {
		return r.ReverseString()
	}

	if !r.ShowPrefix {
		return fmt.Sprintf("%s\t%s", r.Text, r.Key)
	}

	key := fmt.Sprintf("%s %s %s", r.Prefix, r.PrefixSep, r.Key)
	return fmt.Sprintf("%s\t%s", r.Text, key)
}

func (r *Row) ReverseString() string {
	if !r.ShowPrefix {
		return fmt.Sprintf("%s\t%s", r.Key, r.Text)
	}

	key := fmt.Sprintf("%s %s %s", r.Prefix, r.PrefixSep, r.Key)
	return fmt.Sprintf("%s\t%s", key, r.Text)
}

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
