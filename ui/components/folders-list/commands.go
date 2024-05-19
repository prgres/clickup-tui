package folderslist

import tea "github.com/charmbracelet/bubbletea"

type FolderChangeMsg string

func FolderChangeCmd(space string) tea.Cmd {
	return func() tea.Msg {
		return FolderChangeMsg(space)
	}
}
