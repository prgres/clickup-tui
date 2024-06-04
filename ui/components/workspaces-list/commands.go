package workspaceslist

import tea "github.com/charmbracelet/bubbletea"

type WorkspaceChangedMsg string

func WorkspaceChangedCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return WorkspaceChangedMsg(id)
	}
}
