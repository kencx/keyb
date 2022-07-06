package list

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit key.Binding

	Up            key.Binding
	Down          key.Binding
	HalfUp        key.Binding
	HalfDown      key.Binding
	FullUp        key.Binding
	FullDown      key.Binding
	GoToFirstLine key.Binding
	GoToLastLine  key.Binding
	GoToTop       key.Binding
	GoToMiddle    key.Binding
	GoToBottom    key.Binding

	CenterCursor key.Binding

	Search key.Binding
	Normal key.Binding
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
		FullUp: key.NewBinding(
			key.WithKeys("ctrl+b"),
		),
		FullDown: key.NewBinding(
			key.WithKeys("ctrl+f"),
		),
		GoToFirstLine: key.NewBinding(
			key.WithKeys("g"),
		),
		GoToLastLine: key.NewBinding(
			key.WithKeys("G"),
		),
		GoToTop: key.NewBinding(
			key.WithKeys("H"),
		),
		GoToMiddle: key.NewBinding(
			key.WithKeys("M"),
		),
		GoToBottom: key.NewBinding(
			key.WithKeys("L"),
		),

		CenterCursor: key.NewBinding(
			key.WithKeys("Z"),
		),

		Search: key.NewBinding(
			key.WithKeys("/"),
		),
		Normal: key.NewBinding(
			key.WithKeys("esc"),
		),
	}
}
