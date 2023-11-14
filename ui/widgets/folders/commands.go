package folders

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
)

type FoldersListReloadedMsg []clickup.Folder

type FoldersListReadyMsg bool

func FoldersListReadyCmd() tea.Cmd {
	return func() tea.Msg {
		return FoldersListReadyMsg(true)
	}
}

func (m Model) getFoldersCmd(space string) tea.Cmd {
	return func() tea.Msg {
		folders, err := m.ctx.Api.GetFolders(space)
		if err != nil {
			return common.ErrMsg(err)
		}

		return FoldersListReloadedMsg(folders)
	}
}
