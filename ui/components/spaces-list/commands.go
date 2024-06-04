package spaceslist

import tea "github.com/charmbracelet/bubbletea"

type (
	SpaceChangedMsg  string
	SpacePreviewMsg  string
	SpaceSelectedMsg string
)

func SpaceChangedCmd(id string) tea.Cmd {
	return func() tea.Msg { return SpaceChangedMsg(id) }
}

func SpacePreviewCmd(id string) tea.Cmd {
	return func() tea.Msg { return SpacePreviewMsg(id) }
}

func SpaceSelectedCmd(id string) tea.Cmd {
	return func() tea.Msg { return SpaceSelectedMsg(id) }
}
