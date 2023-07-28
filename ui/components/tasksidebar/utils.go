package tasksidebar

import tea "github.com/charmbracelet/bubbletea"

type InitMsg bool

func InitCmd() tea.Cmd {
	return func() tea.Msg {
		return InitMsg(true)
	}
}
