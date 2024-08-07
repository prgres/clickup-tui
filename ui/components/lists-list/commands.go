package listslist

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	ListChangedMsg  string
	ListPreviewMsg  string
	ListSelectedMsg string
)

func ListChangedCmd(id string) tea.Cmd {
	return func() tea.Msg { return ListChangedMsg(id) }
}

func ListPreviewCmd(id string) tea.Cmd {
	return func() tea.Msg { return ListPreviewMsg(id) }
}

func ListSelectedCmd(id string) tea.Cmd {
	return func() tea.Msg { return ListSelectedMsg(id) }
}
