package tasktable

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/widgets/viewtabs"
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

type TabChangedMsg viewtabs.Tab

func TabChangedCmd(tab viewtabs.Tab) tea.Cmd {
	return func() tea.Msg {
		return TabChangedMsg(tab)
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
