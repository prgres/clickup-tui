package folders

import tea "github.com/charmbracelet/bubbletea"

type LoadingFoldersFromSpaceMsg string

func LoadingFoldersFromSpaceCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return LoadingFoldersFromSpaceMsg(id)
	}
}
