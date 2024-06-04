package folderslist

import tea "github.com/charmbracelet/bubbletea"

type FolderChangedMsg string

func FolderChangeCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return FolderChangedMsg(id)
	}
}

type FolderSelectedMsg string

func FolderSelectedCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return FolderSelectedCmd(id)
	}
}
