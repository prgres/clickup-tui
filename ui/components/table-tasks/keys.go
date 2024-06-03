package tabletasks

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
	"github.com/prgrs/clickup/ui/common"
)

type KeyMap struct {
	table.KeyMap
	Select key.Binding
}

func DefaultKeyMap() KeyMap {
	km := table.DefaultKeyMap()

	return KeyMap{
		KeyMap: table.KeyMap{
			RowDown:         common.KeyBindingWithHelp(km.RowDown, "down"),
			RowUp:           common.KeyBindingWithHelp(km.RowUp, "up"),
			RowSelectToggle: common.KeyBindingWithHelp(km.RowSelectToggle, "select"),
			PageDown:        common.KeyBindingWithHelp(km.PageDown, "next page"),
			PageUp:          common.KeyBindingWithHelp(km.PageUp, "previous page"),
			PageFirst:       common.KeyBindingWithHelp(km.PageFirst, "first page"),
			PageLast:        common.KeyBindingWithHelp(km.PageLast, "last page"),
			Filter:          common.KeyBindingWithHelp(km.Filter, "filter"),
			FilterBlur:      common.KeyBindingWithHelp(km.FilterBlur, "filter blur"),
			FilterClear:     common.KeyBindingWithHelp(km.FilterClear, "filter clear"),
			ScrollRight:     common.KeyBindingWithHelp(km.ScrollRight, "scroll right"),
			ScrollLeft:      common.KeyBindingWithHelp(km.ScrollLeft, "scroll left"),
		},
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

func (m *Model) handleKeys(msg tea.KeyMsg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch {
	case key.Matches(msg, m.keyMap.Select):
		index := m.table.GetHighlightedRowIndex()
		if m.table.TotalRows() == 0 {
			m.log.Info("Table is empty")
			break
		}
		taskId := m.tasks[index].Id
		m.log.Infof("Receive enter: %d", index)
		cmds = append(cmds, TaskSelectedCmd(taskId))
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}
