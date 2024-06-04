package workspaceslist

import tea "github.com/charmbracelet/bubbletea"

type (
	WorkspaceChangedMsg  string
	WorkspacePreviewMsg  string
	WorkspaceSelectedMsg string
)

func WorkspaceChangedCmd(id string) tea.Cmd {
	return func() tea.Msg { return WorkspaceChangedMsg(id) }
}

func WorkspacePreviewCmd(workspace string) tea.Cmd {
	return func() tea.Msg { return WorkspacePreviewMsg(workspace) }
}

func WorkspaceSelectedCmd(id string) tea.Cmd {
	return func() tea.Msg { return WorkspaceSelectedMsg(id) }
}
