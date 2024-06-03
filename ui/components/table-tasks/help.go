package tabletasks

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
					km.RowDown,
					km.RowUp,
					km.RowSelectToggle,
				},
				{
					km.PageDown,
					km.PageUp,
					km.PageFirst,
					km.PageLast,
				},
				{
					km.Filter,
					km.FilterBlur,
					km.FilterClear,
				},
				{
					km.ScrollRight,
					km.ScrollLeft,
					m.keyMap.Select,
				},
			}
		},
		func() []key.Binding {
			return []key.Binding{
				km.RowDown,
				km.RowUp,
				km.RowSelectToggle,
				km.PageDown,
				km.PageUp,
				m.keyMap.Select,
			}
		},
	)
}
