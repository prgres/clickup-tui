package tasktable

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

type Model struct {
	ctx          *context.UserContext
	table        table.Model
	columns      []table.Column
	requiredCols []table.Column
	tickets      map[string][]clickup.Task
	SelectedView string

	autoColumns bool
}

func InitialModel(ctx *context.UserContext) Model {
	columns := []table.Column{}
	requiredCols := []table.Column{
		{
			Title: "id",
			Width: 0,
		},
		{
			Title: "name",
			Width: 30,
		},
		{
			Title: "status",
			Width: 30,
		},
		{
			Title: "assignee",
			Width: 30,
		},
	}

	return Model{
		ctx: ctx,
		table: table.New(
			table.WithFocused(true),
		),
		columns:      columns,
		requiredCols: requiredCols,
		tickets:      map[string][]clickup.Task{},
		SelectedView: SPACE_SRE_LIST_COOL,
		autoColumns:  false,
	}
}

func (m Model) syncTable(tasks []clickup.Task) Model {
	m.ctx.Logger.Info("Synchonizing table")
	m.tickets[m.SelectedView] = tasks

	items := taskListToRows(tasks, m.columns)
	// m.ctx.Logger.Infof("Values: %v", items)
	// m.ctx.Logger.Infof("Columns: %v", m.columns)

	m.table = table.New(
		table.WithColumns(m.columns),
		table.WithRows(items),
		table.WithFocused(true),
	)

	return m
}

func taskToRow(task clickup.Task, columns []table.Column) table.Row {
	values := table.Row{}
	for _, column := range columns {
		switch column.Title {
		case "status":
			values = append(values, task.Status.Status)
		case "name":
			values = append(values, task.Name)
		case "assignee":
			values = append(values, task.GetAssignees())
		case "list":
			values = append(values, task.List.String())
		case "tags":
			values = append(values, task.GetTags())
		case "folder":
			values = append(values, task.Folder.String())
		case "url":
			values = append(values, task.Url)
		case "space":
			values = append(values, task.Space.Id)
		case "id":
			values = append(values, task.Id)
		default:
			values = append(values, "XXX")
		}
	}

	return values
}

func taskListToRows(tasks []clickup.Task, columns []table.Column) []table.Row {
	rows := make([]table.Row, len(tasks))
	for i, task := range tasks {
		rows[i] = taskToRow(task, columns)
	}
	return rows
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			row := m.table.SelectedRow()

			m.ctx.Logger.Infof("TaskTable receive enter: %v", row)
			cmds = append(cmds, TaskSelectedCmd(strings.Join(row, " ")))
		}

	case ViewChangedMsg:
		m.ctx.Logger.Infof("TaskTable receive ViewChangedMsg: %s", string(msg))
		m.SelectedView = string(msg)
		cmds = append(cmds, m.getTicketsCmd(string(msg)))

	case TasksListReloadedMsg:
		m.ctx.Logger.Infof("TaskTable receive TasksListReloadedMsg: %d", len(msg))
		m = m.syncTable(msg)
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
		cmds = append(cmds, TasksListReadyCmd())

	case common.WindowSizeMsg:
		m.ctx.Logger.Info("TaskTable receive tea.WindowSizeMsg")
		h, v := docStyle.GetFrameSize()
		m.table.SetWidth(msg.Width - h)
		m.table.SetHeight(msg.Height - v)

	case common.FocusMsg:
		m.ctx.Logger.Info("TaskTable received FocusMsg")

	case FetchTasksForViewMsg:
		m.ctx.Logger.Infof("TaskTable received FetchViewMsg: %s", string(msg))
		view := string(msg)
		tasks, err := m.getTickets(view)
		if err != nil {
			return m, common.ErrCmd(err)
		}
		m.tickets[view] = tasks

	case ViewLoadedMsg:
		m.ctx.Logger.Infof("TaskTable received ViewLoadedMsg")
		view := clickup.View(msg)

		columns := []table.Column{}
		columns = append(columns, m.requiredCols...)

		if m.autoColumns {
			for _, field := range view.Columns.Fields {
				if field.Field == "name" || field.Field == "status" { // TODO: check if in requiredCols
					continue
				}
				columns = append(columns, table.Column{
					Title: field.Field,
					Width: 30,
				})
			}
		}

		m.ctx.Logger.Infof("Columns: %d", len(columns))
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
