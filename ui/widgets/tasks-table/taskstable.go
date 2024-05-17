package taskstable

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/evertras/bubble-table/table"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

const WidgetId = "widgetTasksTable"

type Model struct {
	tasks             []clickup.Task
	log               *log.Logger
	ctx               *context.UserContext
	WidgetId          common.WidgetId
	requiredColsKeys  []string
	columns           []table.Column
	requiredCols      []table.Column
	table             table.Model
	size              common.Size
	SelectedTaskIndex int
	autoColumns       bool
	Focused           bool
	Hidden            bool
	ifBorders         bool
}

func (m Model) SetSize(s common.Size) common.Widget {
	if m.ifBorders {
		s.Width -= 2  // two borders
		s.Height -= 2 // two borders
	}

	if m.size.Width == s.Width && m.size.Height == s.Height {
		return m
	}

	m.size = s
	m.log.Info("here")
	m.refreshTable()
	return m
}

func (m Model) KeyMap() help.KeyMap {
	km := m.table.KeyMap()

	switchFocusToListView := key.NewBinding(
		key.WithKeys("escape"),
		key.WithHelp("escape", "switch to list view"),
	)

	return common.NewKeyMap(
		func() [][]key.Binding {
			return [][]key.Binding{
				{
					common.KeyBindingWithHelp(km.RowDown, "down"),
					common.KeyBindingWithHelp(km.RowUp, "up"),
					common.KeyBindingWithHelp(km.RowSelectToggle, "select"),
				},
				{
					common.KeyBindingWithHelp(km.PageDown, "next page"),
					common.KeyBindingWithHelp(km.PageUp, "previous page"),
					common.KeyBindingWithHelp(km.PageFirst, "first page"),
					common.KeyBindingWithHelp(km.PageLast, "last page"),
				},
				{
					common.KeyBindingWithHelp(km.Filter, "filter"),
					common.KeyBindingWithHelp(km.FilterBlur, "filter blur"),
					common.KeyBindingWithHelp(km.FilterClear, "filter clear"),
				},
				{
					common.KeyBindingWithHelp(km.ScrollRight, "scroll right"),
					common.KeyBindingWithHelp(km.ScrollLeft, "scroll left"),
					common.KeyBindingOpenInBrowser,
					switchFocusToListView,
				},
			}
		},
		func() []key.Binding {
			return []key.Binding{
				common.KeyBindingWithHelp(km.RowDown, "down"),
				common.KeyBindingWithHelp(km.RowUp, "up"),
				common.KeyBindingWithHelp(km.RowSelectToggle, "select"),
				common.KeyBindingWithHelp(km.PageDown, "next page"),
				common.KeyBindingWithHelp(km.PageUp, "previous page"),
				common.KeyBindingOpenInBrowser,
				switchFocusToListView,
			}
		},
	)
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

	size := common.Size{
		Width:  0,
		Height: 0,
	}

	t := table.New(columns).
		WithTargetWidth(size.Width).
		SelectableRows(true).
		WithSelectedText(" ", "✓").
		Focused(true).
		WithPageSize(0).
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
		requiredCols:     requiredCols,
		requiredColsKeys: requiredColsKeys,
		tasks:            []clickup.Task{},
		autoColumns:      false,
		size:             size,
		Focused:          true,
		Hidden:           false,
		log:              log,
		ifBorders:        true,
	}
}

func (m *Model) RefreshTable() tea.Cmd {
	return m.refreshTable()
}

func (m *Model) GetColumnsKey() []string {
	r := make([]string, len(m.columns))
	for i, c := range m.columns {
		r[i] = c.Key()
	}

	return r
}

