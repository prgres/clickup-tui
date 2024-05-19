package navigator

import tea "github.com/charmbracelet/bubbletea"

type LoadingSpacesFromWorkspaceMsg string

func LoadingSpacesFromWorkspaceCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return LoadingSpacesFromWorkspaceMsg(id)
	}
}

type LoadingFoldersFromSpaceMsg string

func LoadingFoldersFromSpaceCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return LoadingFoldersFromSpaceMsg(id)
	}
}

type LoadingListsFromFolderMsg string

func LoadingListsFromFolderCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return LoadingListsFromFolderMsg(id)
	}
}
