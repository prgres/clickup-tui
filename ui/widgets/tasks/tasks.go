package tasks

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	tabletasks "github.com/prgrs/clickup/ui/components/table-tasks"
	taskssidebar "github.com/prgrs/clickup/ui/components/tasks-sidebar"
	"github.com/prgrs/clickup/ui/context"
	"golang.design/x/clipboard"
)

const (
	id                  = "tasks"
	editorIdDescription = "description"
	editorIdName        = "name"
	editorIdStatus      = "status"
)

type Model struct {
	log                *log.Logger
	ctx                *context.UserContext
	id                 common.Id
	size               common.Size
	Focused            bool
	Hidden             bool
	ifBorders          bool
	keyMap             KeyMap
	state              common.Id
	spinner            spinner.Model
	showSpinner        bool
	SelectedViewListId string

	copyMode bool // TODO make as a widget
	editMode bool

	componenetTasksTable   *tabletasks.Model
	componenetTasksSidebar *taskssidebar.Model
}

func (m Model) Id() common.Id {
	return m.id
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	log := logger.WithPrefix(logger.GetPrefix() + "/widget/" + id)

	size := common.Size{
		Width:  0,
		Height: 0,
	}
	var (
		componenetTasksTable   = tabletasks.InitialModel(ctx, log)
		componenetTasksSidebar = taskssidebar.InitialModel(ctx, log).WithHidden(true)
	)

	return Model{
		id:                     id,
		ctx:                    ctx,
		size:                   size,
		Focused:                false,
		Hidden:                 false,
		keyMap:                 DefaultKeyMap(),
		log:                    log,
		ifBorders:              true,
		state:                  componenetTasksTable.Id(),
		spinner:                s,
		showSpinner:            false,
		copyMode:               false,
		editMode:               false,
		componenetTasksTable:   &componenetTasksTable,
		componenetTasksSidebar: &componenetTasksSidebar,
	}
}

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

func (m Model) Help() help.KeyMap {
	var help help.KeyMap

	if m.copyMode {
		return common.NewHelp(
			func() [][]key.Binding {
				return [][]key.Binding{
					{
						m.keyMap.CopyTaskId,
						m.keyMap.CopyTaskUrl,
						m.keyMap.CopyTaskUrlMd,
						m.keyMap.LostFocus,
						m.keyMap.Refresh,
					},
				}
			},
			func() []key.Binding {
				return []key.Binding{
					m.keyMap.CopyTaskId,
					m.keyMap.CopyTaskUrl,
					m.keyMap.CopyTaskUrlMd,
					m.keyMap.LostFocus,
					m.keyMap.Refresh,
				}
			},
		)
	}

	if m.editMode {
		return common.NewHelp(
			func() [][]key.Binding {
				return [][]key.Binding{
					{
						m.keyMap.EditDescription,
						m.keyMap.EditName,
						m.keyMap.EditStatus,
						m.keyMap.EditAssigness,
						m.keyMap.EditQuit,
					},
				}
			},
			func() []key.Binding {
				return []key.Binding{
					m.keyMap.EditDescription,
					m.keyMap.EditName,
					m.keyMap.EditStatus,
					m.keyMap.EditAssigness,
					m.keyMap.EditQuit,
				}
			},
		)
	}
	switch m.state {
	case m.componenetTasksSidebar.Id():
		help = m.componenetTasksSidebar.Help()
	case m.componenetTasksTable.Id():
		help = m.componenetTasksTable.Help()
	}

	return common.NewHelp(
		func() [][]key.Binding {
			return append(
				help.FullHelp(),
				[]key.Binding{
					m.keyMap.OpenTicketInWebBrowser,
					m.keyMap.ToggleSidebar,
					m.keyMap.EditMode,
				},
			)
		},
		func() []key.Binding {
			return append(
				help.ShortHelp(),
				m.keyMap.OpenTicketInWebBrowser,
				m.keyMap.CopyMode,
				m.keyMap.EditMode,
			)
		},
	)
}

