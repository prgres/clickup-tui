package folderslist

import tea "github.com/charmbracelet/bubbletea"

type (
	FolderChangedMsg  string
	FolderPreviewMsg  string
	FolderSelectedMsg string
)

func FolderChangedCmd(id string) tea.Cmd {
	return func() tea.Msg { return FolderChangedMsg(id) }
}

func FolderPreviewCmd(id string) tea.Cmd {
	return func() tea.Msg { return FolderPreviewMsg(id) }
}

func FolderSelectedCmd(id string) tea.Cmd {
	return func() tea.Msg { return FolderSelectedMsg(id) }
}
