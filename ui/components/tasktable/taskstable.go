package tasktable

import (
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

	SelectedView      string
	SelectedTaskIndex int

	Focused     bool
	autoColumns bool
	size        size
}

type size struct {
	Width  int
	Height int
}

func (m Model) getSelectedViewTaskIdByIndex(index int) string {
	return m.getSelectedViewTasks()[index].Id
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
	}

	size := size{
		Width:  0,
		Height: 0,
	}

	tablesStyles := table.DefaultStyles()
	tablesStyles.Cell = tablesStyles.Cell.Copy().
		PaddingBottom(1)

	tablesStyles.Header = tablesStyles.Header.Copy().
		PaddingBottom(1)

	t := table.New(
		table.WithFocused(true),
		table.WithHeight(size.Height),
		table.WithWidth(size.Width),
		table.WithStyles(tablesStyles),
		table.WithCellsWrap(true),
	)

	return Model{
		ctx:          ctx,
		table:        t,
		columns:      columns,
		requiredCols: requiredCols,
		tasks:        map[string][]clickup.Task{},
		SelectedView: SPACE_SRE_LIST_COOL,
		autoColumns:  false,
		size:         size,
	}
}

func (m *Model) refreshTable() {
	m.ctx.Logger.Info("Synchonizing table")
	tasks := m.getSelectedViewTasks()
	items := taskListToRows(tasks, m.columns)

	m.table.SetColumns(m.columns)
	m.table.SetRows(items)
	m.SelectedTaskIndex = m.table.Cursor()

	m.table.SetWidth(m.size.Width)
	m.table.SetHeight(m.size.Height)
}

func (m *Model) loadTasks(tasks []clickup.Task) {
	m.ctx.Logger.Info("Reloading table")
	m.tasks[m.SelectedView] = tasks
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			index := m.table.Cursor()
			taskId := m.getSelectedViewTaskIdByIndex(index)
			m.ctx.Logger.Infof("TaskTable receive enter: %d", index)
			cmds = append(cmds, TaskSelectedCmd(taskId))
		}

	case ViewChangedMsg:
		m.ctx.Logger.Infof("TaskTable receive ViewChangedMsg: %s", string(msg))
		m.SelectedView = string(msg)
		cmds = append(cmds, m.getTasksCmd(string(msg)))

	case TasksListReloadedMsg:
		m.ctx.Logger.Infof("TaskTable receive TasksListReloadedMsg: %d", len(msg))
		tasks := msg
		m.loadTasks(tasks)
		m.refreshTable()

		taskId := m.getSelectedViewTaskIdByIndex(m.SelectedTaskIndex)

		cmds = append(cmds, cmd, TasksListReadyCmd(taskId))
		// cmds = append(cmds, cmd, TasksListReadyCmd())
		// cmds = append(cmds, cmd, TasksListReadyCmd(), TaskSelectedCmd(taskId))

	case tea.WindowSizeMsg:
		m.ctx.Logger.Infof("TaskTable receive tea.WindowSizeMsg Width: %d Height %d", msg.Width, msg.Height)
		m.size.Width = int(0.4 * float32(m.ctx.WindowSize.Width))
		m.size.Height = int(0.7 * float32(m.ctx.WindowSize.Height))

		m.ctx.Logger.Infof("TaskTable set width: %d height: %d",
			m.size.Width, m.size.Height)

		m.refreshTable()

	case FetchTasksForViewMsg:
		m.ctx.Logger.Infof("TaskTable received FetchViewMsg: %s", string(msg))
		view := string(msg)
		tasks, err := m.ctx.Api.GetTasks(view)
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
	bColor := lipgloss.Color("#FFF")
	if m.Focused {
		bColor = lipgloss.Color("#8909FF")
	}

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(bColor).
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
	m.refreshTable()
	return nil
}
