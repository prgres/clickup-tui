package listslist

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ListChangedMsg string

func ListChangedCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return ListChangedMsg(id)
	}
}

type ListSelectedMsg string

func ListSelectedCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return ListSelectedMsg(id)
	}
}

type ListPreviewMsg string

func ListPreviewCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return ListPreviewMsg(id)
	}
}
