package common

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

func NewEmptyHelp() help.KeyMap {
	return NewHelp(
		func() [][]key.Binding { return [][]key.Binding{} },
		func() []key.Binding { return []key.Binding{} },
	)
}

var KeyBindingBack = key.NewBinding(
	key.WithKeys("escape"),
	key.WithHelp("escape", "back to previous view"),
)

type Help struct {
	fullHelp  func() [][]key.Binding
	shortHelp func() []key.Binding
}

func NewHelp(fullHelp func() [][]key.Binding, shortHelp func() []key.Binding) Help {
	return Help{
		fullHelp:  fullHelp,
		shortHelp: shortHelp,
	}
}

func (km Help) With(kb key.Binding) Help {
	return NewHelp(
		func() [][]key.Binding {
			return append(km.FullHelp(), []key.Binding{kb})
		},
		func() []key.Binding {
			return append(km.ShortHelp(), kb)
		},
	)
}

func (km Help) FullHelp() [][]key.Binding {
	return km.fullHelp()
}

func (km Help) ShortHelp() []key.Binding {
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
