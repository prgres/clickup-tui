package tasks

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/prgrs/clickup/ui/common"
)

func (m Model) Help() help.KeyMap {
	var help help.KeyMap

	if m.copyMode {
		return common.NewHelp(
			func() [][]key.Binding {
				return [][]key.Binding{
					{
						m.keyMap.CopyTaskId,
						m.keyMap.CopyTaskUrl,
						m.keyMap.CopyTaskUrlMd,
						m.keyMap.LostFocus,
						m.keyMap.Refresh,
					},
				}
			},
			func() []key.Binding {
				return []key.Binding{
					m.keyMap.CopyTaskId,
					m.keyMap.CopyTaskUrl,
					m.keyMap.CopyTaskUrlMd,
					m.keyMap.LostFocus,
					m.keyMap.Refresh,
				}
			},
		)
	}

	if m.editMode {
		return common.NewHelp(
			func() [][]key.Binding {
				return [][]key.Binding{
					{
						m.keyMap.EditDescription,
						m.keyMap.EditName,
						m.keyMap.EditStatus,
						m.keyMap.EditAssigness,
						m.keyMap.EditQuit,
					},
				}
			},
			func() []key.Binding {
				return []key.Binding{
					m.keyMap.EditDescription,
					m.keyMap.EditName,
					m.keyMap.EditStatus,
					m.keyMap.EditAssigness,
					m.keyMap.EditQuit,
				}
			},
		)
	}
	switch m.state {
	case m.componenetTasksSidebar.Id():
		help = m.componenetTasksSidebar.Help()
	case m.componenetTasksTable.Id():
		help = m.componenetTasksTable.Help()
	}

	return common.NewHelp(
		func() [][]key.Binding {
			return append(
				help.FullHelp(),
				[]key.Binding{
					m.keyMap.OpenTicketInWebBrowser,
					m.keyMap.ToggleSidebar,
					m.keyMap.EditMode,
				},
			)
		},
		func() []key.Binding {
			return append(
				help.ShortHelp(),
				m.keyMap.OpenTicketInWebBrowser,
				m.keyMap.CopyMode,
				m.keyMap.EditMode,
			)
		},
	)
}
