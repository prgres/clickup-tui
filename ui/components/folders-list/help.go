package folderslist

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/prgrs/clickup/ui/common"
)

func (m Model) Help() help.KeyMap {
	return common.NewHelp(
		func() [][]key.Binding {
			return append(
				m.list.FullHelp(),
				[]key.Binding{
					m.keyMap.CursorUp,
					m.keyMap.CursorUpAndSelect,
					m.keyMap.CursorDown,
					m.keyMap.CursorDownAndSelect,
					m.keyMap.Select,
				},
			)
		},
		func() []key.Binding {
			return append(
				m.list.ShortHelp(),
				m.keyMap.CursorUp,
				m.keyMap.CursorDown,
				m.keyMap.Select,
			)
		},
	).With(common.KeyBindingBack)
}
