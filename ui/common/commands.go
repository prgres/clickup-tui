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

func SpaceChangeCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return SpaceChangeMsg(id)
	}
}

type SpacePreviewMsg string

func SpacePreviewCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return SpacePreviewMsg(id)
	}
}

type FolderChangeMsg string

func FolderChangeCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return FolderChangeMsg(id)
	}
}

type FolderPreviewMsg string

func FolderPreviewCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return FolderPreviewMsg(id)
	}
}

type ListChangeMsg string

func ListChangeCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return ListChangeMsg(id)
	}
}

type ListPreviewMsg string

func ListPreviewCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return ListPreviewMsg(id)
	}
}

type WorkspaceChangeMsg string

func WorkspaceChangeCmd(workspace string) tea.Cmd {
	return func() tea.Msg {
		return WorkspaceChangeMsg(workspace)
	}
}

type WorkspacePreviewMsg string

func WorkspacePreviewCmd(workspace string) tea.Cmd {
	return func() tea.Msg {
		return WorkspacePreviewMsg(workspace)
	}
}

type BackToPreviousViewMsg Id

func BackToPreviousViewCmd(currentView Id) tea.Cmd {
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

type UITickMsg int64

func UITickCmd(ts int64) tea.Cmd {
	return func() tea.Msg {
		return UITickMsg(ts)
	}
}

type RefreshMsg string

func RefreshCmd() tea.Cmd {
	return func() tea.Msg {
		return RefreshMsg("")
	}
}
