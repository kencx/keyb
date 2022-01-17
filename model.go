package main

import (
	"fmt"
	"strings"
	"text/tabwriter"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true)
	lineStyle  = lipgloss.NewStyle().Margin(0, 1)
	docStyle   = lipgloss.NewStyle().Margin(1, 4)
)

type Categories map[string]Program
type model struct {
	categories Categories // map of programs from config file
	headings   []string   // *ordered* slice of category names

	headingMap map[int]string // map of headings with line no.
	lineMap    map[int]string // map of all key bindings with line no.
	lineCount  int            // total number of lines

	// body is the final output split into lines. It is split to update the cursor
	// for each line.
	body []string

	height, width int
	maxWidth      int // for word wrapping
	padding       int // bottom padding
	cursor        int
	offset        int

	// styling
	title string
	curFg string
	curBg string

	mouseEnabled bool
	mouseDelta   int
}

func New(cat Categories, config map[string]string) *model {
	m := model{
		categories: cat,
		headings:   sortKeys(cat), // ordered slices of names

		height:   40,
		width:    60,
		maxWidth: 88,
		padding:  4,

		title: config["title"],
		curFg: config["curFg"],
		curBg: config["curBg"],

		mouseEnabled: true,
		mouseDelta:   3,
	}
	if len(m.headings) > 0 {
		m.initBody()
	}
	return &m
}

func (m *model) initBody() {
	m.headingMap, m.lineMap, m.lineCount = m.splitHeadingsAndKeys()

	// generate properly aligned table
	keyTable := m.alignKeyTable() // TODO word wrapped lines do not align well

	// split body into lines to insert headings
	m.body = strings.Split(keyTable, "\n")
	m.insertHeadings()
}

// generate 2 maps from categories: unbtabbed headings and tabbed lines
func (m *model) splitHeadingsAndKeys() (map[int]string, map[int]string, int) {

	headingMap := make(map[int]string)
	lineMap := make(map[int]string)

	var (
		headingIdx       int
		totalLineCount   int // track total number of lines
		wrappedLineCount int // account for wrapping of lines > maxWidth
	)

	for _, h := range m.headings {
		headingMap[headingIdx] = titleStyle.Render(h)

		p := m.categories[h]
		for i, key := range p.KeyBinds {

			// handle word wrapping for long lines
			if len(key.Command) >= m.maxWidth {
				wrappedLineCount = (len(key.Command) / m.maxWidth) + 1
				key.Command = wordwrap.String(key.Command, min(m.width, m.maxWidth))
			}

			line := fmt.Sprintf("%s\t%s", key.Command, key.Key)
			lineMap[headingIdx+i+1] = lineStyle.Render(line)
		}
		// total num of keys + heading + num of (extra) wrapped lines
		totalLineCount += len(p.KeyBinds) + 1 + wrappedLineCount

		headingIdx += len(p.KeyBinds) + max(1, wrappedLineCount)
		wrappedLineCount = 0 // reset
	}
	return headingMap, lineMap, totalLineCount
}

// generate properly aligned table
func (m *model) alignKeyTable() string {

	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 20, 8, 10, ' ', 0)

	// order keys by line no.
	for i := 1; i < m.lineCount; i++ {
		if line, ok := m.lineMap[i]; ok {
			fmt.Fprintln(tw, line)
		}
	}
	tw.Flush()
	return sb.String()
}

// insert headings at respective line numbers
func (m *model) insertHeadings() {
	for i := 0; i < m.lineCount; i++ {
		if heading, ok := m.headingMap[i]; ok {
			m.body = insertAtIndex(i, heading, m.body)
		}
	}
}

// bubbletea
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
			m.cursor = m.lineCount - 1
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
		m.cursor = 0

	} else if m.cursor >= m.lineCount-1 {
		m.cursor = m.lineCount - 1
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
	view := fmt.Sprintf("%s\n\n%s", m.title, m.updateBody())
	return docStyle.Render(view)
}

func (m *model) updateBody() string {

	if len(m.body) <= 0 {
		return "No key bindings found"
	}

	// deep copy slice
	cpy := make([]string, len(m.body))
	copy(cpy, m.body)

	body := m.renderCursor(cpy)

	if m.lineCount >= m.offset+m.height {
		body = body[m.offset : m.offset+m.height]
	}
	result := strings.Join(body, "\n")
	return result
}

// render cursor style at position
func (m *model) renderCursor(lines []string) []string {
	cursorStyle := lipgloss.NewStyle().Background(lipgloss.Color(m.curBg)).Foreground(lipgloss.Color(m.curFg))

	for i, line := range lines {
		if m.cursor == i {
			lines[i] = cursorStyle.Render(line)
		}
	}
	return lines
}
