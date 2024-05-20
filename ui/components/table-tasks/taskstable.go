package tabletasks

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

const ComponentId = "componenetTasksTable"

type Model struct {
	tasks             []clickup.Task
	log               *log.Logger
	ctx               *context.UserContext
	ComponentId       common.ComponentId
	columns           []table.Column
	table             table.Model
	size              common.Size
	SelectedTaskIndex int
	Focused           bool
	Hidden            bool
	ifBorders         bool
}

func (m Model) GetTasks() []clickup.Task {
	return m.tasks
}

func (m Model) GetSize() common.Size {
	return m.size
}

func (m *Model) GetTable() table.Model {
	return m.table
}

func (m *Model) SetSize(s common.Size) {
	m.size = s
	m.setTableSize(s)
}

func (m *Model) setTableSize(s common.Size) {
	pageSize := s.Height - 2
	if m.table.GetHeaderVisibility() {
		pageSize -= 2
	}

	if m.table.GetFooterVisibility() {
		pageSize -= 2
	}

	if pageSize < 0 {
		pageSize = 1
	}

	m.table = m.table.
		WithTargetWidth(s.Width).
		WithMaxTotalWidth(s.Width).
		WithMinimumHeight(s.Height).
		WithPageSize(pageSize)
}

func (m Model) KeyMap() help.KeyMap {
	km := m.table.KeyMap()

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
			}
		},
	)
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	columns := []table.Column{
		table.NewFlexColumn("name", "Name", 70),
		table.NewFlexColumn("status", "Status", 5).
			WithStyle(
				lipgloss.NewStyle().Align(lipgloss.Center),
			),
	}

	size := common.Size{
		Width:  0,
		Height: 0,
	}

	tableKeyMap := table.DefaultKeyMap()
	tableKeyMap.RowSelectToggle = key.NewBinding(
		key.WithKeys(" "),
	)

	t := table.New(columns).
		WithKeyMap(tableKeyMap).
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

	log := logger.WithPrefix(logger.GetPrefix() + "/" + ComponentId)

	return Model{
		ComponentId: ComponentId,
		ctx:         ctx,
		table:       t,
		columns:     columns,
		tasks:       []clickup.Task{},
		size:        size,
		Focused:     false,
		Hidden:      false,
		log:         log,
		ifBorders:   true,
	}
}

func (m *Model) GetColumnsKey() []string {
	r := make([]string, len(m.columns))
	for i, c := range m.columns {
		r[i] = c.Key()
	}

	return r
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
			taskId := m.tasks[index].Id
			m.log.Infof("Receive enter: %d", index)
			cmds = append(cmds, TaskSelectedCmd(taskId))
		}
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) GetHighlightedTask() *clickup.Task {
	index := m.table.GetHighlightedRowIndex()
	if m.table.TotalRows() == 0 {
		m.log.Info("Table is empty")
		return nil
	}

	return &m.tasks[index]
}

func (m Model) TotalRows() int {
	return m.table.TotalRows()
}

func (m Model) View() string {
	style := lipgloss.NewStyle()

	return style.Render(m.table.View())
}

func (m Model) Init() tea.Cmd {
	m.log.Info("Initializing...")
	return nil
}

func (m Model) GetFocused() bool {
	return m.Focused
}

func (m Model) SetFocused(f bool) Model {
	m.Focused = f
	return m
}

func (m Model) GetHidden() bool {
	return m.Hidden
}

func (m Model) SetHidden(h bool) Model {
	m.Hidden = h
	return m
}

func (m Model) TabChanged(tabId string) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	m.log.Infof("Received TabChangedMsg: %s", tabId)

	m.SetTasks(m.tasks)

	if len(m.tasks) != 0 { // TODO: store tasks list in var
		taskId := m.tasks[m.SelectedTaskIndex].Id
		cmds = append(cmds, TaskSelectedCmd(taskId))
	}

	cmds = append(cmds, TasksListReadyCmd())

	return m, tea.Batch(cmds...)
}

func (m *Model) setColumns(columns []table.Column) {
	m.columns = columns
	m.table = m.table.
		WithColumns(m.columns)
}

func (m *Model) SetTasks(tasks []clickup.Task) {
	m.tasks = tasks
	items := taskListToRows(tasks, m.GetColumnsKey())
	m.table = m.table.WithRows(items)
	m.log.Info("Table synchonized", "size", len(m.table.GetVisibleRows()))
}
