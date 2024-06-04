package compact

import tea "github.com/charmbracelet/bubbletea"

type (
	InitCompactMsg          string
	LoadingTasksFromViewMsg string
)

func InitCompactCmd() tea.Cmd {
	return func() tea.Msg {
		return InitCompactMsg("")
	}
}

func LoadingTasksFromViewCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return LoadingTasksFromViewMsg(id)
	}
}
