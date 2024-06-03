package tasks

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/ui/common"
	"golang.design/x/clipboard"
)

type KeyMap struct {
	OpenTicketInWebBrowserBatch key.Binding
	OpenTicketInWebBrowser      key.Binding
	ToggleSidebar               key.Binding
	CopyMode                    key.Binding
	CopyTaskId                  key.Binding
	CopyTaskUrl                 key.Binding
	CopyTaskUrlMd               key.Binding
	LostFocus                   key.Binding
	EditMode                    key.Binding
	EditDescription             key.Binding
	EditName                    key.Binding
	EditStatus                  key.Binding
	EditAssigness               key.Binding
	EditQuit                    key.Binding
	Refresh                     key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		OpenTicketInWebBrowserBatch: key.NewBinding(
			key.WithKeys("U"),
			key.WithHelp("U", "batch open in web browser"),
		),
		OpenTicketInWebBrowser: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "open in web browser"),
		),
		ToggleSidebar: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "toggle sidebar"),
		),
		CopyMode: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "toggle copy mode"),
		),
		CopyTaskId: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "copy task id to clipboard"),
		),
		CopyTaskUrl: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "copy task url to clipboard"),
		),
		CopyTaskUrlMd: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "copy task url as markdown to clipboard"),
		),
		LostFocus: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "lost pane focus"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "go to refresh"),
		),
		EditMode: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit mode"),
		),
		EditDescription: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit description"),
		),
		EditName: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "edit name"),
		),
		EditStatus: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "edit status"),
		),
		EditAssigness: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "edit assigness"),
		),
		EditQuit: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "quit edit mode"),
		),
	}
}

func (m *Model) handleKeys(msg tea.KeyMsg) tea.Cmd {
	var cmds []tea.Cmd

	if m.copyMode {
		return m.handleKeysCopyMode(msg)
	}

	if m.editMode {
		return m.handleKeysEditMode(msg)
	}

	switch {
	case key.Matches(msg, m.keyMap.OpenTicketInWebBrowserBatch):
		tasks := m.componenetTasksTable.GetSelectedTasks()
		for _, task := range tasks {
			m.log.Debug("Opening task in the web browser", "url", task.Url)
			if err := common.OpenUrlInWebBrowser(task.Url); err != nil {
				m.log.Fatal(err)
			}
		}
		return nil

	case key.Matches(msg, m.keyMap.Refresh):
		m.log.Info("Refreshing...")
		if err := m.ctx.Api.Sync(); err != nil {
			m.log.Error("Failed to sync", "error", err)
		}
		m.log.Debug("API sync")
		return nil

	case key.Matches(msg, m.keyMap.OpenTicketInWebBrowser):
		task := m.componenetTasksTable.GetHighlightedTask()
		m.log.Debug("Opening task in the web browser", "url", task.Url)
		if err := common.OpenUrlInWebBrowser(task.Url); err != nil {
			m.log.Fatal(err)
		}
		return nil

	case key.Matches(msg, m.keyMap.ToggleSidebar):
		m.log.Debug("Toggle sidebar")
		m.componenetTasksSidebar.SetHidden(!m.componenetTasksSidebar.GetHidden())
		return nil

	case key.Matches(msg, m.keyMap.CopyMode):
		m.log.Debug("Toggle copy mode")
		m.copyMode = true
		return nil

	case key.Matches(msg, m.keyMap.EditMode):
		m.log.Debug("Toggle edit mode")
		m.editMode = true
		return nil

	case key.Matches(msg, m.keyMap.LostFocus):
		switch m.state {
		case m.componenetTasksSidebar.Id():
			m.state = m.componenetTasksTable.Id()
			m.componenetTasksSidebar.SetFocused(false)
			m.componenetTasksTable.SetFocused(true)

		case m.componenetTasksTable.Id():
			m.componenetTasksSidebar.SetFocused(false)
			m.componenetTasksTable.SetFocused(false)

			cmds = append(cmds, LostFocusCmd())
		}

		return tea.Batch(append(cmds,
			m.componenetTasksSidebar.Update(msg),
			m.componenetTasksTable.Update(msg),
		)...)

	}

	var cmd tea.Cmd
	switch m.state {
	case m.componenetTasksSidebar.Id():
		cmd = m.componenetTasksSidebar.Update(msg)
	case m.componenetTasksTable.Id():
		cmd = m.componenetTasksTable.Update(msg)
	}

	return tea.Batch(append(cmds, cmd)...)
}

func (m *Model) handleKeysCopyMode(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keyMap.CopyTaskId):
		task := m.componenetTasksTable.GetHighlightedTask()
		clipboard.Write(clipboard.FmtText, []byte(task.Id))
		m.copyMode = false

	case key.Matches(msg, m.keyMap.CopyTaskUrl):
		task := m.componenetTasksTable.GetHighlightedTask()
		clipboard.Write(clipboard.FmtText, []byte(task.Url))
		m.copyMode = false

	case key.Matches(msg, m.keyMap.CopyTaskUrlMd):
		task := m.componenetTasksTable.GetHighlightedTask()
		md := fmt.Sprintf("[[#%s] - %s](%s)", task.Id, task.Name, task.Url)
		clipboard.Write(clipboard.FmtText, []byte(md))
		m.copyMode = false

	case key.Matches(msg, m.keyMap.LostFocus):
		m.copyMode = false
	}

	return nil
}

func (m *Model) handleKeysEditMode(msg tea.KeyMsg) tea.Cmd {
	switch {
	case key.Matches(msg, m.keyMap.EditDescription):
		data := m.componenetTasksSidebar.SelectedTask.MarkdownDescription
		m.editMode = false
		return common.OpenEditor(editorIdDescription, data)

	case key.Matches(msg, m.keyMap.EditName):
		data := m.componenetTasksSidebar.SelectedTask.Name
		m.editMode = false
		return common.OpenEditor(editorIdName, data)

	case key.Matches(msg, m.keyMap.EditStatus):
		data := m.componenetTasksSidebar.SelectedTask.Status.Status
		m.editMode = false
		return common.OpenEditor(editorIdStatus, data)

	// case key.Matches(msg, m.keyMap.EditAssigness):
	// 	data := m.SelectedTask.Assignees
	// 	cmds = append(cmds, common.OpenEditor(data))
	// 	m.editMode = false

	case key.Matches(msg, m.keyMap.EditQuit):
		m.editMode = false
	}

	return nil
}
