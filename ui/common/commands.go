package common

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type FocusMsg bool

func FocusCmd() tea.Cmd {
	return func() tea.Msg {
		return FocusMsg(true)
	}
}

type SpaceChangedMsg string

func SpaceChangedCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return SpaceChangedMsg(id)
	}
}

type SpacePreviewMsg string

func SpacePreviewCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return SpacePreviewMsg(id)
	}
}

type FolderChangedMsg string

func FolderChangedCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return FolderChangedMsg(id)
	}
}

type FolderPreviewMsg string

func FolderPreviewCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return FolderPreviewMsg(id)
	}
}

type ListChangedMsg string

func ListChangedCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return ListChangedMsg(id)
	}
}

type ListPreviewMsg string

func ListPreviewCmd(id string) tea.Cmd {
	return func() tea.Msg {
		return ListPreviewMsg(id)
	}
}

type WorkspaceChangedMsg string

func WorkspaceChangedCmd(workspace string) tea.Cmd {
	return func() tea.Msg {
		return WorkspaceChangedMsg(workspace)
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

func (m UITickMsg) Tick() tea.Cmd {
	return func() tea.Msg {
		return m
	}
}

func UITickCmd(ts int64) tea.Cmd {
	return func() tea.Msg {
		return UITickMsg(time.Now().Unix() + ts)
	}
}

type RefreshMsg string

func RefreshCmd() tea.Cmd {
	return func() tea.Msg {
		return RefreshMsg("")
	}
}
