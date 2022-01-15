package main

import (
	"fmt"
	"keyb/pkg"
	"log"
	"sort"
	"strings"
	"text/tabwriter"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().
			Margin(1, 4).
			Padding(1, 2).
			Width(88).
			Height(20).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#ebdbb2"))

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			MarginTop(1)

	cursorStyle = lipgloss.NewStyle().Background(lipgloss.Color("#448488"))
)

type (
	Categories map[string]pkg.Program

	model struct {
		title      string     // title text
		categories Categories // map of programs from config file
		headings   []string   // *ordered* slice of category names

		headingMap map[int]string // map of headings with line no.
		lineMap    map[int]string // map of all key bindings with line no.
		lineCount  int            // total number of lines
		body       []string

		height, width int
		offset        int
		cursor        int

		// builder strings.Builder
	}
)

func New(cat Categories) *model {
	m := model{
		title:      "All your key bindings in one place\n",
		categories: cat,
		headings:   sortKeys(cat), // maintain an ordered slices of names
		cursor:     1,
	}

	m.headingMap, m.lineMap = m.splitHeadingsAndKeys()
	m.lineCount = len(m.headingMap) + len(m.lineMap) - 1

	keyTable := alignKeyTable(m.lineMap) // generate properly aligned table

	// split body into lines to insert headings
	m.body = strings.Split(keyTable, "\n")
	m.insertHeadings()

	return &m
}

// generate 2 maps from categories: unbtabbed headings and tabbed lines
func (m *model) splitHeadingsAndKeys() (map[int]string, map[int]string) {

	headingMap := make(map[int]string)
	lineMap := make(map[int]string)
	headingIdx := 1

	for _, h := range m.headings {
		headingMap[headingIdx] = titleStyle.Render(h)
		p := m.categories[h]

		for i, key := range p.KeyBinds {
			lineMap[headingIdx+i+1] = fmt.Sprintf("%s\t%s", key.Command, key.Key)
		}
		headingIdx += len(p.KeyBinds)
	}
	return headingMap, lineMap
}

// insert headings at respective line numbers
func (m *model) insertHeadings() {
	for i := 1; i < m.lineCount; i++ {
		if heading, ok := m.headingMap[i]; ok {
			m.body = insertAtIndex(i, heading, m.body)
		}
	}
}

func main() {
	prog := pkg.GetConfig()
	m := New(prog)

	p := tea.NewProgram(m)
	p.EnterAltScreen()
	defer p.ExitAltScreen()

	if err := p.Start(); err != nil {
		log.Fatalln(err)
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
		case "up":
			m.cursor--
			m.updateCursor()
			m.updateOffset()
		case "down":
			m.cursor++
			m.updateCursor()
			m.updateOffset()
		}
		// case tea.WindowSizeMsg:
		// 	m.width = msg.Width
		// 	m.height = msg.Height
	}
	return m, nil
}

func (m *model) updateCursor() {
	if m.cursor < 0 {
		m.cursor = 0
	} else if m.cursor > m.lineCount {
		m.cursor = m.lineCount - 1
	}
}

func (m *model) updateOffset() {
	// scroll down
	if m.cursor >= m.offset+m.height {
		m.offset = m.cursor - m.height + 1
	}
	// scroll up
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
}

func (m *model) View() string {
	v := fmt.Sprintf("%s %s", m.title, m.generateBody())
	return docStyle.Render(v)
}

func (m *model) generateBody() string {

	// deep copy slice
	cpy := make([]string, len(m.body))
	copy(cpy, m.body)

	renderedBody := renderCursor(cpy, m.cursor)
	body := strings.Join(renderedBody, "\n")
	return body
}

// render cursor style at position
func renderCursor(lines []string, cursorPosition int) []string {
	for i, line := range lines {
		if cursorPosition == i && strings.Contains(line, " ") {
			lines[i] = cursorStyle.Render(line)
		}
	}
	return lines
}

// returns slice of sorted keys
func sortKeys(cat Categories) []string {
	keys := make([]string, 0, len(cat))
	for k := range cat {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// generate properly aligned table
func alignKeyTable(lineMap map[int]string) string {

	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 20, 8, 10, ' ', 0)

	// order keys by line no.
	for i := 1; i < len(lineMap); i++ {
		line := lineMap[i]
		fmt.Fprintln(tw, line)
	}
	tw.Flush()
	return sb.String()
}

func insertAtIndex(index int, element string, array []string) []string {
	array = append(array[:index+1], array[index:]...)
	array[index] = element
	return array
}
