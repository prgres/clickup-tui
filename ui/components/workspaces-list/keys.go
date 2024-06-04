package workspaceslist

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
)

type KeyMap struct {
	CursorUp            key.Binding
	CursorUpAndSelect   key.Binding
	CursorDown          key.Binding
	CursorDownAndSelect key.Binding
	Select              key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		CursorUp: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k, up", "up"),
		),
		CursorUpAndSelect: key.NewBinding(
			key.WithKeys("K", "shift+up"),
			key.WithHelp("K, shift+up", "up and select"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j, down", "down"),
		),
		CursorDownAndSelect: key.NewBinding(
			key.WithKeys("J", "shift+down"),
			key.WithHelp("J, down", "down and select"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

func (m *Model) handleKeys(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keyMap.Select):
		if m.list.SelectedItem() == nil {
			m.log.Info("List is empty")
			break
		}
		selected := m.list.SelectedItem().(listitem.Item).Data().(clickup.Workspace)
		m.log.Info("Selected workspace", "id", selected.Id, "name", selected.Name)
		m.Selected = selected
		return WorkspaceChangedCmd(selected.Id)

	case key.Matches(msg, m.keyMap.CursorDown):
		m.list.CursorDown()

	case key.Matches(msg, m.keyMap.CursorDownAndSelect):
		m.list.CursorDown()
		if m.list.SelectedItem() == nil {
			m.log.Info("List is empty")
			break
		}
		selectedWorkspace := m.list.SelectedItem().(listitem.Item).Data().(clickup.Workspace)
		m.log.Info("Selected workspace", "id", selectedWorkspace.Id, "name", selectedWorkspace.Name)
		m.Selected = selectedWorkspace
		return WorkspacePreviewCmd(selectedWorkspace.Id)

	case key.Matches(msg, m.keyMap.CursorUp):
		m.list.CursorUp()

	case key.Matches(msg, m.keyMap.CursorUpAndSelect):
		m.list.CursorUp()
		if m.list.SelectedItem() == nil {
			m.log.Info("List is empty")
			break
		}
		selectedWorkspace := m.list.SelectedItem().(listitem.Item).Data().(clickup.Workspace)
		m.log.Info("Selected workspace", "id", selectedWorkspace.Id, "name", selectedWorkspace.Name)
		m.Selected = selectedWorkspace
		return WorkspacePreviewCmd(selectedWorkspace.Id)
	}

	return nil
}
