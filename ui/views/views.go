package views

import tea "github.com/charmbracelet/bubbletea"

type WindowSizeMsg tea.WindowSizeMsg

func WindowSizeCmd(msg tea.WindowSizeMsg) tea.Cmd {
	return func() tea.Msg {
		return WindowSizeMsg(msg)
	}
}
