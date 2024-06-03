package navigator

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) handleKeys(msg tea.KeyMsg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch keypress := msg.String(); keypress {
	case "esc":
		m.log.Info("Received: Go to previous view")

		switch m.state {
		case m.componentSpacesList.Id():
			m.state = m.componentWorkspacesList.Id()
		case m.componentFoldersList.Id():
			m.state = m.componentSpacesList.Id()
		case m.componentListsList.Id():
			m.state = m.componentFoldersList.Id()
		}

		cmds = append(cmds, cmd)
		return tea.Batch(cmds...)
	}

	switch m.state {
	case m.componentWorkspacesList.Id():
		cmd = m.componentWorkspacesList.Update(msg)
	case m.componentSpacesList.Id():
		cmd = m.componentSpacesList.Update(msg)
	case m.componentFoldersList.Id():
		cmd = m.componentFoldersList.Update(msg)
	case m.componentListsList.Id():
		cmd = m.componentListsList.Update(msg)
	}

	cmds = append(cmds, cmd)
	return tea.Batch(cmds...)
}
