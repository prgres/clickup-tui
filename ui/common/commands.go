package common

import tea "github.com/charmbracelet/bubbletea"

type FocusMsg bool

func FocusCmd() tea.Cmd {
	return func() tea.Msg {
		return FocusMsg(true)
	}
}

type WindowSizeMsg tea.WindowSizeMsg

func WindowSizeCmd(msg tea.WindowSizeMsg) tea.Cmd {
	return func() tea.Msg {
		return WindowSizeMsg(msg)
	}
}

type SpaceChangeMsg string

func SpaceChangeCmd(space string) tea.Cmd {
	return func() tea.Msg {
		return SpaceChangeMsg(space)
	}
}
