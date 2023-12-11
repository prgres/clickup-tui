package common

import (
	tea "github.com/charmbracelet/bubbletea"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
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

type ListChangeMsg listitem.Item

func ListChangeCmd(list listitem.Item) tea.Cmd {
	return func() tea.Msg {
		return ListChangeMsg(list)
	}
}

type WorkspaceChangeMsg string

func WorkspaceChangeCmd(workspace string) tea.Cmd {
	return func() tea.Msg {
		return WorkspaceChangeMsg(workspace)
	}
}

type BackToPreviousViewMsg ViewId

func BackToPreviousViewCmd(currentView ViewId) tea.Cmd {
	return func() tea.Msg {
		return BackToPreviousViewMsg(currentView)
	}
}

type ErrMsg error

func ErrCmd(err ErrMsg) tea.Cmd {
	return func() tea.Msg {
		return err
	}
}
