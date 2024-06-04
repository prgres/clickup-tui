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
