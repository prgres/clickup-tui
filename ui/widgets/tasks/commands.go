package tasks

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
)

type LostFocusMsg string

func LostFocusCmd() tea.Cmd {
	return func() tea.Msg {
		return LostFocusMsg("")
	}
}

type UpdateTaskMsg clickup.Task

func UpdateTaskCmd(task clickup.Task) tea.Cmd {
	return func() tea.Msg {
		return UpdateTaskMsg(task)
	}
}
