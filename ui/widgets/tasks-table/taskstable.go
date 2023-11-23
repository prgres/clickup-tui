package taskstable

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/evertras/bubble-table/table"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
	"github.com/prgrs/clickup/ui/widgets/tasks-tabs"
)

const WidgetId = "widgetTasksTable"

type Model struct {
	WidgetId          common.WidgetId
	Hidden            bool
	SelectedTab       taskstabs.Tab
	SelectedTaskIndex int
	Focused           bool

	ctx          *context.UserContext
	table        table.Model
	columns      []table.Column
	requiredCols []table.Column
	tasks        map[string][]clickup.Task
	autoColumns  bool
	size         size
	log          *log.Logger

	// TODO: ugly hack since table.Column does not expose any Getters. Waits for https://github.com/Evertras/bubble-table/issues/157
	columnsKeys      []string
	requiredColsKeys []string
}

type size struct {
	Width  int
	Height int
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	columns := []table.Column{}

	requiredCols := []table.Column{
		table.NewFlexColumn("name", "Name", 70),
		table.NewFlexColumn("status", "Status", 5).
			WithStyle(
				lipgloss.NewStyle().Align(lipgloss.Center),
			),
	}
	requiredColsKeys := []string{"name", "status"}

	size := size{
		Width:  0,
		Height: 0,
	}

	t := table.New(columns).
		WithTargetWidth(size.Width).
		SelectableRows(true).
		WithSelectedText(" ", "✓").
		Focused(true).
		WithBaseStyle(
			lipgloss.NewStyle().
				Align(lipgloss.Left),
		).
		HighlightStyle(
			lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("212")))

	log := logger.WithPrefix(logger.GetPrefix() + "/" + WidgetId)

	return Model{
		WidgetId:         WidgetId,
		ctx:              ctx,
		table:            t,
		columns:          columns,
		columnsKeys:      []string{},
		requiredCols:     requiredCols,
		requiredColsKeys: requiredColsKeys,
		tasks:            map[string][]clickup.Task{},
		autoColumns:      false,
		size:             size,
		Focused:          true,
		Hidden:           false,
		log:              log,
	}
}

func (m *Model) refreshTable() tea.Cmd {
	m.log.Info("Synchonizing table...")
	tasks := m.getSelectedViewTasks()

	m.Hidden = false
	if len(tasks) == 0 {
		m.log.Info("Table is empty")
		m.Hidden = true
		return HideTableCmd()
	}

	items := taskListToRows(tasks, m.columnsKeys)

	m.SelectedTaskIndex = m.table.GetHighlightedRowIndex()

	m.table = m.table.
		WithRows(items).
		WithColumns(m.columns).
		SelectableRows(true).
		WithTargetWidth(m.size.Width).
		WithPageSize(m.size.Height).
		WithMinimumHeight(m.size.Height)

	m.log.Info("Table synchonized")

	return nil
}

func (m *Model) loadTasks(tasks []clickup.Task) {
	m.log.Info("Reloading table")
	m.tasks[m.SelectedTab.Id] = tasks
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			index := m.table.GetHighlightedRowIndex()
			if m.table.TotalRows() == 0 {
				m.log.Info("Table is empty")
				break
			}
			taskId := m.getSelectedViewTaskIdByIndex(index)
			m.log.Infof("Receive enter: %d", index)
			cmds = append(cmds, TaskSelectedCmd(taskId))
		}

	case TabChangedMsg:
		tab := taskstabs.Tab(msg)
		m.log.Infof("Received TabChangedMsg: %s", tab.Name)

		columns := []table.Column{}
		columns = append(columns, m.requiredCols...)

		columnsKeys := []string{}
		columnsKeys = append(columnsKeys, m.requiredColsKeys...)

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

		m.log.Infof("Columns: %d", len(columns))
		m.columns = columns
		m.columnsKeys = columnsKeys

		m.SelectedTab = tab
		tasks := m.tasks[tab.Id]

		m.loadTasks(tasks)
		cmd = m.refreshTable()
		cmds = append(cmds, cmd)

		if len(m.tasks[m.SelectedTab.Id]) != 0 { //TODO: store tasks list in var
			taskId := m.getSelectedViewTaskIdByIndex(m.SelectedTaskIndex)
			cmds = append(cmds, TaskSelectedCmd(taskId))
		}

		cmds = append(cmds, TasksListReadyCmd())

	case tea.WindowSizeMsg:
		m.log.Debug("Received: tea.WindowSizeMsg",
			"width", msg.Width,
			"height", msg.Height)
		m.size.Width = int(0.4 * float32(m.ctx.WindowSize.Width))
		m.size.Height = int(0.7 * float32(m.ctx.WindowSize.Height))

		cmds = append(cmds, m.refreshTable())

	case taskstabs.FetchTasksForTabsMsg:
		m.log.Infof("Received: viewtabs.FetchTasksForTabsMsg")
		tabs := []taskstabs.Tab(msg)
		for _, tab := range tabs {
			m.log.Infof("Received: FetchTasksForTabsMsg: type %v", tab.Type)
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
		cmd = m.refreshTable()
		cmds = append(cmds, cmd)
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
	m.log.Info("Initializing...")
	return m.refreshTable()
}

func (m *Model) fetchTasksForViewId(viewId string) error {
	m.log.Infof("Fetching tasks for the view: %s", viewId)
	tasks, err := m.ctx.Api.GetTasksFromView(viewId)
	if err != nil {
		return err
	}
	m.tasks[viewId] = tasks
	return nil
}

func (m *Model) fetchTasksForList(listId string) error {
	m.log.Infof("Fetching tasks for the list: %s", listId)
	tasks, err := m.ctx.Api.GetTasksFromList(listId)
	if err != nil {
		return err
	}
	m.tasks[listId] = tasks
	return nil
}
