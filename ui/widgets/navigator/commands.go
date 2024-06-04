package navigator

import (
	tea "github.com/charmbracelet/bubbletea"
	folderslist "github.com/prgrs/clickup/ui/components/folders-list"
	listslist "github.com/prgrs/clickup/ui/components/lists-list"
	spaceslist "github.com/prgrs/clickup/ui/components/spaces-list"
	workspaceslist "github.com/prgrs/clickup/ui/components/workspaces-list"
)

type (
	LoadingSpacesFromWorkspaceMsg string
	LoadingFoldersFromSpaceMsg    string
	LoadingListsFromFolderMsg     string
)

func LoadingSpacesFromWorkspaceCmd(id string) tea.Cmd {
	return func() tea.Msg { return LoadingSpacesFromWorkspaceMsg(id) }
}

func LoadingFoldersFromSpaceCmd(id string) tea.Cmd {
	return func() tea.Msg { return LoadingFoldersFromSpaceMsg(id) }
}

func LoadingListsFromFolderCmd(id string) tea.Cmd {
	return func() tea.Msg { return LoadingListsFromFolderMsg(id) }
}

type (
	ListPreviewMsg  listslist.ListPreviewMsg
	ListChangedMsg  listslist.ListChangedMsg
	ListSelectedMsg listslist.ListSelectedMsg

	FolderPreviewMsg  folderslist.FolderPreviewMsg
	FolderChangedMsg  folderslist.FolderChangedMsg
	FolderSelectedMsg folderslist.FolderSelectedMsg

	SpacePreviewMsg  spaceslist.SpacePreviewMsg
	SpaceChangedMsg  spaceslist.SpaceChangedMsg
	SpaceSelectedMsg spaceslist.SpaceSelectedMsg

	WorkspacePreviewMsg  workspaceslist.WorkspacePreviewMsg
	WorkspaceChangedMsg  workspaceslist.WorkspaceChangedMsg
	WorkspaceSelectedMsg workspaceslist.WorkspaceSelectedMsg
)

func FolderChangedCmd(id string) tea.Cmd {
	return func() tea.Msg { return FolderChangedMsg(id) }
}

func FolderPreviewCmd(id string) tea.Cmd {
	return func() tea.Msg { return FolderPreviewMsg(id) }
}

func ListChangedCmd(id string) tea.Cmd {
	return func() tea.Msg { return ListChangedMsg(id) }
}

func ListPreviewCmd(id string) tea.Cmd {
	return func() tea.Msg { return ListPreviewMsg(id) }
}

func SpaceChangedCmd(id string) tea.Cmd {
	return func() tea.Msg { return SpaceChangedMsg(id) }
}

func SpacePreviewCmd(id string) tea.Cmd {
	return func() tea.Msg { return SpacePreviewMsg(id) }
}

func WorkspaceChangedCmd(workspace string) tea.Cmd {
	return func() tea.Msg { return WorkspaceChangedMsg(workspace) }
}

func WorkspacePreviewCmd(workspace string) tea.Cmd {
	return func() tea.Msg { return WorkspacePreviewMsg(workspace) }
}
