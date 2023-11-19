package taskstabs

import (
	tea "github.com/charmbracelet/bubbletea"
)

type FetchTasksForTabsMsg []Tab

func FetchTasksForTabsCmd(tabs []Tab) tea.Cmd {
	return func() tea.Msg {
		return FetchTasksForTabsMsg(tabs)
	}
}

type TabChangedMsg Tab

func TabChangedCmd(tab Tab) tea.Cmd {
	return func() tea.Msg {
		return TabChangedMsg(tab)
	}
}

type TabLoadedMsg Tab

func TabLoadedCmd(tab Tab) tea.Cmd {
	return func() tea.Msg {
		return TabLoadedMsg(tab)
	}
}
