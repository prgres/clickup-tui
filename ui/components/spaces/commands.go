package spaces

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
)

type SpaceListReloadedMsg []clickup.Space

type SpaceListReadyMsg bool

func SpaceListReadyCmd() tea.Cmd {
	return func() tea.Msg {
		return SpaceListReadyMsg(true)
	}
}
