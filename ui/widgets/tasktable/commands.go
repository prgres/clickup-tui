package tasktable

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
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

type TasksListReloadedMsg []clickup.Task

func TasksListReloadedCmd(tasks []clickup.Task) tea.Cmd {
	return func() tea.Msg {
		return TasksListReloadedMsg(tasks)
	}
}

type ViewChangedMsg string

func ViewChangedCmd(space string) tea.Cmd {
	return func() tea.Msg {
		return ViewChangedMsg(space)
	}
}

type FetchTasksForViewMsg string

func FetchTasksForViewCmd(view string) tea.Cmd {
	return func() tea.Msg {
		return FetchTasksForViewMsg(view)
	}
}
