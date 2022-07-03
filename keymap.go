package main

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit key.Binding

	Up         key.Binding
	Down       key.Binding
	HalfUp     key.Binding
	HalfDown   key.Binding
	GoToTop    key.Binding
	GoToBottom key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
		),
		HalfUp: key.NewBinding(
			key.WithKeys("ctrl+u"),
		),
		HalfDown: key.NewBinding(
			key.WithKeys("ctrl+d"),
		),
		GoToTop: key.NewBinding(
			key.WithKeys("g"),
		),
		GoToBottom: key.NewBinding(
			key.WithKeys("G"),
		),
	}
}
