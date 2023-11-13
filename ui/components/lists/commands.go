package lists

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
)

type ListsListReloadedMsg []clickup.List

type ListsListReadyMsg bool

func ListsListReadyCmd() tea.Cmd {
	return func() tea.Msg {
		return ListsListReadyMsg(true)
	}
}
