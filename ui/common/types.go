package common

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
)

type Size struct {
	Width  int
	Height int
}

type ViewId string

type WidgetId string

type ComponentId string

type View interface {
	View() string
	KeyMap() help.KeyMap
	Ready() bool
	SetSize(Size) View
	GetSize() Size
	Update(msg tea.Msg) (View, tea.Cmd)
	GetViewId() ViewId
	Init() tea.Cmd
}

type Widget interface {
	View() string
	KeyMap() help.KeyMap
	GetFocused() bool
	SetFocused(f bool) Widget
	Update(msg tea.Msg) (Widget, tea.Cmd)
	SetSize(s Size) Widget
	Init() tea.Cmd
	SetHidden(h bool) Widget
	GetHidden() bool
}
