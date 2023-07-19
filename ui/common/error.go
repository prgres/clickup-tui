package common

import tea "github.com/charmbracelet/bubbletea"

type ErrMsg error

func ErrCmd(err ErrMsg) tea.Cmd {
	return func() tea.Msg {
		return err
	}
}
