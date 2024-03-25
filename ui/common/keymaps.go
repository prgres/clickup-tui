package common

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

func NewEmptyKeyMap() help.KeyMap {
	return NewKeyMap(
		func() [][]key.Binding { return [][]key.Binding{} },
		func() []key.Binding { return []key.Binding{} },
	)
}

var (
	KeyBindingBack = key.NewBinding(
		key.WithKeys("escape"),
		key.WithHelp("escape", "back to previous view"),
	)
	KeyBindingOpenInBrowser = key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "open in browser"),
	)
)

type KeyMap struct {
	fullHelp  func() [][]key.Binding
	shortHelp func() []key.Binding
}

func NewKeyMap(fullHelp func() [][]key.Binding, shortHelp func() []key.Binding) KeyMap {
	return KeyMap{
		fullHelp:  fullHelp,
		shortHelp: shortHelp,
	}
}

func (km KeyMap) With(kb key.Binding) KeyMap {
	return NewKeyMap(
		func() [][]key.Binding {
			return append(km.FullHelp(), []key.Binding{kb})
		},
		func() []key.Binding {
			return append(km.ShortHelp(), kb)
		},
	)
}

func (km KeyMap) FullHelp() [][]key.Binding {
	return km.fullHelp()
}

func (km KeyMap) ShortHelp() []key.Binding {
	return km.shortHelp()
}

func KeyBindingWithHelp(kb key.Binding, description string) key.Binding {
	return key.NewBinding(
		key.WithKeys(kb.Keys()...),
		key.WithHelp(
			strings.ReplaceAll(
				strings.Join(kb.Keys(), ","),
				" ",
				"space",
			),
			description),
	)
}
