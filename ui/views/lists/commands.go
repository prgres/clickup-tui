package lists

import tea "github.com/charmbracelet/bubbletea"

type HideListsViewMsg bool

func HideListsViewCmd() tea.Cmd {
	return func() tea.Msg {
		return HideListsViewMsg(true)
	}
}
