package ui

import (
	"github.com/kencx/keyb/config"
	"github.com/kencx/keyb/ui/list"
	"github.com/kencx/keyb/ui/table"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	List list.Model
	GlobalStyle
}

type GlobalStyle struct {
	CursorForeground string
	CursorBackground string
	Border           string
	BorderColor      string
}

func NewModel(binds config.Bindings, config *config.Config) *Model {
	table := bindingsToTable(binds)

	m := Model{
		List: list.New(config.Title, table),
		GlobalStyle: GlobalStyle{
			CursorForeground: config.Cursor_fg,
			CursorBackground: config.Cursor_bg,
			Border:           config.Border,
			BorderColor:      config.Border_color,
		},
	}
	return &m
}

func bindingsToTable(bindings config.Bindings) *table.Model {
	keys := bindings.SortedKeys()
	parent := appToTable(keys[0], bindings[keys[0]])

	if len(keys) > 1 {
		for _, k := range keys[1:] {
			child := appToTable(k, bindings[k])
			parent.Join(child)
		}
	}
	parent.Align()
	return parent
}

func appToTable(heading string, app config.App) *table.Model {
	var rows []table.Row
	for _, kb := range app.Keybinds {
		rows = append(rows, table.NewRow(kb.Comment, kb.Key, ""))
	}
	h := table.NewHeading(heading)
	return table.New(h, rows)
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
