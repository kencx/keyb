package list

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {

	if m.table.Empty() {
		m.viewport.SetContent("\nNo key bindings found")
	}

	header := m.title
	counter := formCounter(m)

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		m.searchBar.View(),
		counter,
		m.viewport.View(),
	)

	style := m.border.
		Margin(m.margin).
		Padding(m.padding).
		// border should occupy entire width
		Width(m.viewport.Width - m.margin)
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
	return lipgloss.NewStyle().Faint(true).Margin(0, 1).Render(counter)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
