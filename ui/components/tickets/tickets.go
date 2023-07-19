package tickets

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"

	"github.com/charmbracelet/bubbles/table"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

const (
	SPACE_SRE_LIST_COOL = "q5kna-61288"
	SPACE_SRE           = "48458830"
)

type Model struct {
	ctx               *context.UserContext
	table             table.Model
	columns           []table.Column
	tickets           map[string][]clickup.Task
	SelectedSpace     string
	PrevSelectedSpace string
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

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

	return Model{
		ctx:               ctx,
		table:             t,
		columns:           columns,
		tickets:           map[string][]clickup.Task{},
		SelectedSpace:     SPACE_SRE,
		PrevSelectedSpace: SPACE_SRE,
	}
}

func (m Model) syncTable(tasks []clickup.Task) Model {
	m.ctx.Logger.Info("Synchoning table")
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

	if m.SelectedSpace != m.PrevSelectedSpace {
		m.ctx.Logger.Info("space changed")
		m.PrevSelectedSpace = m.SelectedSpace
		tasks, err := m.getTickets(m.SelectedSpace)
		if err != nil {
			return m, common.ErrCmd(err)
		}
		m = m.syncTable(tasks)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case TasksListReloadedMsg:
		m.ctx.Logger.Info("TaskView receive TasksListReloadedMsg")
		m = m.syncTable(msg)

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.table.SetWidth(msg.Width - h)
		m.table.SetHeight(msg.Height - v)
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}

func (m Model) Init() tea.Msg {
	tasks, err := m.getTickets(m.SelectedSpace)
	if err != nil {
		return common.ErrMsg(err)
	}

	return TasksListReloadedMsg(tasks)
}

func (m Model) getTickets(space string) ([]clickup.Task, error) {
	m.ctx.Logger.Info("Getting tasks for space: " + space)
	if m.tickets[m.SelectedSpace] != nil {
		m.ctx.Logger.Info("Tasks found in cache")
		return m.tickets[m.SelectedSpace], nil
	}

	m.ctx.Logger.Info("Fetching tasks from API")
	client := m.ctx.Clickup

	m.ctx.Logger.Info("Getting views from space: " + space)
	views, err := client.GetViewsFromSpace(m.SelectedSpace)
	if err != nil {
		return nil, err
	}
	m.ctx.Logger.Info("Found %d views in space %s", len(views), space)

	m.ctx.Logger.Info("Getting tasks from view ID: %s NAME: %s", views[0].Id, views[0].Name)
	tasks, err := client.GetTasksFromView(views[0].Id)
	if err != nil {
		return nil, err
	}
	m.ctx.Logger.Info("Found %d tasks in view %s", len(tasks), views[0].Name)
	return tasks, nil
}

type TasksListReloadedMsg []clickup.Task
