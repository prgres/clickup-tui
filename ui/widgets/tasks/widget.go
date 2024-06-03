package tasks

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	tabletasks "github.com/prgrs/clickup/ui/components/table-tasks"
	taskssidebar "github.com/prgrs/clickup/ui/components/tasks-sidebar"
	"github.com/prgrs/clickup/ui/context"
)

const (
	id = "tasks"

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

	log := common.NewLogger(logger, common.ResourceTypeRegistry.WIDGET, id)
	size := common.NewEmptySize()

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

func (m *Model) SetSpinner(f bool) {
	m.showSpinner = f
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeys(msg)

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
		tableTasks[m.componenetTasksTable.SelectedIdx] = m.componenetTasksSidebar.SelectedTask
		m.componenetTasksTable.SetTasks(tableTasks)

	case UpdateTaskMsg:
		m.log.Debug("Received: UpdateTaskMsg")
		t, err := m.ctx.Api.UpdateTask(m.componenetTasksSidebar.SelectedTask)
		if err != nil {
			return common.ErrCmd(err)
		}

		if err := m.componenetTasksSidebar.SetTask(t); err != nil {
			return common.ErrCmd(err)
		}

		tableTasks := m.componenetTasksTable.GetTasks()
		tableTasks[m.componenetTasksTable.SelectedIdx] = m.componenetTasksSidebar.SelectedTask
		m.componenetTasksTable.SetTasks(tableTasks)

		tasks, err := m.ctx.Api.SyncTasksFromView(m.SelectedViewListId)
		if err != nil {
			return common.ErrCmd(err)
		}
		m.componenetTasksTable.SetTasks(tasks)

	case common.RefreshMsg:
		m.log.Debug("Received: common.RefreshMsg")
		t, err := m.ctx.Api.SyncTask(m.componenetTasksSidebar.SelectedTask.Id)
		if err != nil {
			return common.ErrCmd(err)
		}

		if err := m.componenetTasksSidebar.SetTask(t); err != nil {
			return common.ErrCmd(err)
		}

		tasks, err := m.ctx.Api.SyncTasksFromView(m.SelectedViewListId)
		if err != nil {
			return common.ErrCmd(err)
		}
		m.componenetTasksTable.SetTasks(tasks)
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
