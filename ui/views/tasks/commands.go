package tasks

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/ui/components/viewtabs"
)

type SpaceChangedMsg viewtabs.SpaceChangedMsg

func SpaceChangedCmd(space string) tea.Cmd {
	return viewtabs.SpaceChangedCmd(space)
}
