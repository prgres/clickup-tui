package spaces

import tea "github.com/charmbracelet/bubbletea"

type SpaceChangeMsg string

func SpaceChangeCmd(space string) tea.Cmd {
	return func() tea.Msg {
		return SpaceChangeMsg(space)
	}
}

type HideSpaceViewMsg bool

func HideSpaceViewCmd() tea.Cmd {
	return func() tea.Msg {
		return HideSpaceViewMsg(true)
	}
}
