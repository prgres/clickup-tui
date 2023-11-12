package tasktable

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	autoColumns  bool
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
			table.WithHeight(0),
			table.WithWidth(0),
		),
		columns:      columns,
		requiredCols: requiredCols,
		tasks:        map[string][]clickup.Task{},
		SelectedView: SPACE_SRE_LIST_COOL,
		autoColumns:  false,
	}
}

func (m Model) syncTable(tasks []clickup.Task) Model {
	m.ctx.Logger.Info("Synchonizing table")
	m.tasks[m.SelectedView] = tasks

	items := taskListToRows(tasks, m.columns)

	m.table.SetColumns(m.columns)
	m.table.SetRows(items)

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

	case tea.WindowSizeMsg:
		m.ctx.Logger.Infof("TaskTable receive tea.WindowSizeMsg Width: %d Height %d", msg.Width, msg.Height)
		w := int(0.6 * float32(m.ctx.WindowSize.Width))
		h := int(0.7 * float32(m.ctx.WindowSize.Height))
		m.table.SetWidth(w)
		m.table.SetHeight(h)

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
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderBottom(true).
		BorderRight(true).
		BorderTop(true).
		BorderLeft(true).
		Render(
			m.table.View(),
		)
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
