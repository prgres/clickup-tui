package tasks

import tea "github.com/charmbracelet/bubbletea"

type LostFocusMsg string

func LostFocusCmd() tea.Cmd {
	return func() tea.Msg {
		return LostFocusMsg("")
	}
}
