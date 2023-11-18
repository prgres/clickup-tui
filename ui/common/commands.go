package common

import (
	tea "github.com/charmbracelet/bubbletea"
)

type FocusMsg bool

func FocusCmd() tea.Cmd {
	return func() tea.Msg {
		return FocusMsg(true)
	}
}

type SpaceChangeMsg string

func SpaceChangeCmd(space string) tea.Cmd {
	return func() tea.Msg {
		return SpaceChangeMsg(space)
	}
}

type FolderChangeMsg string

func FolderChangeCmd(space string) tea.Cmd {
	return func() tea.Msg {
		return FolderChangeMsg(space)
	}
}

type ListChangeMsg string

func ListChangeCmd(space string) tea.Cmd {
	return func() tea.Msg {
		return ListChangeMsg(space)
	}
}

type TeamChangeMsg string

func TeamChangeCmd(team string) tea.Cmd {
	return func() tea.Msg {
		return TeamChangeMsg(team)
	}
}

type BackToPreviousViewMsg ViewId

func BackToPreviousViewCmd(currentView ViewId) tea.Cmd {
	return func() tea.Msg {
		return BackToPreviousViewMsg(currentView)
	}
}
