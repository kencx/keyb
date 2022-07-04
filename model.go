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

	padding  int // vertical padding - necessary to stabilize scrolling
	cursor   int
	maxWidth int // for word wrapping
	maxRows  int
	Settings
	GlobalStyle
}

type GlobalStyle struct {
	Title            string
	CursorForeground string
	CursorBackground string
	Border           string
	BorderColor      string
}

type Settings struct {
	mouseEnabled bool
}

func NewModel(binds Bindings, config *Config) *model {
	m := model{
		keys:  DefaultKeyMap(),
		debug: true,

		padding:  4,
		maxWidth: 88,

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
		LineStyle:    lipgloss.NewStyle().Margin(0, 2),
	}

	filteredStyles := table.Styles{
		BodyStyle:    lipgloss.NewStyle(),
		HeadingStyle: lipgloss.NewStyle(),
		LineStyle:    lipgloss.NewStyle(),
	}

	if len(binds) == 0 {
		m.table = table.New("", []string{"No key bindings found"}, tableStyles)
		return &m
	}

	m.table = bindingsToTable(binds, tableStyles)
	m.maxRows = m.table.LineCount
	m.filteredTable = table.New("", make([]string, 1, m.table.LineCount), filteredStyles)
	return &m
}

func bindingsToTable(bindings Bindings, style table.Styles) *table.Table {
	keys := bindings.sortedKeys()
	parentTable := appToTable(keys[0], bindings[keys[0]], style)

	if len(keys) > 0 {
		for _, k := range keys[1:] {
			table := appToTable(k, bindings[k], style)
			parentTable.Join(table)
		}
	}

	parentTable.Align()
	return parentTable
}

func appToTable(heading string, app App, styles table.Styles) *table.Table {
	var rows []string
	for _, kb := range app.Keybinds {
		rows = append(rows, kb.String())
	}
	heading = fmt.Sprintf("%s\t%s", heading, " ")
	return table.New(heading, rows, styles)
}
