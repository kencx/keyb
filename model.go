package main

import (
	"fmt"
	"keyb/table"
	"sort"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
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

type filterState int

const (
	unfiltered filterState = iota
	filtering
)

type model struct {
	keys     KeyMap
	viewport viewport.Model
	table    *table.Table

	ready bool
	debug bool

	search    bool
	searchBar textinput.Model

	filterState   filterState
	filteredTable *table.Table

	padding int // vertical padding - necessary to stabilize scrolling
	cursor  int
	// maxWidth int // for word wrapping
	maxRows int
	Settings
	GlobalStyle
}

type Settings struct {
	mouseEnabled bool
}

type GlobalStyle struct {
	Title            string
	CursorForeground string
	CursorBackground string
	Border           string
	BorderColor      string
}

func NewModel(binds Bindings, config *Config) *model {
	m := model{
		keys:  DefaultKeyMap(),
		debug: true,

		padding: 4,
		// maxWidth: 88,

		Settings: Settings{
			mouseEnabled: true,
		},
		GlobalStyle: GlobalStyle{
			Title:            config.Title,
			CursorForeground: config.Cursor_fg,
			CursorBackground: config.Cursor_bg,
			Border:           config.Border,
			BorderColor:      config.Border_color,
		},
	}

	tableStyles := table.Styles{
		BodyStyle:    lipgloss.NewStyle(),
		HeadingStyle: lipgloss.NewStyle().Bold(true).Margin(0, 1),
		RowStyle:     lipgloss.NewStyle().Margin(0, 2),
	}

	if len(binds) == 0 {
		m.table = table.NewEmpty(1)
		m.table.AppendRow("No key bindings found")
		return &m
	}

	m.table = bindingsToTable(binds, tableStyles)
	m.maxRows = m.table.LineCount
	m.filteredTable = table.NewEmpty(m.table.LineCount)
	m.filteredTable.Styles = tableStyles
	return &m
}

func bindingsToTable(bindings Bindings, style table.Styles) *table.Table {
	keys := bindings.sortedKeys()
	parent := appToTable(keys[0], bindings[keys[0]], style)

	if len(keys) > 1 {
		for _, k := range keys[1:] {
			child := appToTable(k, bindings[k], style)
			parent.Join(child)
		}
	}
	parent.Align()
	return parent
}

func appToTable(heading string, app App, styles table.Styles) *table.Table {
	var rows []string
	for _, kb := range app.Keybinds {
		rows = append(rows, kb.String())
	}
	heading = fmt.Sprintf("%s\t%s", heading, " ")
	return table.New(heading, rows, styles)
}
