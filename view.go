package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) View() string {
	if !m.ready {
		return "\n Initializing..."
	}

	view := fmt.Sprintf("%s\n"+
		"%s\n"+
		"%s\n"+
		"\n"+
		" keys: %d",
		m.Title,
		m.searchBar.View(),
		m.viewport.View(),
		m.table.LineCount)

	if m.debug {
		view = fmt.Sprintf("%s\tLine: %d YOffset: %d Height: %d",
			view, m.cursor, m.viewport.YOffset, m.viewport.Height)
	}
	// m.setStyle()
	return m.table.BodyStyle.Render(view)
}

// set global styles
func (m *model) setStyle() {
	m.table.BodyStyle = m.table.BodyStyle.Margin(1, 2)
	m.handleBorder()
}

// TODO issues with border:
// does not resize with window size
// width changes when lines shrink and grow
// set border fixed to window size like fzf
func (m *model) handleBorder() {
	var borderStyle lipgloss.Border

	switch m.Border {
	case "normal":
		borderStyle = lipgloss.NormalBorder()
	case "rounded":
		borderStyle = lipgloss.RoundedBorder()
	case "double":
		borderStyle = lipgloss.DoubleBorder()
	case "thick":
		borderStyle = lipgloss.ThickBorder()
	default:
		borderStyle = lipgloss.HiddenBorder()
	}

	m.table.BodyStyle = m.table.BodyStyle.Border(borderStyle)
	m.padding += m.table.BodyStyle.GetBorderTopWidth() + m.table.BodyStyle.GetBorderBottomSize()
}
