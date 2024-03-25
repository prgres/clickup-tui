package workspaces

import tea "github.com/charmbracelet/bubbletea"

type InitWorkspacesMsg string

func InitWorkspacesCmd() tea.Cmd {
	return func() tea.Msg {
		return InitWorkspacesMsg("")
	}
}
