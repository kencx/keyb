package list

import (
	"keyb/table"

	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Title     lipgloss.Style
	Cursor    lipgloss.Style
	Highlight lipgloss.Style
	Border    lipgloss.Style
	Table     table.Styles
}

func DefaultStyle() Styles {
	table := table.Styles{
		BodyStyle:    lipgloss.NewStyle(),
		HeadingStyle: lipgloss.NewStyle().Bold(true).Margin(0, 1),
		RowStyle:     lipgloss.NewStyle().Margin(0, 2),
	}
	return Styles{
		Cursor:    lipgloss.NewStyle().Background(lipgloss.Color("#448448")).Foreground(lipgloss.Color("#edb")),
		Highlight: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA066")),
		Table:     table,
	}
}
