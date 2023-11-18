package tasktable

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
	"github.com/prgrs/clickup/ui/widgets/viewtabs"
)

type Model struct {
	ctx               *context.UserContext
	table             table.Model
	columns           []table.Column
	requiredCols      []table.Column
	tasks             map[string][]clickup.Task
	SelectedTab       viewtabs.Tab // TODO: make it SelectedTabId
	SelectedTaskIndex int
	Focused           bool
	autoColumns       bool
	size              size
}

type size struct {
	Width  int
	Height int
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
		// SelectedView: "",
		autoColumns: false,
		size:        size,
		Focused:     true,
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
	m.tasks[m.SelectedTab.Id] = tasks
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

	case TabChangedMsg:
		tab := viewtabs.Tab(msg)
		m.ctx.Logger.Infof("TaskTable receive TabChangedMsg: %s", tab.Name)

		columns := []table.Column{}
		columns = append(columns, m.requiredCols...)

		// if m.autoColumns {
		//      tab := viewtabs.Tab(msg)
		// 	for _, field := range view.Columns.Fields {
		// 		if field.Field == "name" || field.Field == "status" { // TODO: check if in requiredCols
		// 			continue
		// 		}
		// 		columns = append(columns, table.Column{
		// 			Title: field.Field,
		// 			Width: 30,
		// 		})
		// 	}
		// }
		m.ctx.Logger.Infof("Columns: %d", len(columns))
		m.columns = columns

		m.SelectedTab = tab
		tasks := m.tasks[tab.Id]
		cmds = append(cmds, TasksListReloadedCmd(tasks))

		// case TasksListReloadedMsg:
		// m.ctx.Logger.Infof("TaskTable receive TasksListReloadedMsg: %d", len(msg))
		// tasks := msg
		m.loadTasks(tasks)
		m.refreshTable()

		if len(m.tasks[m.SelectedTab.Id]) != 0 { //TODO: store tasks list in var
			taskId := m.getSelectedViewTaskIdByIndex(m.SelectedTaskIndex)
			cmds = append(cmds, TaskSelectedCmd(taskId))
		}

		cmds = append(cmds, cmd, TasksListReadyCmd())

	case tea.WindowSizeMsg:
		m.ctx.Logger.Infof("TaskTable receive tea.WindowSizeMsg Width: %d Height %d", msg.Width, msg.Height)
		m.size.Width = int(0.4 * float32(m.ctx.WindowSize.Width))
		m.size.Height = int(0.7 * float32(m.ctx.WindowSize.Height))

		m.ctx.Logger.Infof("TaskTable set width: %d height: %d",
			m.size.Width, m.size.Height)

		m.refreshTable()

	case viewtabs.FetchTasksForTabsMsg:
		m.ctx.Logger.Infof("TaskTable receive viewtabs.FetchTasksForTabsMsg")
		tabs := []viewtabs.Tab(msg)
		for _, tab := range tabs {
			m.ctx.Logger.Infof("TaskTable received FetchTasksForTabsMsg: type %v", tab.Type)
			switch tab.Type {
			case "list":
				if err := m.fetchTasksForList(tab.Id); err != nil {
					return m, common.ErrCmd(err)
				}
			case "view":
				if err := m.fetchTasksForViewId(tab.Id); err != nil {
					return m, common.ErrCmd(err)
				}
			}
		}
		m.refreshTable()
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

func (m *Model) fetchTasksForViewId(viewId string) error {
	m.ctx.Logger.Infof("TaskTable: fetching tasks for the view: %s", viewId)
	tasks, err := m.ctx.Api.GetTasksFromView(viewId)
	if err != nil {
		return err
	}
	m.tasks[viewId] = tasks
	return nil
}

func (m *Model) fetchTasksForList(listId string) error {
	m.ctx.Logger.Infof("TaskTable: fetching tasks for the list: %s", listId)
	tasks, err := m.ctx.Api.GetTasksFromList(listId)
	if err != nil {
		return err
	}
	m.tasks[listId] = tasks
	return nil
	// 	m.refreshTable()
}
