package compact

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) handleKeys(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch keypress := msg.String(); keypress {
	case "tab":
		switch m.state {
		case m.widgetNavigator.Id():
			m.state = m.widgetTasks.Id()
			m.widgetTasks.SetFocused(true)
			m.widgetViewsTabs.SetFocused(false)
			m.widgetNavigator.SetFocused(false)
		case m.widgetViewsTabs.Id():
			m.state = m.widgetNavigator.Id()
			m.widgetTasks.SetFocused(false)
			m.widgetViewsTabs.SetFocused(false)
			m.widgetNavigator.SetFocused(true)
		case m.widgetTasks.Id():
			m.state = m.widgetViewsTabs.Id()
			m.widgetTasks.SetFocused(false)
			m.widgetViewsTabs.SetFocused(true)
			m.widgetNavigator.SetFocused(false)
		}
	}

	switch m.state {
	case m.widgetNavigator.Id():
		cmd = m.widgetNavigator.Update(msg)
	case m.widgetViewsTabs.Id():
		cmd = m.widgetViewsTabs.Update(msg)
	case m.widgetTasks.Id():
		cmd = m.widgetTasks.Update(msg)
	}

	m.widgetViewsTabs.Path = m.widgetNavigator.GetPath()

	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}
