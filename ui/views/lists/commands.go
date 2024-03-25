package lists

import tea "github.com/charmbracelet/bubbletea"

type LoadingListsFromFolderMsg string

func LoadingFoldersFromSpaceCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return LoadingListsFromFolderMsg(id)
	}
}
