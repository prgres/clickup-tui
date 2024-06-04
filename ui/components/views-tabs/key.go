package viewstabs

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	CursorLeft           key.Binding
	CursorLeftAndSelect  key.Binding
	CursorRight          key.Binding
	CursorRightAndSelect key.Binding
	Select               key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		CursorLeft: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("h, left", "previous tab"),
		),
		CursorLeftAndSelect: key.NewBinding(
			key.WithKeys("H", "left"),
			key.WithHelp("H, shift+left", "select tab"),
		),
		CursorRight: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("l, right", "next tab"),
		),
		CursorRightAndSelect: key.NewBinding(
			key.WithKeys("L", "shift+right"),
			key.WithHelp("L, shift+right", "select tab"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

func (m *Model) handleKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keyMap.CursorLeft):
		index := prevTab(m.tabs, m.SelectedIdx)
		if m.SelectedIdx == index {
			break
		}
		m.SelectedIdx = index
		m.Selected = m.tabs[index].Id
		return nil

	case key.Matches(msg, m.keyMap.CursorRight):
		index := nextTab(m.tabs, m.SelectedIdx)
		if m.SelectedIdx == index {
			break
		}
		m.SelectedIdx = index
		m.Selected = m.tabs[index].Id
		return nil

	case key.Matches(msg, m.keyMap.Select):
		m.Selected = m.tabs[m.SelectedIdx].Id
		return TabChangedCmd(m.Selected)

	case key.Matches(msg, m.keyMap.CursorLeftAndSelect):
		index := prevTab(m.tabs, m.SelectedIdx)
		if m.SelectedIdx == index {
			break
		}
		m.SelectedIdx = index
		m.Selected = m.tabs[index].Id
		return TabChangedCmd(m.Selected)

	case key.Matches(msg, m.keyMap.CursorRightAndSelect):
		index := nextTab(m.tabs, m.SelectedIdx)
		if m.SelectedIdx == index {
			break
		}
		m.SelectedIdx = index
		m.Selected = m.tabs[index].Id
		return TabChangedCmd(m.Selected)
	}

	return nil
}
