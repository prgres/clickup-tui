package folders

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
)

type FoldersListReloadedMsg []clickup.Folder

type FoldersListReadyMsg bool

func FoldersListReadyCmd() tea.Cmd {
	return func() tea.Msg {
		return FoldersListReadyMsg(true)
	}
}
