package main

import (
	"fmt"
	"gokeys/pkg"
	"strings"
	"text/tabwriter"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	Title string

	width, height int
	cursor        int

	Apps       map[string]pkg.App
	numOfApps  int
	numOfLines int
}

func New(apps map[string]pkg.App) model {
	return model{
		Apps:       apps,
		numOfApps:  len(apps),
		numOfLines: countLines(apps),
	}
}

func countLines(apps map[string]pkg.App) int {
	var result int
	for _, app := range apps {
		result += len(app.KeyBinds) + 1 // +1 for app name
	}
	return result
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "ctrl + c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.Apps)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	m.Title = "Key bindings\n\n"
	s := m.Title

	for name, app := range m.Apps {
		str, err := appView(name, app)
		if err != nil {
			fmt.Println(err)
		}
		s += str
		s += "\n"
	}
	s += "\nPress q to Quit\n"
	return s
}

func appView(name string, app pkg.App) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintln(name))
	w := tabwriter.NewWriter(&sb, 1, 1, 1, ' ', 0)

	for _, key := range app.KeyBinds {
		fmt.Fprintf(w, "%s\t%s\n", key.Command, key.Key)
	}

	if err := w.Flush(); err != nil {
		return "", fmt.Errorf("error: %w", err)
	}
	return sb.String(), nil
}

func main() {

	appMap := pkg.GetConfig()
	m := New(appMap)

	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Errorf("unexpected error: %w", err)
	}
}
