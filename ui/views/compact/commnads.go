package compact

import tea "github.com/charmbracelet/bubbletea"

type InitCompactMsg string

func InitCompactCmd() tea.Cmd {
	return func() tea.Msg {
		return InitCompactMsg("")
	}
}
