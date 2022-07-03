package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	lastLine := m.Table.LineCount - 1

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			m.cursor--
			if m.cursorAtViewTop() {
				m.Viewport.LineUp(1)
			}
		case "down", "j":
			m.cursor++
			if m.cursorAtViewBottom() {
				m.Viewport.LineDown(1)
			}

		case "G":
			m.cursor = lastLine
			m.Viewport.GotoBottom()

		case "ctrl+u":
			m.cursor -= m.Viewport.Height / 2
			m.Viewport.HalfViewUp()

			// don't loop around
			if m.cursor < 0 {
				m.cursor = 0
				m.Viewport.GotoTop()
			}
		case "ctrl+d":
			m.cursor += m.Viewport.Height / 2
			m.Viewport.HalfViewDown()

			// don't loop around
			if m.cursor > lastLine {
				m.cursor = lastLine
				m.Viewport.GotoBottom()
			}
		}

	case tea.MouseMsg:
		if !m.mouseEnabled {
			break
		}
		switch msg.Type {
		case tea.MouseWheelUp:
			m.cursor -= m.Viewport.MouseWheelDelta
			if m.cursorAtViewTop() {
				m.Viewport.LineUp(m.Viewport.MouseWheelDelta)
			}
		case tea.MouseWheelDown:
			m.cursor += m.Viewport.MouseWheelDelta
			if m.cursorAtViewBottom() {
				m.Viewport.LineDown(m.Viewport.MouseWheelDelta)
			}
		}

	case tea.WindowSizeMsg:
		if !m.ready {
			m.Viewport = viewport.New(msg.Width, msg.Height-m.padding)
			m.Viewport.SetContent(strings.Join(m.Table.Output, "\n"))
			m.Viewport.MouseWheelEnabled = m.mouseEnabled
			m.ready = true
			m.Viewport.SetYOffset(0)

		} else {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height - m.padding
		}

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

	// cursor loop around
	if m.cursorAtTop() {
		m.cursor = lastLine
		m.Viewport.GotoBottom()
	} else if m.cursorAtBottom() {
		m.cursor = 0
		m.Viewport.GotoTop()
	}

	m.updateYOffset()
	m.renderCursor()
	return m, nil
}

func (m *model) updateYOffset() {
	if m.cursor < m.Viewport.YOffset {
		m.Viewport.YOffset = m.cursor
	}
	if m.cursor >= m.Viewport.YOffset+m.Viewport.Height {
		m.Viewport.YOffset = m.cursor - m.Viewport.Height
	}
}

func (m *model) renderCursor() {
	cursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(m.CursorBackground)).
		Foreground(lipgloss.Color(m.CursorForeground))

	// make a deep copy to not preserve cursor position
	cpy := make([]string, len(m.Table.Output))
	copy(cpy, m.Table.Output)

	cpy[m.cursor] = cursorStyle.Render(cpy[m.cursor])
	result := strings.Join(cpy, "\n")
	m.Viewport.SetContent(result)
}

func (m *model) cursorAtViewTop() bool {
	return m.cursor < m.Viewport.YOffset
}

func (m *model) cursorAtViewBottom() bool {
	return m.cursor > m.Viewport.YOffset+m.Viewport.Height-1
}

func (m *model) cursorAtTop() bool {
	return m.cursor < 0
}

func (m *model) cursorAtBottom() bool {
	lastLine := m.Table.LineCount - 1
	return m.cursor > lastLine
}
