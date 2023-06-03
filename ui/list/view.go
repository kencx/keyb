package list

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	if m.table.Empty() {
		m.viewport.SetContent("\nNo key bindings found")
	}

	counter := formCounter(m)

	var view string
	if m.promptLocation == "bottom" {
		view = lipgloss.JoinVertical(
			lipgloss.Left,
			m.viewport.View(),
			counter,
			m.searchBar.View(),
		)
	} else {
		view = lipgloss.JoinVertical(
			lipgloss.Left,
			m.searchBar.View(),
			counter,
			m.viewport.View(),
		)
	}

	style := m.border.
		Margin(m.margin).
		Padding(m.padding).
		Width(m.viewport.Width)
	return style.Render(view)
}

func formCounter(m *Model) string {
	var counter string

	if m.filterState == filtering && m.searchBar.Value() != "" {
		counter = fmt.Sprintf("%d/%d %s", m.filteredTable.LineCount, m.table.LineCount, m.currentHeading)
	} else {
		counter = fmt.Sprintf("%d/%d %s", m.table.LineCount, m.table.LineCount, m.currentHeading)
	}

	if m.debug {
		counter = fmt.Sprintf("%s\tLine: %d YOffset: %d Height: %d",
			counter, m.cursor, m.viewport.YOffset, m.viewport.Height)
	}
	return m.counterStyle.Render(counter)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
