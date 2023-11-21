package taskstable

import (
	// "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
	"github.com/prgrs/clickup/ui/widgets/tasks-tabs"
)

const WidgetId = "widgetTasksTable"

type Model struct {
	WidgetId          common.WidgetId
	ctx               *context.UserContext
	table             *table.Table
	tableData         table.Data
	tableHeigh        int
	columns           []string
	requiredCols      []string
	tasks             map[string][]clickup.Task
	SelectedTab       taskstabs.Tab
	SelectedTaskIndex int
	Focused           bool
	autoColumns       bool
	size              size
	Hidden            bool
	log               *log.Logger
}

type size struct {
	Width  int
	Height int
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	columns := []string{}
	requiredCols := []string{
		"name",
		"status",
	}

	size := size{
		Width:  0,
		Height: 0,
	}

	log := logger.WithPrefix(logger.GetPrefix() + "/" + WidgetId)

	m := Model{
		WidgetId:          WidgetId,
		ctx:               ctx,
		table:             table.New(),
		columns:           columns,
		requiredCols:      requiredCols,
		tasks:             map[string][]clickup.Task{},
		SelectedTaskIndex: 1,
		autoColumns:       false,
		size:              size,
		Focused:           true,
		Hidden:            false,
		log:               log,

		// SelectedView: "",
	}

	// m.table = *t

	return m
}
func (m *Model) renderTable() string {
	var (
		tablesStyles = lipgloss.NewStyle().
				PaddingBottom(1)

		tablesStylesCellActive = tablesStyles.Copy().
					Bold(true).
					Foreground(lipgloss.Color("212"))

		tablesStylesCellInactive = tablesStyles.Copy()

		tablesStylesHeader = tablesStyles.Copy().
					Bold(true)
	)

	return m.table.
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			switch row {
			case 0:
				return tablesStylesHeader
			case m.SelectedTaskIndex:
				return tablesStylesCellActive
			default:
				style = tablesStylesCellInactive
			}

			// Make the second column a little wider.
			if col == 1 {
				style = style.Copy().Width(22)
			}

			return style
		}).String()
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

	items := taskListToRows(tasks, m.columns)
	m.tableData = table.NewStringData(items...)

	for i := 0; i < m.size.Height-7-m.tableData.Rows()-1*m.tableData.Rows(); i++ {
		items = append(items, []string{})
	}

	m.table = m.table.Data(table.NewStringData(items...))
	m.table = m.table.Headers(m.columns...)

	m.SelectedTaskIndex = 1 // 0 is header

	m.table.Width(m.size.Width)
	m.table.Height(m.size.Height)

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
		case "j", "down":
			if m.SelectedTaskIndex+1 <= m.tableData.Rows() {
				m.SelectedTaskIndex++
			}

		case "k", "up":
			if m.SelectedTaskIndex-1 > 0 {
				m.SelectedTaskIndex--
			}

		case "enter":
			index := m.SelectedTaskIndex
			if m.tableData == nil || m.tableData.Rows() == 0 {
				m.log.Info("Table is empty")
				break
			}
			m.log.Infof("Receive enter: %d", index)
			taskId := m.getSelectedViewTaskIdByIndex(index)
			cmds = append(cmds, TaskSelectedCmd(taskId))
		}

	case TabChangedMsg:
		tab := taskstabs.Tab(msg)
		m.log.Infof("Received TabChangedMsg: %s", tab.Name)

		columns := []string{}
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

		m.log.Infof("Columns: %d", len(columns))
		m.columns = columns

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
		m.log.Info("Received: tea.WindowSizeMsg",
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
			m.renderTable(),
			// m.table.String(),
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
