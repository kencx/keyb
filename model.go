package main

import (
	"fmt"
	"strings"
	"text/tabwriter"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type model struct {
	categories Categories // map of programs from keyb file
	headings   []string   // *ordered* slice of category names

	headingMap map[int]string // map of headings with line no.
	lineMap    map[int]string // map of all key bindings with line no.
	lineCount  int            // total number of lines

	// body is the final output split into lines. It is split to update the cursor
	// in updateBody()
	body []string

	height, width int
	maxWidth      int // for word wrapping
	padding       int // vertical padding - necessary to stabilize scrolling
	offset        int // for vertical scrolling
	cursor        int

	// customization
	bodyStyle   lipgloss.Style
	title       string
	curFg       string
	curBg       string
	border      string
	borderColor string

	mouseEnabled bool
	mouseDelta   int
}

func New(cat Categories, config *Config) *model {

	m := model{
		categories: cat,
		headings:   sortKeys(cat),

		height:   40,
		width:    60,
		maxWidth: 88,
		padding:  4,

		bodyStyle:   lipgloss.NewStyle(),
		title:       config.Title,
		curFg:       config.Cursor_fg,
		curBg:       config.Cursor_bg,
		border:      config.Border,
		borderColor: config.Border_color,

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
	keyTable := m.alignKeyTable() // TODO word wrapped lines do not align well

	m.body = strings.Split(keyTable, "\n")
	m.insertHeadings()
}

// generate 2 maps from categories: unbtabbed headings and tabbed lines
func (m *model) splitHeadingsAndKeys() (map[int]string, map[int]string, int) {

	var (
		headingStyle = lipgloss.NewStyle().Bold(true).Margin(0, 1)
		lineStyle    = lipgloss.NewStyle().Margin(0, 2)

		headingIdx     int
		totalLineCount int // track total number of lines
		line           string
	)

	headingMap := make(map[int]string)
	lineMap := make(map[int]string)

	for _, h := range m.headings {
		var wrappedLineCount int // account for wrapping of lines > maxWidth
		headingMap[headingIdx] = headingStyle.Render(h)

		p := m.categories[h]
		for i, key := range p.KeyBinds {

			// handle word wrapping for long lines
			if len(key.Desc) >= m.maxWidth {
				wrappedLineCount = (len(key.Desc) / m.maxWidth) + 1
				key.Desc = wordwrap.String(key.Desc, min(m.width, m.maxWidth))
			}

			if p.Prefix != "" && !key.Ignore_Prefix {
				line = fmt.Sprintf("%s\t%s ; %s", key.Desc, p.Prefix, key.Key)
			} else {
				line = fmt.Sprintf("%s\t%s", key.Desc, key.Key)
			}
			lineMap[headingIdx+i+1] = lineStyle.Render(line)
		}
		// each category contributes: num of keys + heading + num of (extra) wrapped lines
		totalLineCount += len(p.KeyBinds) + 1 + wrappedLineCount

		// required offset
		headingIdx += len(p.KeyBinds) + max(1, wrappedLineCount)
	}
	return headingMap, lineMap, totalLineCount
}

// generate aligned table
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
		m.cursor = m.lineCount - 1

	} else if m.cursor > m.lineCount-1 {
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
	view := fmt.Sprintf("%s\n\n%s", m.title, body)

	m.setStyle()
	return m.bodyStyle.Render(view)

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
	cursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(m.curBg)).
		Foreground(lipgloss.Color(m.curFg))

	for i, line := range lines {
		if m.cursor == i {
			lines[i] = cursorStyle.Render(line)
		}
	}
	return lines
}

func (m *model) setStyle() {
	m.bodyStyle = m.bodyStyle.Margin(1, 2)
	m.handleBorder()
}

// issues with border:
// does not resize with window size
// width changes when lines shrink and grow
func (m *model) handleBorder() {
	var borderStyle lipgloss.Border

	switch m.border {
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

	m.bodyStyle = m.bodyStyle.Border(borderStyle)
	m.padding += m.bodyStyle.GetBorderTopWidth() + m.bodyStyle.GetBorderBottomSize()
}
