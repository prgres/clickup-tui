package tasks

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/ui/components/views"
)

type SpaceChangedMsg views.SpaceChangedMsg

func SpaceChangedCmd(space string) tea.Cmd {
	return views.SpaceChangedCmd(space)
}