func (m *Model) SetSpinner(f bool) {
	m.showSpinner = f
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.copyMode {
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

			return tea.Batch(cmds...)
		}

		if m.editMode {
			switch {
			case key.Matches(msg, m.keyMap.EditDescription):
				data := m.componenetTasksSidebar.SelectedTask.MarkdownDescription
				cmds = append(cmds, common.OpenEditor(editorIdDescription, data))
				m.editMode = false

			case key.Matches(msg, m.keyMap.EditName):
				data := m.componenetTasksSidebar.SelectedTask.Name
				cmds = append(cmds, common.OpenEditor(editorIdName, data))
				m.editMode = false

			case key.Matches(msg, m.keyMap.EditStatus):
				data := m.componenetTasksSidebar.SelectedTask.Status.Status
				cmds = append(cmds, common.OpenEditor(editorIdStatus, data))
				m.editMode = false

			// case key.Matches(msg, m.keyMap.EditAssigness):
			// 	data := m.SelectedTask.Assignees
			// 	cmds = append(cmds, common.OpenEditor(data))
			// 	m.editMode = false

			case key.Matches(msg, m.keyMap.EditQuit):
				m.editMode = false
			}

			return tea.Batch(cmds...)
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

		case key.Matches(msg, m.keyMap.Refresh):
			m.log.Info("Refreshing...")
			if err := m.ctx.Api.Sync(); err != nil {
				m.log.Error("Failed to sync", "error", err)
			}
			m.log.Debug("API sync")

		case key.Matches(msg, m.keyMap.OpenTicketInWebBrowser):
			task := m.componenetTasksTable.GetHighlightedTask()
			m.log.Debug("Opening task in the web browser", "url", task.Url)
			if err := common.OpenUrlInWebBrowser(task.Url); err != nil {
				m.log.Fatal(err)
			}

		case key.Matches(msg, m.keyMap.ToggleSidebar):
			m.log.Debug("Toggle sidebar")
			m.componenetTasksSidebar.SetHidden(!m.componenetTasksSidebar.GetHidden())

		case key.Matches(msg, m.keyMap.CopyMode):
			m.log.Debug("Toggle copy mode")
			m.copyMode = true

		case key.Matches(msg, m.keyMap.EditMode):
			m.log.Debug("Toggle edit mode")
			m.editMode = true

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

			cmds = append(cmds,
				m.componenetTasksSidebar.Update(msg),
				m.componenetTasksTable.Update(msg),
			)

			return tea.Batch(cmds...)
		}

		switch m.state {
		case m.componenetTasksSidebar.Id():
			cmd = m.componenetTasksSidebar.Update(msg)
		case m.componenetTasksTable.Id():
			cmd = m.componenetTasksTable.Update(msg)
		}

		cmds = append(cmds, cmd)

		return tea.Batch(cmds...)

	case tabletasks.TaskSelectedMsg:
		id := string(msg)
		m.log.Infof("Received: taskstable.TaskSelectedMsg: %s", id)

		m.state = m.componenetTasksSidebar.Id()

		m.componenetTasksSidebar.
			SetFocused(true).
			SetHidden(false)

		m.componenetTasksTable.SetFocused(false)

		if err := m.componenetTasksSidebar.SelectTask(id); err != nil {
			cmds = append(cmds, common.ErrCmd(err))
		}

		cmds = append(cmds, cmd)

	case common.EditorFinishedMsg:
		err := msg.Err
		id := msg.Id
		if err := err; err != nil {
			return common.ErrCmd(err)
		}

		switch id {
		case editorIdDescription:
			data := msg.Data.(string)
			m.componenetTasksSidebar.SelectedTask.Description = data
		case editorIdName:
			data := msg.Data.(string)
			m.componenetTasksSidebar.SelectedTask.Name = data
		case editorIdStatus:
			data := msg.Data.(string)
			m.componenetTasksSidebar.SelectedTask.Status.Status = data
		}

		cmds = append(cmds, UpdateTaskCmd(m.componenetTasksSidebar.SelectedTask))

		if err := m.componenetTasksSidebar.SetTask(m.componenetTasksSidebar.SelectedTask); err != nil {
			return common.ErrCmd(err)
		}

		tableTasks := m.componenetTasksTable.GetTasks()
		tableTasks[m.componenetTasksTable.SelectedTaskIndex] = m.componenetTasksSidebar.SelectedTask
		m.componenetTasksTable.SetTasks(tableTasks)

	case UpdateTaskMsg:
		t, err := m.ctx.Api.UpdateTask(m.componenetTasksSidebar.SelectedTask)
		if err != nil {
			return common.ErrCmd(err)
		}

		if err := m.componenetTasksSidebar.SetTask(t); err != nil {
			return common.ErrCmd(err)
		}

		tableTasks := m.componenetTasksTable.GetTasks()
		tableTasks[m.componenetTasksTable.SelectedTaskIndex] = m.componenetTasksSidebar.SelectedTask
		m.componenetTasksTable.SetTasks(tableTasks)

		// TODO: this is temp solution withouth err checking
		// because we are not able to distinguish upstream
		m.ctx.Api.SyncTasksFromView(m.SelectedViewListId) //nolint:errcheck
		m.ctx.Api.SyncTasksFromList(m.SelectedViewListId) //nolint:errcheck
	}

	cmds = append(cmds,
		m.componenetTasksTable.Update(msg),
		m.componenetTasksSidebar.Update(msg),
	)

	return tea.Batch(cmds...)
}

