package workspaceslist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
)

type WorkspaceListReloadedMsg []clickup.Workspace

type WorkspaceListReadyMsg bool

func WorkspaceListReadyCmd() tea.Cmd {
	return func() tea.Msg {
		return WorkspaceListReadyMsg(true)
	}
}

func (m Model) getWorkspacesCmd() tea.Cmd {
	return func() tea.Msg {
		workspaces, err := m.ctx.Api.GetWorkspaces()
		if err != nil {
			return common.ErrMsg(err)
		}

		return WorkspaceListReloadedMsg(workspaces)
	}
}
