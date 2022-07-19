package list

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/kencx/keyb/config"
)

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

	Search      key.Binding
	ClearSearch key.Binding
	Normal      key.Binding
}

func CreateKeyMap(keys config.Keys) KeyMap {
	return KeyMap{
		Quit:          SetKey(keys.Quit),
		Up:            SetKey(keys.Up),
		Down:          SetKey(keys.Down),
		HalfUp:        SetKey(keys.HalfUp),
		HalfDown:      SetKey(keys.HalfDown),
		FullUp:        SetKey(keys.FullUp),
		FullDown:      SetKey(keys.FullDown),
		GoToFirstLine: SetKey(keys.GoToFirstLine),
		GoToLastLine:  SetKey(keys.GoToLastLine),
		GoToTop:       SetKey(keys.GoToTop),
		GoToMiddle:    SetKey(keys.GoToMiddle),
		GoToBottom:    SetKey(keys.GoToBottom),

		Search:      SetKey(keys.Search),
		ClearSearch: SetKey(keys.ClearSearch),
		Normal:      SetKey(keys.Normal),
	}
}

func SetKey(s string) key.Binding {
	return key.NewBinding(
		key.WithKeys(splitAndTrim(s, ",")...),
	)
}

func splitAndTrim(s, sep string) []string {
	sl := strings.Split(s, sep)
	for i := range sl {
		sl[i] = strings.TrimSpace(sl[i])
	}
	return sl
}
