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
	tasks        map[string][]clickup.Task

	SelectedView string

	autoColumns bool
}

func (m Model) getSelectedViewTasks() []clickup.Task {
	return m.tasks[m.SelectedView]
}

func InitialModel(ctx *context.UserContext) Model {
	columns := []table.Column{}
	requiredCols := []table.Column{
		{
			Title: "name",
			Width: 40,
		},
		{
			Title: "status",
			Width: 15,
		},
		{
			Title: "assignee",
			Width: 40,
		},
	}

	return Model{
		ctx: ctx,
		table: table.New(
			table.WithFocused(true),
		),
		columns:      columns,
		requiredCols: requiredCols,
		tasks:      map[string][]clickup.Task{},
		SelectedView: SPACE_SRE_LIST_COOL,
		autoColumns:  false,
	}
}

func (m Model) syncTable(tasks []clickup.Task) Model {
	m.ctx.Logger.Info("Synchonizing table")
	m.tasks[m.SelectedView] = tasks

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

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			index := m.table.Cursor()
			task := m.getSelectedViewTasks()[index]
			m.ctx.Logger.Infof("TaskTable receive enter: %d", index)
			cmds = append(cmds, TaskSelectedCmd(task.Id))
		}

	case ViewChangedMsg:
		m.ctx.Logger.Infof("TaskTable receive ViewChangedMsg: %s", string(msg))
		m.SelectedView = string(msg)
		cmds = append(cmds, m.getTasksCmd(string(msg)))

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
		tasks, err := m.getTasks(view)
		if err != nil {
			return m, common.ErrCmd(err)
		}
		m.tasks[view] = tasks

	case common.ViewLoadedMsg:
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
	tasks, err := m.getTasks(m.SelectedView)
	if err != nil {
		return common.ErrCmd(err)
	}
	var cmd tea.Cmd
	m = m.syncTable(tasks)
	m.table, cmd = m.table.Update(tasks)
	row := m.table.SelectedRow()
	return tea.Batch(cmd, TasksListReadyCmd(), TaskSelectedCmd(strings.Join(row, " ")))
}
