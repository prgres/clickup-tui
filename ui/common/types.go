package common

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type Size struct {
	Width  int
	Height int
}

func NewEmptySize() Size {
	return Size{
		Width:  0,
		Height: 0,
	}
}

type Id string

type UIElement interface {
	BubblesElem

	Id() Id
	Help() help.KeyMap
	SetSize(Size)
	Size() Size
}

type BubblesElem interface {
	Init() tea.Cmd
	Update(msg tea.Msg) tea.Cmd
	View() string
}
