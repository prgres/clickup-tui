package taskssidebar

import tea "github.com/charmbracelet/bubbletea"

type InitMsg bool

func InitCmd() tea.Cmd {
	return func() tea.Msg {
		return InitMsg(true)
	}
}

type TaskSelectedMsg string

func TaskSelectedCmd(task string) tea.Cmd {
	return func() tea.Msg {
		return TaskSelectedMsg(task)
	}
}
