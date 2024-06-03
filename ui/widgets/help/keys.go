package help

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	ShowHelp key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		ShowHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "show help"),
		),
	}
}

func (m *Model) handleKeys(msg tea.KeyMsg) tea.Cmd {
	m.lastKey = msg.String()

	switch {
	case key.Matches(msg, m.keyMap.ShowHelp):
		m.ShowHelp = !m.ShowHelp
		m.help.ShowAll = !m.help.ShowAll
	}

	return nil
}
