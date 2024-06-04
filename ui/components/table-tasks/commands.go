package tabletasks

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	TaskSelectedMsg   string
	TasksListReadyMsg bool
	HideTableMsg      bool
)

func TaskSelectedCmd(task string) tea.Cmd {
	return func() tea.Msg {
		return TaskSelectedMsg(task)
	}
}

func TasksListReadyCmd() tea.Cmd {
	return func() tea.Msg {
		return TasksListReadyMsg(true)
	}
}

func HideTableCmd() tea.Cmd {
	return func() tea.Msg {
		return HideTableMsg(true)
	}
}