func (m *Model) refreshTable() tea.Cmd {
	m.log.Info("Synchonizing table...")
	tasks := m.tasks

	m.Hidden = false
	if len(tasks) == 0 {
		m.log.Info("Table is empty")
		m.Hidden = true
		return HideTableCmd()
	}

	items := taskListToRows(tasks, m.GetColumnsKey())

	m.SelectedTaskIndex = m.table.GetHighlightedRowIndex()

	pageSize := m.size.Height
	m.log.Infof("pageSize: %d", pageSize)
	if m.table.GetHeaderVisibility() {
		pageSize -= 1
	}
	if m.table.GetFooterVisibility() {
		pageSize -= 1
	}
	pageSize -= 3 // TODO: why 3? fix
	if pageSize < 0 {
		pageSize = 1
	}
	m.log.Infof("pageSize: %d", pageSize)
	m.table = m.table.
		WithRows(items).
		WithColumns(m.columns).
		WithTargetWidth(m.size.Width).
		WithMaxTotalWidth(m.size.Width).
		WithMinimumHeight(m.size.Height).
		WithPageSize(pageSize)

	m.log.Info("Table synchonized")

	return nil
}

func (m *Model) loadTasks(tasks []clickup.Task) {
	m.log.Info("Reloading table")
	m.tasks = tasks
}

func (m Model) Update(msg tea.Msg) (common.Widget, tea.Cmd) {
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
			taskId := m.tasks[index].Id
			m.log.Infof("Receive enter: %d", index)
			cmds = append(cmds, TaskSelectedCmd(taskId))
		case "p":
			index := m.table.GetHighlightedRowIndex()
			if m.table.TotalRows() == 0 {
				m.log.Info("Table is empty")
				break
			}
			task := m.tasks[index]
			m.log.Infof("Receive p: %d", index)
			if err := common.OpenUrlInWebBrowser(task.Url); err != nil {
				m.log.Fatal(err)
			}
		}
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	bColor := m.ctx.Theme.BordersColorInactive
	if m.Focused {
		bColor = m.ctx.Theme.BordersColorActive
	}

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(bColor).
		BorderBottom(m.ifBorders).
		BorderRight(m.ifBorders).
		BorderTop(m.ifBorders).
		BorderLeft(m.ifBorders).
		Render(
			m.table.View(),
		)
}

func (m Model) Init() tea.Cmd {
	m.log.Info("Initializing...")
	return m.refreshTable()
}

func (m Model) GetFocused() bool {
	return m.Focused
}

func (m Model) SetFocused(f bool) common.Widget {
	m.Focused = f
	return m
}

func (m Model) GetHidden() bool {
	return m.Hidden
}

func (m Model) SetHidden(h bool) common.Widget {
	m.Hidden = h
	return m
}

func (m Model) TabChanged(tabId string) (common.Widget, tea.Cmd) {
	var cmds []tea.Cmd

	m.log.Infof("Received TabChangedMsg: %s", tabId)

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

	m.log.Infof("Columns: %d", len(columns))
	m.columns = columns
	tasks := m.tasks

	m.loadTasks(tasks)
	cmds = append(cmds, m.refreshTable())

	if len(m.tasks) != 0 { // TODO: store tasks list in var
		taskId := m.tasks[m.SelectedTaskIndex].Id
		cmds = append(cmds, TaskSelectedCmd(taskId))
	}

	cmds = append(cmds, TasksListReadyCmd())

	return m, tea.Batch(cmds...)
}

func (m Model) FetchTasksForView(viewId string) (common.Widget, tea.Cmd) {
	m.log.Infof("Fetching tasks for the view: %s", viewId)
	tasks, err := m.ctx.Api.GetTasksFromView(viewId)
	if err != nil {
		return m, common.ErrCmd(err)
	}
	m.tasks = tasks

	return m, m.refreshTable()
}

func (m Model) FetchTasksForList(listId string) (common.Widget, tea.Cmd) {
	m.log.Infof("Fetching tasks for the list: %s", listId)
	tasks, err := m.ctx.Api.GetTasksFromList(listId)
	if err != nil {
		return m, common.ErrCmd(err)
	}
	m.tasks = tasks

	return m, m.refreshTable()
}
