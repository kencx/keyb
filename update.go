package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	lastLine := m.Table.LineCount - 1

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			m.cursor--
			if m.cursorAtViewTop() {
				m.Viewport.LineUp(1)
			}
		case key.Matches(msg, m.keys.Down):
			m.cursor++
			if m.cursorAtViewBottom() {
				m.Viewport.LineDown(1)
			}

		case key.Matches(msg, m.keys.GoToTop):
			m.cursor = 0
			m.Viewport.GotoTop()
		case key.Matches(msg, m.keys.GoToBottom):
			m.cursor = lastLine
			m.Viewport.GotoBottom()

		case key.Matches(msg, m.keys.HalfUp):
			m.cursor -= m.Viewport.Height / 2
			if m.cursorAtViewTop() {
				m.Viewport.HalfViewUp()
			}

			// don't loop around
			if m.cursor < 0 {
				m.cursor = 0
				m.Viewport.GotoTop()
			}
		case key.Matches(msg, m.keys.HalfDown):
			m.cursor += m.Viewport.Height / 2
			if m.cursorAtViewBottom() {
				m.Viewport.HalfViewDown()
			}

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
	}

	// cursor loop around
	if m.cursorAtTop() {
		m.cursor = lastLine
		m.Viewport.GotoBottom()
	} else if m.cursorAtBottom() {
		m.cursor = 0
		m.Viewport.GotoTop()
	}

	m.renderCursor()
	return m, nil
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
