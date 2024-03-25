package spaces

import tea "github.com/charmbracelet/bubbletea"

type LoadingSpacesFromWorkspaceMsg string

func LoadingSpacesFromWorkspaceCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return LoadingSpacesFromWorkspaceMsg(id)
	}
}
