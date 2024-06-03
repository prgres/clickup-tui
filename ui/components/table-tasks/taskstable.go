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

const id = "tasks-table"

type Model struct {
	id             common.Id
	tasks          []clickup.Task
	log            *log.Logger
	ctx            *context.UserContext
	columns        []Column
	columnsVisible []Column
	columnsHidden  []Column
	table          table.Model
	size           common.Size
	Focused        bool
	Hidden         bool
	ifBorders      bool
	keyMap         KeyMap

	SelectedIdx int
}

func (m Model) Id() common.Id {
	return m.id
}

type Column struct {
	table.Column
	Hidden bool
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

type KeyMap struct {
	table.KeyMap
	Select key.Binding
}

func (m Model) KeyMap() KeyMap {
	return m.keyMap
}

func DefaultKeyMap() KeyMap {
	km := table.DefaultKeyMap()

	return KeyMap{
		KeyMap: table.KeyMap{
			RowDown:         common.KeyBindingWithHelp(km.RowDown, "down"),
			RowUp:           common.KeyBindingWithHelp(km.RowUp, "up"),
			RowSelectToggle: common.KeyBindingWithHelp(km.RowSelectToggle, "select"),
			PageDown:        common.KeyBindingWithHelp(km.PageDown, "next page"),
			PageUp:          common.KeyBindingWithHelp(km.PageUp, "previous page"),
			PageFirst:       common.KeyBindingWithHelp(km.PageFirst, "first page"),
			PageLast:        common.KeyBindingWithHelp(km.PageLast, "last page"),
			Filter:          common.KeyBindingWithHelp(km.Filter, "filter"),
			FilterBlur:      common.KeyBindingWithHelp(km.FilterBlur, "filter blur"),
			FilterClear:     common.KeyBindingWithHelp(km.FilterClear, "filter clear"),
			ScrollRight:     common.KeyBindingWithHelp(km.ScrollRight, "scroll right"),
			ScrollLeft:      common.KeyBindingWithHelp(km.ScrollLeft, "scroll left"),
		},
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

func (m Model) Help() help.KeyMap {
	km := m.keyMap

	return common.NewHelp(
		func() [][]key.Binding {
			return [][]key.Binding{
				{
					km.RowDown,
					km.RowUp,
					km.RowSelectToggle,
				},
				{
					km.PageDown,
					km.PageUp,
					km.PageFirst,
					km.PageLast,
				},
				{
					km.Filter,
					km.FilterBlur,
					km.FilterClear,
				},
				{
					km.ScrollRight,
					km.ScrollLeft,
					m.keyMap.Select,
				},
			}
		},
		func() []key.Binding {
			return []key.Binding{
				km.RowDown,
				km.RowUp,
				km.RowSelectToggle,
				km.PageDown,
				km.PageUp,
				m.keyMap.Select,
			}
		},
	)
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	// TODO: do better
	columnsHidden := []Column{
		{
			Column: table.NewFlexColumn("url", "url", 0),
			Hidden: true,
		},
		{
			Column: table.NewFlexColumn("id", "id", 0),
			Hidden: true,
		},
	}
	columnsVisible := []Column{
		{
			Column: table.NewFlexColumn("name", "Name", 70),
			Hidden: false,
		},
		{
			Column: table.NewFlexColumn("status", "Status", 5).
				WithStyle(lipgloss.NewStyle().Align(lipgloss.Center)),
			Hidden: false,
		},
	}
	columns := append(columnsVisible, columnsHidden...)

	tableColumns := make([]table.Column, len(columnsVisible))
	for i := range columnsVisible {
		tableColumns[i] = columns[i].Column
	}

	size := common.NewEmptySize()

	tableKeyMap := table.DefaultKeyMap()
	tableKeyMap.RowSelectToggle = key.NewBinding(
		key.WithKeys(" "),
	)

	t := table.New(tableColumns).
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

	log := common.NewLogger(logger, common.ResourceTypeRegistry.COMPONENT, id)

	return Model{
		id:             id,
		ctx:            ctx,
		table:          t,
		columns:        columns,
		columnsVisible: columnsVisible,
		columnsHidden:  columnsHidden,
		tasks:          []clickup.Task{},
		size:           size,
		Focused:        false,
		Hidden:         false,
		log:            log,
		ifBorders:      true,
		keyMap:         DefaultKeyMap(),
	}
}

func (m *Model) GetVisibleColumnsKey() []string {
	return m.getColumnsKey(m.columnsVisible)
}

func (m *Model) GetColumnsKey() []string {
	return m.getColumnsKey(m.columns)
}

func (m *Model) getColumnsKey(cols []Column) []string {
	r := make([]string, len(cols))
	for i := range cols {
		r[i] = cols[i].Key()
	}

	return r
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Select):
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

	return tea.Batch(cmds...)
}

func (m Model) GetHighlightedTask() *clickup.Task {
	index := m.table.GetHighlightedRowIndex()
	if m.table.TotalRows() == 0 {
		m.log.Info("Table is empty")
		return nil
	}

	return &m.tasks[index]
}

func (m Model) GetSelectedTasks() []*clickup.Task {
	rows := m.table.SelectedRows()
	tasks := make([]*clickup.Task, len(rows))

	for i := range rows {
		task := rowToTask(rows[i], m.GetColumnsKey())
		tasks[i] = &task
	}

	return tasks
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

func (m Model) WithFocused(f bool) Model {
	m.Focused = f
	return m
}

func (m *Model) SetFocused(f bool) *Model {
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
		taskId := m.tasks[m.SelectedIdx].Id
		cmds = append(cmds, TaskSelectedCmd(taskId))
	}

	cmds = append(cmds, TasksListReadyCmd())

	return m, tea.Batch(cmds...)
}

func (m *Model) SetTasks(tasks []clickup.Task) {
	m.tasks = tasks
	items := taskListToRows(tasks, m.GetColumnsKey())
	m.table = m.table.WithRows(items)
	m.log.Info("Table synchonized", "size", len(m.table.GetVisibleRows()))
}
