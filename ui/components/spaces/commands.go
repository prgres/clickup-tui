package spaces

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
)

type SpaceChangeMsg string

func SpaceChangeCmd(space string) tea.Cmd {
	return func() tea.Msg {
		return SpaceChangeMsg(space)
	}
}

type TeamChangeMsg string

type SpaceListReloadedMsg []clickup.Space
