package spaces

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
)

type SpaceListReloadedMsg []clickup.Space

type SpaceListReadyMsg bool

func SpaceListReadyCmd() tea.Cmd {
	return func() tea.Msg {
		return SpaceListReadyMsg(true)
	}
}

func (m Model) getSpacesCmd() tea.Cmd {
	return func() tea.Msg {
		spaces, err := m.ctx.Api.GetSpaces(m.ctx.Config.DefaultTeam)
		if err != nil {
			return common.ErrMsg(err)
		}

		return SpaceListReloadedMsg(spaces)
	}
}
