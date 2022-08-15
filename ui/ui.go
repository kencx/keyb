package ui

import (
	"sort"
	"strings"

	"github.com/kencx/keyb/config"
	"github.com/kencx/keyb/ui/list"
	"github.com/kencx/keyb/ui/table"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	List list.Model
}

func NewModel(a config.Apps, config *config.Config) *Model {

	table := createParentTable(a, config.SortKeys)
	return &Model{
		List: list.New(table, config),
	}
}

func createParentTable(a config.Apps, sortKeys bool) *table.Model {

	if len(a) <= 0 {
		t := table.NewEmpty(1)
		return t
	}

	sort.Slice(a, func(i, j int) bool {
		return a[i].Name < a[j].Name
	})

	parent := appToTable(a[0].Name, *a[0], sortKeys)

	if len(a) > 1 {
		for _, k := range a[1:] {
			child := appToTable(k.Name, *k, sortKeys)
			parent.Join(child)
		}
	}
	return parent
}

func appToTable(heading string, app config.App, sortKeys bool) *table.Model {
	var rows []*table.Row

	h := table.NewHeading(heading)
	rows = append(rows, h)

	if sortKeys {
		sort.Slice(app.Keybinds, func(i, j int) bool {
			return strings.ToLower(app.Keybinds[i].Name) < strings.ToLower(app.Keybinds[j].Name)
		})
	}

	// convert Keybind to Row
	for _, kb := range app.Keybinds {
		row := table.NewRow(kb.Name, kb.Key, app.Prefix, heading)

		// KeyBind's ignore prefix defaults to false
		// so user can choose to ignore prefix for a specific kb
		if kb.IgnorePrefix {
			// Row's show prefix defaults to false
			// so prefix is not shown unless defined
			row.ShowPrefix = false
		}
		rows = append(rows, row)
	}
	return table.New(rows)
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.List.Resize(msg.Width, msg.Height)
	}

	m.List, cmd = m.List.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return m.List.View()
}
