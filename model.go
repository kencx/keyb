package main

import (
	"fmt"
	"keyb/list"
	"keyb/table"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
)

type Bindings map[string]App

func (b Bindings) sortedKeys() []string {
	keys := make([]string, 0, len(b))
	for k := range b {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// TODO replace bindings with []App
// name of App should be a field instead
type App struct {
	Prefix   string    `yaml:"prefix,omitempty"`
	Keybinds []KeyBind `yaml:"keybinds"`
}

type KeyBind struct {
	Comment      string `yaml:"desc"`
	Key          string `yaml:"key"`
	IgnorePrefix bool   `yaml:"ignore_prefix,omitempty"`
}

func (k KeyBind) String() string {
	return fmt.Sprintf("%s\t%s", k.Comment, k.Key)
}

type model struct {
	list list.Model
	GlobalStyle
}

type GlobalStyle struct {
	CursorForeground string
	CursorBackground string
	Border           string
	BorderColor      string
}

func NewModel(binds Bindings, config *Config) *model {
	table := bindingsToTable(binds)

	m := model{
		list: list.New(config.Title, table),
		GlobalStyle: GlobalStyle{
			CursorForeground: config.Cursor_fg,
			CursorBackground: config.Cursor_bg,
			Border:           config.Border,
			BorderColor:      config.Border_color,
		},
	}
	return &m
}

func bindingsToTable(bindings Bindings) *table.Model {
	keys := bindings.sortedKeys()
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

func appToTable(heading string, app App) *table.Model {
	var rows []string
	for _, kb := range app.Keybinds {
		rows = append(rows, kb.String())
	}
	heading = fmt.Sprintf("%s\t%s", heading, " ")
	return table.New(heading, rows)
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.Resize(msg.Width, msg.Height)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *model) View() string {
	return m.list.View()
}
