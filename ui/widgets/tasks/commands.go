package tasks

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
)

type (
	LostFocusMsg  string
	UpdateTaskMsg clickup.Task
)

func LostFocusCmd() tea.Cmd {
	return func() tea.Msg { return LostFocusMsg("") }
}

func UpdateTaskCmd(task clickup.Task) tea.Cmd {
	return func() tea.Msg { return UpdateTaskMsg(task) }
}
