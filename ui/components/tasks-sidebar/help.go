package taskssidebar

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/prgrs/clickup/ui/common"
)

func (m Model) Help() help.KeyMap {
	km := m.keyMap

	return common.NewHelp(
		func() [][]key.Binding {
			return [][]key.Binding{
				{
					km.Down,
					km.Up,
				},
				{
					km.PageDown,
					km.PageUp,
				},
				{
					km.HalfPageUp,
					km.HalfPageDown,
				},
			}
		},
		func() []key.Binding {
			return []key.Binding{
				km.Down,
				km.Up,
				km.PageDown,
				km.PageUp,
			}
		},
	)
}
