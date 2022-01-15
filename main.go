package main

import (
	"fmt"
	"gokeys/pkg"
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

	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FA7921"))
)

type Programs map[string]pkg.Program
type model struct {
	title    string
	programs Programs
	names    []string

	headings   map[int]string
	keyLines   map[int]string
	numOfLines int

	height, width int
	offset        int
	cursor        int

	builder strings.Builder
}

func New(programs Programs) model {
	return model{
		title:    "All your key bindings in one place\n\n",
		programs: programs,
		names:    orderNames(programs), // maintain an ordered slices of names
	}
}

func main() {

	prog := pkg.GetConfig()
	m := New(prog)

	p := tea.NewProgram(&m)
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
		case "ctrl + c", "q":
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

func (m *model) View() string {
	v := fmt.Sprintf("%s %s", m.title, m.viewContent())
	return docStyle.Render(v)
}

func (m *model) viewContent() string {

	n, headings, keyLines := m.orderLines()
	tabbedLines := m.alignKeyLines(keyLines)
	tabbedLinesSlice := strings.Split(tabbedLines, "\n")

	// insert headings
	for i := 1; i < n; i++ {
		if heading, ok := headings[i]; ok {
			tabbedLinesSlice = append(tabbedLinesSlice[:i+1], tabbedLinesSlice[i:]...)
			tabbedLinesSlice[i] = heading
		}
	}

	// insert color
	for i, line := range tabbedLinesSlice {
		if m.cursor == i && strings.Contains(line, " ") {
			tabbedLinesSlice[i] = cursorStyle.Render(line)
		}
	}
	final := strings.Join(tabbedLinesSlice, "\n")
	return final
}

// generate table for key lines
func (m *model) alignKeyLines(lines map[int]string) string {

	defer m.builder.Reset()
	tw := tabwriter.NewWriter(&m.builder, 20, 8, 10, ' ', 0)

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		fmt.Fprintln(tw, line)
	}
	tw.Flush()
	return m.builder.String()
}

func (m *model) updateCursor() {
	if m.cursor < 0 {
		m.cursor = 0
	} else if m.cursor > m.numOfLines {
		m.cursor = m.numOfLines - 1
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

// returns slice of sorted keys
func orderNames(p Programs) []string {
	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (m *model) orderLines() (int, map[int]string, map[int]string) {
	keyLines := make(map[int]string)
	headings := make(map[int]string)
	headingIdx := 1

	for _, name := range m.names {
		p := m.programs[name]
		headings[headingIdx] = titleStyle.Render(name)

		for j, key := range p.KeyBinds {
			line := fmt.Sprintf("%s\t%s", key.Command, key.Key)
			keyLines[headingIdx+j+1] = line
		}
		numOfKeys := len(p.KeyBinds)
		headingIdx += numOfKeys
	}

	m.numOfLines = len(keyLines) + len(headings)
	return len(keyLines) + len(headings), headings, keyLines
}
