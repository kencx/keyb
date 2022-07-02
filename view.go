package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.cursor--
			m.updateCursor()
		case "down", "j":
			m.cursor++
			m.updateCursor()
		case "G":
			m.cursor = m.Table.LineCount - 1
		case "ctrl+d":
			m.cursor += m.height / 2
			m.updateCursor()
		case "ctrl+u":
			m.cursor -= m.height / 2
			m.updateCursor()
		}
	case tea.MouseMsg:
		if !m.mouseEnabled {
			break
		}
		switch msg.Type {

		// TODO smoother scrolling
		case tea.MouseWheelUp:
			m.cursor -= m.mouseDelta
			m.updateCursor()
		case tea.MouseWheelDown:
			m.cursor += m.mouseDelta
			m.updateCursor()
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - m.padding
		m.offset = 0

		// refresh line wrapping (WIP)
		// if msg.Width/2 < m.maxWidth {
		// 	m.maxWidth = msg.Width / 2
		// 	m.initBody()
		// }
		// if msg.Width > m.maxWidth*2 {
		// 	m.maxWidth = msg.Width
		// 	m.initBody()
		// }
	}
	m.updateOffset()
	return m, nil
}

func (m *model) updateCursor() {
	if m.cursor < 0 {
		m.cursor = m.Table.LineCount - 1

	} else if m.cursor > m.Table.LineCount-1 {
		m.cursor = 0
	}
}

// scrolling
func (m *model) updateOffset() {
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+m.height {
		m.offset = m.cursor - m.height + 1
	}
}
func (m *model) View() string {
	body := m.updateBody()
	view := fmt.Sprintf("%s\n\n%s", m.Table.Title, body)

	m.setStyle()
	return m.Table.BodyStyle.Render(view)

}

func (m *model) updateBody() string {

	if len(m.Table.Output) == 0 {
		return "No key bindings found"
	}

	// deep copy slice
	cpy := make([]string, len(m.Table.Output))
	copy(cpy, m.Table.Output)

	body := m.renderCursor(cpy)

	if m.Table.LineCount >= m.offset+m.height {
		body = body[m.offset : m.offset+m.height]
	}
	result := strings.Join(body, "\n")
	return result
}

// render cursor style at position
func (m *model) renderCursor(lines []string) []string {
	cursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(m.Table.CursorBackground)).
		Foreground(lipgloss.Color(m.Table.CursorForeground))

	for i, line := range lines {
		if m.cursor == i {
			lines[i] = cursorStyle.Render(line)
		}
	}
	return lines
}

func (m *model) setStyle() {
	m.Table.BodyStyle = m.Table.BodyStyle.Margin(1, 2)
	m.handleBorder()
}

// TODO issues with border:
// does not resize with window size
// width changes when lines shrink and grow
// set border fixed to window size like fzf
func (m *model) handleBorder() {
	var borderStyle lipgloss.Border

	switch m.Table.Border {
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

	m.Table.BodyStyle = m.Table.BodyStyle.Border(borderStyle)
	m.padding += m.Table.BodyStyle.GetBorderTopWidth() + m.Table.BodyStyle.GetBorderBottomSize()
}
