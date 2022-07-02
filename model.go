package main

import (
	"keyb/table"

	"github.com/charmbracelet/lipgloss"
)

type Bindings map[string]App

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
}

type Settings struct {
	mouseEnabled bool
	mouseDelta   int
}

func NewModel(binds Bindings, config *Config) *model {
	header := "bspwm"
	keys := binds[header].Keybinds
	var rows []string

	for _, k := range keys {
		rows = append(rows, k.Key)
	}

	styles := table.Styles{
		BodyStyle:        lipgloss.NewStyle(),
		HeadingStyle:     lipgloss.NewStyle().Bold(true).Margin(0, 1),
		LineStyle:        lipgloss.NewStyle().Margin(0, 2),
		Title:            config.Title,
		CursorForeground: config.Cursor_fg,
		CursorBackground: config.Cursor_bg,
		Border:           config.Border,
		BorderColor:      config.Border_color,
	}

	table := table.New(header, rows, styles)
	table.Render()

	s := Settings{
		mouseEnabled: true,
		mouseDelta:   3,
	}

	m := model{
		Table:    table,
		height:   40,
		width:    60,
		maxWidth: 88,
		padding:  4,
		Settings: s,
	}
	return &m
}
