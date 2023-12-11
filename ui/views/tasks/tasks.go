package tasks

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
	taskssidebar "github.com/prgrs/clickup/ui/widgets/tasks-sidebar"
	taskstable "github.com/prgrs/clickup/ui/widgets/tasks-table"
	taskstabs "github.com/prgrs/clickup/ui/widgets/tasks-tabs"
)

const ViewId = "viewTasks"

type TasksState uint

const (
	TasksStateLoading TasksState = iota
	TasksStateTasksTable
	TasksStateViewsTabs
	TasksStateTaskSidebar
)

type Model struct {
	widgetViewsTabs   common.Widget
	widgetTasksTable  common.Widget
	widgetTaskSidebar common.Widget
	ctx               *context.UserContext
	log               *log.Logger
	ViewId            common.ViewId
	widgetsList       []common.Widget
	spinner           spinner.Model
	size              common.Size
	state             TasksState
	showSpinner       bool
	ifBorders         bool
}

func (m Model) Ready() bool {
	return !m.showSpinner
}

func (m Model) KeyMap() help.KeyMap {
	var km help.KeyMap
	var ()
	switch m.state {
	case TasksStateTasksTable:
		km = m.widgetTasksTable.KeyMap()
	case TasksStateTaskSidebar:
		km = m.widgetTaskSidebar.KeyMap()
	case TasksStateViewsTabs:
		km = m.widgetViewsTabs.KeyMap()
	}
	return km
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) common.View {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	log := logger.WithPrefix(logger.GetPrefix() + "/" + ViewId)

	var (
		widgetViewsTabs   = taskstabs.InitialModel(ctx, log)
		widgetTasksTable  = taskstable.InitialModel(ctx, log)
		widgetTaskSidebar = taskssidebar.InitialModel(ctx, log)
	)

	widgetsList := []common.Widget{
		widgetViewsTabs,
		widgetTasksTable,
		widgetTaskSidebar,
	}

	return Model{
		ViewId:            ViewId,
		ctx:               ctx,
		state:             TasksStateTasksTable,
		widgetViewsTabs:   widgetViewsTabs,
		widgetTasksTable:  widgetTasksTable,
		widgetTaskSidebar: widgetTaskSidebar,
		widgetsList:       widgetsList,
		spinner:           s,
		showSpinner:       true,
		log:               log,
		ifBorders:         true,
		size: common.Size{
			Width:  0,
			Height: 0,
		},
	}
}

