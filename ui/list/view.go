package list

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kencx/keyb/util"
)

func (m *Model) View() string {

	header := m.title
	if m.filterState == filtering {
		separator := strings.Repeat(" ", util.Max(0,
			m.viewport.Width-m.padding-lipgloss.Width(m.title)-lipgloss.Width(m.currentHeading)),
		)
		header = lipgloss.JoinHorizontal(lipgloss.Center, m.title, separator, m.currentHeading)
	}

	footer := fmt.Sprintf(" keys: %d", m.table.LineCount)
	if m.debug {
		footer = fmt.Sprintf("%s\tLine: %d YOffset: %d Height: %d",
			footer, m.cursor, m.viewport.YOffset, m.viewport.Height)
	}

	view := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		m.searchBar.View(),
		m.viewport.View(),
		footer,
	)

	style := lipgloss.NewStyle().Margin(0, 1)
	return style.Render(view)
}
