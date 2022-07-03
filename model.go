package main

import (
	"fmt"
	"keyb/table"
	"sort"

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

type model struct {
	Table         *table.Table
	height, width int
	maxWidth      int // for word wrapping
	padding       int // vertical padding - necessary to stabilize scrolling
	offset        int // for vertical scrolling
	cursor        int
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
	mouseDelta   int
}

func NewModel(binds Bindings, config *Config) *model {
	m := model{
		height:   40,
		width:    60,
		maxWidth: 88,
		padding:  4,
		Settings: Settings{
			mouseEnabled: true,
			mouseDelta:   3,
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

	if len(binds) == 0 {
		m.Table = table.New("", []string{""}, tableStyles)
		m.Table.Output = []string{"No key bindings found"}
		return &m
	}

	m.Table = bindingsToTable(binds, tableStyles)
	return &m
}

func bindingsToTable(bindings Bindings, style table.Styles) *table.Table {
	keys := bindings.sortedKeys()
	parentTable := appToTable(keys[0], bindings[keys[0]], style)

	for _, k := range keys[1:] {
		table := appToTable(k, bindings[k], style)
		parentTable.Join(table)
	}

	parentTable.Align()
	return parentTable
}

func appToTable(heading string, app App, styles table.Styles) *table.Table {
	var rows []string
	for _, kb := range app.Keybinds {
		rows = append(rows, fmt.Sprintf("%s\t%s", kb.Comment, kb.Key))
	}
	heading = fmt.Sprintf("%s\t%s", heading, " ")
	return table.New(heading, rows, styles)
}
