package taskstable

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"

	"github.com/charmbracelet/bubbles/table"
)

type ViewLoadedMsg clickup.View

func ViewLoadedCmd(view clickup.View) tea.Cmd {
	return func() tea.Msg {
		return ViewLoadedMsg(view)
	}
}

type TasksListReady bool

func TasksListReadyCmd() tea.Cmd {
	return func() tea.Msg {
		return TasksListReady(true)
	}
}

type TasksListReloadedMsg []clickup.Task

func TasksListReloadedCmd(tasks []clickup.Task) tea.Cmd {
	return func() tea.Msg {
		return TasksListReloadedMsg(tasks)
	}
}

type ViewChangedMsg string

func ViewChangedCmd(space string) tea.Cmd {
	return func() tea.Msg {
		return ViewChangedMsg(space)
	}
}

type FetchTasksForViewMsg string

func FetchTasksForViewCmd(view string) tea.Cmd {
	return func() tea.Msg {
		return FetchTasksForViewMsg(view)
	}
}

type Model struct {
	ctx          *context.UserContext
	table        table.Model
	columns      []table.Column
	tickets      map[string][]clickup.Task
	SelectedView string
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

	return Model{
		ctx:          ctx,
		table:        t,
		columns:      columns,
		tickets:      map[string][]clickup.Task{},
		SelectedView: SPACE_SRE_LIST_COOL,
	}
}

func (m Model) syncTable(tasks []clickup.Task) Model {
	m.ctx.Logger.Info("Synchonizing table")
	m.tickets[m.SelectedView] = tasks

	m.table.SetColumns(m.columns)

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
	case ViewChangedMsg:
		m.ctx.Logger.Infof("TaskView receive ViewChangedMsg: %s", string(msg))
		m.SelectedView = string(msg)
		cmds = append(cmds, m.getTicketsCmd(string(msg)))

	case TasksListReloadedMsg:
		m.ctx.Logger.Infof("TaskView receive TasksListReloadedMsg: %d", len(msg))
		m = m.syncTable(msg)
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
		cmds = append(cmds, TasksListReadyCmd())

	case tea.WindowSizeMsg:
		m.ctx.Logger.Info("TaskView receive tea.WindowSizeMsg")
		h, v := docStyle.GetFrameSize()
		m.table.SetWidth(msg.Width - h)
		m.table.SetHeight(msg.Height - v)

	case common.FocusMsg:
		m.ctx.Logger.Info("TaskView received FocusMsg")

	case FetchTasksForViewMsg:
		m.ctx.Logger.Infof("TaskView received FetchViewMsg: %s", string(msg))
		view := string(msg)
		tasks, err := m.getTickets(view)
		if err != nil {
			return m, common.ErrCmd(err)
		}
		m.tickets[view] = tasks

	case ViewLoadedMsg:
		m.ctx.Logger.Infof("ViewsView received ViewLoadedMsg")
		view := clickup.View(msg)
		columnsNames := view.Columns.GetColumnsFields()
		columns := make([]table.Column, len(columnsNames))
		for i, name := range columnsNames {
			columns[i] = table.Column{Title: name, Width: 30}
		}
		m.columns = columns
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Info("Initializing component: TaskTable")
	return ViewChangedCmd(m.SelectedView)
}

func (m Model) getTicketsCmd(view string) tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.getTickets(view)
		if err != nil {
			return common.ErrMsg(err)
		}

		return TasksListReloadedMsg(tasks)
	}
}

func (m Model) getTickets(view string) ([]clickup.Task, error) {
	m.ctx.Logger.Infof("Getting tasks for view: %s", view)

	data, ok := m.ctx.Cache.Get("tasks", view)
	if ok {
		m.ctx.Logger.Infof("Tasks found in cache")
		var tasks []clickup.Task
		if err := m.ctx.Cache.ParseData(data, &tasks); err != nil {
			return nil, err
		}

		return tasks, nil
	}
	m.ctx.Logger.Info("Tasks not found in cache")

	m.ctx.Logger.Info("Fetching tasks from API")
	client := m.ctx.Clickup

	tasks, err := client.GetTasksFromView(view)
	if err != nil {
		return nil, err
	}
	m.ctx.Logger.Infof("Found %d tasks in view %s", len(tasks), view)

	m.ctx.Logger.Info("Caching tasks")
	m.ctx.Cache.Set("tasks", view, tasks)

	return tasks, nil
}
