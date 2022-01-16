package main

import (
	"fmt"
	"keyb/pkg"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// vertical padding & margin ruins scrolling for some reason
	titleStyle  = lipgloss.NewStyle().Bold(true)
	lineStyle   = lipgloss.NewStyle().Margin(0, 1)
	docStyle    = lipgloss.NewStyle().Margin(0, 4)
	cursorStyle = lipgloss.NewStyle().Background(lipgloss.Color("#448488"))
)

type (
	Categories map[string]pkg.Program

	model struct {
		title      string
		categories Categories // map of programs from config file
		headings   []string   // *ordered* slice of category names

		headingMap map[int]string // map of headings with line no.
		lineMap    map[int]string // map of all key bindings with line no.
		lineCount  int            // total number of lines

		// body is the final output split into lines. It is split to update the cursor
		// for each line.
		body []string

		height, width int
		padding       int
		cursor        int
		offset        int
	}
)

func New(cat Categories) *model {
	m := model{
		title:      "All your key bindings in one place\n",
		categories: cat,
		headings:   sortKeys(cat), // maintain an ordered slices of names
		height:     40,
		width:      60,
		padding:    2,
	}

	m.headingMap, m.lineMap = m.splitHeadingsAndKeys()
	m.lineCount = len(m.headingMap) + len(m.lineMap)

	// generate properly aligned table
	keyTable := m.alignKeyTable()

	// split body into lines to insert headings
	m.body = strings.Split(keyTable, "\n")
	m.insertHeadings()
	return &m
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

// generate 2 maps from categories: unbtabbed headings and tabbed lines
func (m *model) splitHeadingsAndKeys() (map[int]string, map[int]string) {

	headingMap := make(map[int]string)
	lineMap := make(map[int]string)
	headingIdx := 0

	for _, h := range m.headings {
		headingMap[headingIdx] = titleStyle.Render(h)
		p := m.categories[h]

		for i, key := range p.KeyBinds {
			lineMap[headingIdx+i+1] = lineStyle.Render(fmt.Sprintf("%s\t%s", key.Command, key.Key))
		}
		headingIdx += len(p.KeyBinds) + 1
	}
	return headingMap, lineMap
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
		case "up", "k":
			m.cursor--
			m.updateCursor()
		case "down", "j":
			m.cursor++
			m.updateCursor()
		case "G":
			m.cursor = m.lineCount - 1
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height - m.padding
		m.offset = 0
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
	view := fmt.Sprintf("%s\n%s", m.title, m.setBody())
	return docStyle.Render(view)
}

func (m *model) setBody() string {
	// deep copy slice
	cpy := make([]string, len(m.body))
	copy(cpy, m.body)

	body := renderCursor(cpy, m.cursor)

	if m.lineCount >= m.offset+m.height {
		body = body[m.offset : m.offset+m.height]
	}
	result := strings.Join(body, "\n")
	return result
}

// render cursor style at position
func renderCursor(lines []string, cursorPosition int) []string {
	for i, line := range lines {
		if cursorPosition == i {
			lines[i] = cursorStyle.Render(line)
		}
	}
	return lines
}

func insertAtIndex(index int, element string, array []string) []string {
	array = append(array[:index+1], array[index:]...)
	array[index] = element
	return array
}

// output
func (m *model) outputBodyToFile(path string) error {

	output := strings.Join(m.body, "\n")
	if err := os.WriteFile(path, []byte(output), 0744); err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

func (m *model) outputBodyToStdout() error {

	output := strings.Join(m.body, "\n")
	_, err := os.Stdout.Write([]byte(output))
	if err != nil {
		return fmt.Errorf("error writing to stdout: %w", err)
	}
	return nil
}
