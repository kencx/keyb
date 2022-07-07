package list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {

	status := fmt.Sprintf(" keys: %d", m.table.LineCount)
	if m.debug {
		status = fmt.Sprintf("%s\tLine: %d YOffset: %d Height: %d",
			status, m.cursor, m.viewport.YOffset, m.viewport.Height)
	}

	topbar := m.title
	if m.filterState == filtering {
		separator := strings.Repeat(" ", m.viewport.Width-m.padding-lipgloss.Width(m.title)-lipgloss.Width(m.currentHeading))
		topbar = lipgloss.JoinHorizontal(lipgloss.Center, m.title, separator, m.currentHeading)
	}

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		topbar,
		m.searchBar.View(),
		m.viewport.View(),
		status,
	)

	style := lipgloss.NewStyle().
		Margin(0, 1)
		// Width(m.viewport.Width)
	return style.Render(view)
}
