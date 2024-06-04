package spaceslist

import tea "github.com/charmbracelet/bubbletea"

type SpaceChangedMsg string

func SpaceChangedCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return SpaceChangedMsg(id)
	}
}
