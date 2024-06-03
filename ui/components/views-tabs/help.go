package viewstabs

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/prgrs/clickup/ui/common"
)

func (m Model) Help() help.KeyMap {
	return common.NewHelp(
		func() [][]key.Binding {
			return [][]key.Binding{
				{
					m.keyMap.CursorLeft,
					m.keyMap.CursorLeftAndSelect,
					m.keyMap.CursorRight,
					m.keyMap.CursorRightAndSelect,
					m.keyMap.Select,
				},
			}
		},
		func() []key.Binding {
			return []key.Binding{
				m.keyMap.CursorLeft,
				m.keyMap.CursorRight,
				m.keyMap.Select,
			}
		},
	)
}
