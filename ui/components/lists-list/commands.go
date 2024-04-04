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
