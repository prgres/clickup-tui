package tickets

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
)

type TasksListReloadedMsg []clickup.Task

func TasksListReloadedCmd(tasks []clickup.Task) tea.Cmd {
	return func() tea.Msg {
		return TasksListReloadedMsg(tasks)
	}
}

type SpaceChangedMsg string

func SpaceChangedCmd(space string) tea.Cmd {
	return func() tea.Msg {
		return SpaceChangedMsg(space)
	}
}

type Model struct {
	ctx           *context.UserContext
	table         table.Model
	columns       []table.Column
	tickets       map[string][]clickup.Task
	SelectedSpace string
	spinner       spinner.Model
	showSpinner   bool

	width  int
	height int
}

func InitialModel(ctx *context.UserContext) Model {
	columns := []table.Column{
		{Title: "Status", Width: 15},
		{Title: "Name", Width: 90},
		{Title: "Assignees", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	s := spinner.New()
	s.Spinner = spinner.Pulse
	// spinner.Dot,
	// spinner.Line,
	// spinner.Pulse,
	// spinner.Points,
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		ctx:           ctx,
		table:         t,
		columns:       columns,
		tickets:       map[string][]clickup.Task{},
		SelectedSpace: SPACE_SRE,
		spinner:       s,
		showSpinner:   true,
	}
}

func (m Model) syncTable(tasks []clickup.Task) Model {
	m.ctx.Logger.Info("Synchonizing table")
	m.tickets[m.SelectedSpace] = tasks

	items := taskListToRows(tasks)
	m.table.SetRows(items)

	return m
}

func taskToRow(task clickup.Task) table.Row {
	return table.Row{
		task.Status.Status,
		task.Name,
		task.GetAssignees(),
	}
}

func taskListToRows(tasks []clickup.Task) []table.Row {
	rows := make([]table.Row, len(tasks))
	for i, task := range tasks {
		rows[i] = taskToRow(task)
	}
	return rows
}
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case SpaceChangedMsg:
		m.ctx.Logger.Info("TaskView receive SpaceChangedMsg")
		m.showSpinner = true
		m.SelectedSpace = string(msg)
		return m, tea.Batch(m.getTicketsCmd(string(msg)), spinner.Tick)

	case TasksListReloadedMsg:
		m.ctx.Logger.Info("TaskView receive TasksListReloadedMsg")
		m = m.syncTable(msg)
		m.table, cmd = m.table.Update(msg)
		m.showSpinner = false
		return m, cmd

	case tea.WindowSizeMsg:
		m.ctx.Logger.Info("TaskView receive tea.WindowSizeMsg")
		h, v := docStyle.GetFrameSize()

		m.table.SetWidth(msg.Width - h)
		m.table.SetHeight(msg.Height - v)

		m.width = msg.Width
		m.height = msg.Height

		return m, nil

	case spinner.TickMsg:
		m.ctx.Logger.Info("TaskView receive spinner.TickMsg")
		if !m.showSpinner {
			return m, nil
		}
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.showSpinner {
		return lipgloss.Place(
			m.width, m.height,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("%s Loading tasks...", m.spinner.View()),
		)
	}

	return m.table.View()
}

func (m Model) Init() tea.Msg {
	return SpaceChangedMsg(m.SelectedSpace)
}

func (m Model) getTicketsCmd(space string) tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.getTickets(space)
		if err != nil {
			return common.ErrMsg(err)
		}

		return TasksListReloadedMsg(tasks)
	}
}

func (m Model) getTickets(space string) ([]clickup.Task, error) {
	m.ctx.Logger.Infof("Getting tasks for space: %s", space)
	if m.tickets[space] != nil {
		m.ctx.Logger.Info("Tasks found in cache")
		return m.tickets[space], nil
	}

	m.ctx.Logger.Info("Fetching tasks from API")
	client := m.ctx.Clickup

	m.ctx.Logger.Infof("Getting views from space: %s", space)
	views, err := client.GetViewsFromSpace(space)
	if err != nil {
		return nil, err
	}
	m.ctx.Logger.Infof("Found %d views in space %s", len(views), space)

	m.ctx.Logger.Infof("Getting tasks from view ID: %s NAME: %s", views[0].Id, views[0].Name)
	tasks, err := client.GetTasksFromView(views[0].Id)
	if err != nil {
		return nil, err
	}
	m.ctx.Logger.Infof("Found %d tasks in view %s", len(tasks), views[0].Name)

	return tasks, nil
}
