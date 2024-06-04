package navigator

import (
	tea "github.com/charmbracelet/bubbletea"
	listslist "github.com/prgrs/clickup/ui/components/lists-list"
)

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

type (
	ListPreviewMsg  = listslist.ListPreviewMsg
	ListChangedMsg  = listslist.ListChangedMsg
	ListSelectedMsg = listslist.ListSelectedMsg
)

var (
	ListChangedCmd  = listslist.ListChangedCmd
	ListSelectedCmd = listslist.ListSelectedCmd
)