func (m Model) Update(msg tea.Msg) (common.View, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {

		case "esc":
			switch m.state {
			case TasksStateTaskSidebar:
				m.state = TasksStateTasksTable
				m.widgetTasksTable = m.widgetTasksTable.SetFocused(true)
				m.widgetTaskSidebar = m.widgetTaskSidebar.SetFocused(false)
				m.widgetViewsTabs = m.widgetViewsTabs.SetFocused(false)

			case TasksStateViewsTabs:
				m.state = TasksStateTasksTable
				m.widgetTasksTable = m.widgetTasksTable.SetFocused(true)
				m.widgetTaskSidebar = m.widgetTaskSidebar.SetFocused(false)
				m.widgetViewsTabs = m.widgetViewsTabs.SetFocused(false)

			case TasksStateTasksTable:
				m.state = TasksStateTasksTable
				m.log.Info("Received: Go to previous view")
				cmds = append(cmds, common.BackToPreviousViewCmd(m.ViewId))
			}

		default:
			switch keypress {
			case "h", "left", "l", "right":
				if m.state == TasksStateTasksTable {
					m.state = TasksStateViewsTabs
					m.widgetTasksTable = m.widgetTasksTable.SetFocused(false)
					m.widgetTaskSidebar = m.widgetTaskSidebar.SetFocused(false)
					m.widgetViewsTabs = m.widgetViewsTabs.SetFocused(true)
				}

			case "j", "down", "k", "up":
				if m.state == TasksStateViewsTabs {
					m.state = TasksStateTasksTable
					m.widgetTasksTable = m.widgetTasksTable.SetFocused(true)
					m.widgetTaskSidebar = m.widgetTaskSidebar.SetFocused(false)
					m.widgetViewsTabs = m.widgetViewsTabs.SetFocused(false)
				}
			}

			switch m.state {
			case TasksStateTasksTable:
				m.widgetTasksTable, cmd = m.widgetTasksTable.Update(msg)
				cmds = append(cmds, cmd)

			case TasksStateTaskSidebar:
				m.widgetTaskSidebar, cmd = m.widgetTaskSidebar.Update(msg)
				cmds = append(cmds, cmd)

			case TasksStateViewsTabs:
				m.widgetViewsTabs, cmd = m.widgetViewsTabs.Update(msg)
				cmds = append(cmds, cmd)

			default:
				m.log.Infof("Received: unhandled keypress: %s", keypress)
			}

			return m, tea.Batch(cmds...)
		}

	case taskstabs.TabChangedMsg:
		tab := taskstabs.Tab(msg)
		m.log.Info("Received: TabChangedMsg", "name", tab.Name, "id", tab.Id)
		m.showSpinner = true

		m.widgetViewsTabs, cmd = m.widgetViewsTabs.Update(msg)

		cmds = append(cmds,
			cmd,
			taskstable.TabChangedCmd(tab),
			m.spinner.Tick,
		)

	case taskstable.TasksListReadyMsg:
		m.log.Info("Received: TasksListReady")
		m.showSpinner = false
		// m = m.Resize(m.size)
		cmds = append(cmds,
			m.spinner.Tick,
		)

	case spinner.TickMsg:
		// m.log.Info("ViewTask spinner.TickMsg")
		if m.showSpinner {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case taskstable.TaskSelectedMsg:
		id := string(msg)
		m.log.Infof("Received: taskstable.TaskSelectedMsg: %s", id)
		if m.state == TasksStateTasksTable {
			m.state = TasksStateTaskSidebar
			m.widgetTasksTable = m.widgetTasksTable.SetFocused(false)
			m.widgetTaskSidebar = m.widgetTaskSidebar.SetFocused(true)
		}
		cmds = append(cmds, taskssidebar.TaskSelectedCmd(id))
	}

	m.widgetViewsTabs, cmd = m.widgetViewsTabs.Update(msg)
	cmds = append(cmds, cmd)

	m.widgetTasksTable, cmd = m.widgetTasksTable.Update(msg)
	cmds = append(cmds, cmd)

	m.widgetTaskSidebar, cmd = m.widgetTaskSidebar.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) Resize(size common.Size) Model {
	if m.ifBorders {
		size.Width -= 2
		size.Height -= 2
	}

	tableWidth := int(0.4 * float32(size.Width))
	m.widgetTasksTable = m.widgetTasksTable.SetSize(common.Size{
		Width:  tableWidth,
		Height: size.Height,
	})

	widgetTasksTableRendered := m.widgetTasksTable.View()
	widgetTasksTableWidth := lipgloss.Width(widgetTasksTableRendered)
	widgetTasksTableHeight := lipgloss.Height(widgetTasksTableRendered)

	m.widgetTaskSidebar = m.widgetTaskSidebar.SetSize(common.Size{
		Width:  size.Width - widgetTasksTableWidth,
		Height: widgetTasksTableHeight,
	})

	return m
}

func (m Model) View() string {
	if m.showSpinner {
		return lipgloss.Place(
			m.ctx.WindowSize.Width,
			m.ctx.WindowSize.Height,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("%s Loading tasks...",
				m.spinner.View()),
		)
	}

	widgetViewsTabsRendered := m.widgetViewsTabs.View()

	tableAndSidebar := ""
	if m.widgetTasksTable.GetHidden() {
		tableAndSidebar = lipgloss.Place(
			int(0.9*float32(m.ctx.WindowSize.Width)),
			m.size.Height,
			// availableHeight,
			lipgloss.Center, lipgloss.Center,
			"No tasks founds",
		)
	} else {
		tableAndSidebar = lipgloss.JoinHorizontal(
			lipgloss.Left,
			m.widgetTasksTable.View(),
			m.widgetTaskSidebar.View(),
		)
	}

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderRight(m.ifBorders).
		BorderBottom(m.ifBorders).
		BorderTop(m.ifBorders).
		BorderLeft(m.ifBorders).
		Render(lipgloss.JoinVertical(
			lipgloss.Top,
			widgetViewsTabsRendered,
			tableAndSidebar,
		))
}

func (m Model) Init() tea.Cmd {
	m.log.Infof("Initializing...")

	cmds := make([]tea.Cmd, len(m.widgetsList))
	for i, w := range m.widgetsList {
		cmds[i] = w.Init()
	}

	cmds = append(cmds, m.spinner.Tick)

	return tea.Batch(
		cmds...,
	)
}

func (m Model) SetSize(s common.Size) common.View {
	availableHeight := s.Height
	// availableHeight := s.Height - m.ctx.WindowSize.MetaHeight
	// availableHeight := m.ctx.WindowSize.Height - m.ctx.WindowSize.MetaHeight
	widgetViewsTabsRendered := m.widgetViewsTabs.View()
	availableHeight -= lipgloss.Height(widgetViewsTabsRendered)
	// availableHeight -= 1
	m = m.Resize(common.Size{
		Width:  s.Width,
		Height: availableHeight,
	})
	m.size = s
	return m
}

func (m Model) GetSize() common.Size {
	return m.size
}

func (m Model) GetViewId() common.ViewId {
	return m.ViewId
}
