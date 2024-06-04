package folderslist

import tea "github.com/charmbracelet/bubbletea"

type FolderChangedMsg string

func FolderChangeCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return FolderChangedMsg(id)
	}
}

