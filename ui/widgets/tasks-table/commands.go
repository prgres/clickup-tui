package taskstable

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/widgets/tasks-tabs"
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

type TabChangedMsg taskstabs.Tab

func TabChangedCmd(tab taskstabs.Tab) tea.Cmd {
	return func() tea.Msg {
		return TabChangedMsg(tab)
	}
}

type HideTableMsg bool

func HideTableCmd() tea.Cmd {
	return func() tea.Msg {
		return HideTableMsg(true)
	}
}

type FetchTasksForViewMsg string

func FetchTasksForViewCmd(view string) tea.Cmd {
	return func() tea.Msg {
		return FetchTasksForViewMsg(view)
	}
}

type FetchTasksForListMsg string

func FetchTasksForListCmd(list string) tea.Cmd {
	return func() tea.Msg {
		return FetchTasksForListMsg(list)
	}
}
