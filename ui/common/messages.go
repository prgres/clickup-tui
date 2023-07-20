package common

import tea "github.com/charmbracelet/bubbletea"

type FocusMsg bool

func FocusCmd() tea.Cmd {
	return func() tea.Msg {
		return FocusMsg(true)
	}
}
