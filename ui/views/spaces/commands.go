package spaces

import tea "github.com/charmbracelet/bubbletea"

type HideSpaceViewMsg bool

func HideSpaceViewCmd() tea.Cmd {
	return func() tea.Msg {
		return HideSpaceViewMsg(true)
	}
}
