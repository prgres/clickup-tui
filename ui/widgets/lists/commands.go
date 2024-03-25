package lists

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
)

type ListsListReloadedMsg []clickup.List

func (m Model) getListsCmd(folderId string) tea.Cmd {
	return func() tea.Msg {
		folders, err := m.ctx.Api.GetLists(folderId)
		if err != nil {
			return common.ErrMsg(err)
		}

		return ListsListReloadedMsg(folders)
	}
}

type ListChangedMsg string

func ListChangedCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return ListChangedMsg(id)
	}
}
