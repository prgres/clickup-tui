package viewstabs

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	TabChangedMsg  string
	TabSelectedMsg string
)

func TabChangedCmd(id string) tea.Cmd {
	return func() tea.Msg { return TabChangedMsg(id) }
}

func TabSelectedCmd(id string) tea.Cmd {
	return func() tea.Msg { return TabSelectedMsg(id) }
}
