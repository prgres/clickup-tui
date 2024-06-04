package folderslist

import tea "github.com/charmbracelet/bubbletea"

type (
	FolderChangedMsg string
	FolderPreviewMsg string
)

func FolderChangedCmd(id string) tea.Cmd {
	return func() tea.Msg { return FolderChangedMsg(id) }
}

func FolderPreviewCmd(id string) tea.Cmd {
	return func() tea.Msg { return FolderPreviewMsg(id) }
}
