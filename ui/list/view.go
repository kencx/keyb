package list

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {

	status := fmt.Sprintf(" keys: %d", m.table.LineCount)
	if m.debug {
		status = fmt.Sprintf("%s\tLine: %d YOffset: %d Height: %d",
			status, m.cursor, m.viewport.YOffset, m.viewport.Height)
	}

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		m.Title,
		m.searchBar.View(),
		m.viewport.View(),
		status,
	)

	style := lipgloss.NewStyle().
		Margin(0, 1)
		// Width(m.viewport.Width)
		// MaxWidth(m.viewport.Width)
	return style.Render(view)
}
