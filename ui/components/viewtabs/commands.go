package viewtabs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
)

type ViewsListLoadedMsg []clickup.View

type FetchViewsMsg []string

func FetchViewsCmd(spaces []string) tea.Cmd {
	return func() tea.Msg {
		return FetchViewsMsg(spaces)
	}
}

type ViewChangedMsg string

func ViewChangedCmd(view string) tea.Cmd {
	return func() tea.Msg {
		return ViewChangedMsg(view)
	}
}
