package folders

import tea "github.com/charmbracelet/bubbletea"

type HideFolderViewMsg bool

func HideFolderViewCmd() tea.Cmd {
	return func() tea.Msg {
		return HideFolderViewMsg(true)
	}
}