func (m *Model) SetTasks(tasks []clickup.Task) {
	m.showSpinner = false
	m.componenetTasksTable.SetTasks(tasks)

	if len(tasks) == 0 {
		m.componenetTasksSidebar.SetHidden(true)
		return
	}

	// TODO: check if it should yield at all or move it to cmd
	id := tasks[0].Id
	if err := m.componenetTasksSidebar.SelectTask(id); err != nil {
		m.log.Fatal(err)
	}
}

func (m Model) View() string {
	bColor := m.ctx.Theme.BordersColorInactive
	if m.Focused {
		bColor = m.ctx.Theme.BordersColorActive
	}

	if m.copyMode {
		bColor = m.ctx.Theme.BordersColorCopyMode
	}

	if m.editMode {
		bColor = m.ctx.Theme.BordersColorEditMode
	}

	style := lipgloss.NewStyle().
		Width(m.size.Width).
		MaxWidth(m.size.Width).
		Height(m.size.Height).
		MaxHeight(m.size.Height)

	styleBorders := m.ctx.Style.Borders.
		BorderForeground(bColor)

	borderMargin := 0
	if m.ifBorders {
		borderMargin = 2
	}

	size := common.Size{
		Width:  m.size.Width - borderMargin,
		Height: m.size.Height - borderMargin,
	}

	if m.showSpinner {
		return style.
			Inherit(styleBorders).
			Width(m.size.Width - borderMargin).
			MaxWidth(m.size.Width + borderMargin).
			Height(m.size.Height - borderMargin).
			MaxHeight(m.size.Height + borderMargin).
			Render(
				lipgloss.Place(
					size.Width, size.Height,
					lipgloss.Center,
					lipgloss.Center,
					fmt.Sprintf("%s Loading lists...", m.spinner.View()),
				),
			)
	}

	if m.componenetTasksTable.TotalRows() == 0 {
		return style.
			Inherit(styleBorders).
			Width(m.size.Width - borderMargin).
			MaxWidth(m.size.Width + borderMargin).
			Height(m.size.Height - borderMargin).
			MaxHeight(m.size.Height + borderMargin).
			Render(
				lipgloss.Place(
					size.Width, size.Height,
					lipgloss.Center,
					lipgloss.Center,
					"No tasks found",
				),
			)
	}

	tasksTableBorders := m.ctx.Theme.BordersColorInactive
	if m.componenetTasksTable.GetFocused() {
		tasksTableBorders = m.ctx.Theme.BordersColorActive
	}

	if m.copyMode {
		tasksTableBorders = m.ctx.Theme.BordersColorCopyMode
	}

	if m.editMode {
		tasksTableBorders = m.ctx.Theme.BordersColorEditMode
	}

	tmpStyle := style.
		Inherit(styleBorders)

	var (
		contentTasksTable   string
		contentTasksSidebar string
	)

	if m.componenetTasksSidebar.Hidden {
		m.componenetTasksTable.SetSize(size)

		tmpStyle = tmpStyle.
			BorderForeground(tasksTableBorders).
			Width(size.Width).
			MaxWidth(size.Width + borderMargin).
			Height(size.Height).
			MaxHeight(m.size.Height + borderMargin)

		m.componenetTasksTable.SetSize(size)
		contentTasksTable = tmpStyle.Render(m.componenetTasksTable.View())
		contentTasksSidebar = ""
	} else {
		// TODO: WTF?!
		size.Width /= 2

		size.Height += borderMargin
		m.componenetTasksSidebar.SetSize(size)
		size.Height -= borderMargin

		size.Width -= borderMargin // size.Width -= 2 * borderMargin

		tmpStyle = tmpStyle.
			Width(size.Width).
			BorderForeground(tasksTableBorders).
			MaxWidth(size.Width + borderMargin).
			Height(size.Height).
			MaxHeight(m.size.Height + borderMargin)

		m.componenetTasksTable.SetSize(size)

		contentTasksTable = tmpStyle.Render(m.componenetTasksTable.View())
		contentTasksSidebar = m.componenetTasksSidebar.View()
	}

	return style.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			contentTasksTable,
			contentTasksSidebar,
		))
}

func (m *Model) SetFocused(f bool) {
	m.Focused = f

	switch m.state {
	case m.componenetTasksSidebar.Id():
		m.componenetTasksSidebar.SetFocused(f)
	case m.componenetTasksTable.Id():
		m.componenetTasksTable.SetFocused(f)
	}
}

func (m *Model) SetSize(s common.Size) {
	m.size = s
}

func (m *Model) Init() error {
	m.log.Info("Initializing...")
	return nil
}

func (m Model) Size() common.Size {
	return m.size
}
