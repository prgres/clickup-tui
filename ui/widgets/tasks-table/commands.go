package taskstable

import (
	tea "github.com/charmbracelet/bubbletea"
)

type TaskSelectedMsg string

func TaskSelectedCmd(task string) tea.Cmd {
	return func() tea.Msg {
		return TaskSelectedMsg(task)
	}
}

type TasksListReadyMsg bool

func TasksListReadyCmd() tea.Cmd {
	return func() tea.Msg {
		return TasksListReadyMsg(true)
	}
}

type HideTableMsg bool

func HideTableCmd() tea.Cmd {
	return func() tea.Msg {
		return HideTableMsg(true)
	}
}
