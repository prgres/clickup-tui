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
	ViewId            common.ViewId
	ctx               *context.UserContext
	state             TasksState
	widgetViewsTabs   taskstabs.Model
	widgetTasksTable  taskstable.Model
	widgetTaskSidebar taskssidebar.Model
	spinner           spinner.Model
	showSpinner       bool
	log               *log.Logger
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

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	log := logger.WithPrefix(logger.GetPrefix() + "/" + ViewId)

	return Model{
		ViewId:            ViewId,
		ctx:               ctx,
		state:             TasksStateTasksTable,
		widgetViewsTabs:   taskstabs.InitialModel(ctx, log),
		widgetTasksTable:  taskstable.InitialModel(ctx, log),
		widgetTaskSidebar: taskssidebar.InitialModel(ctx, log),
		spinner:           s,
		showSpinner:       true,
		log:               log,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {

		case "esc":
			switch m.state {
			case TasksStateTaskSidebar:
				m.state = TasksStateTasksTable
				m.widgetTasksTable.Focused = true
				m.widgetTaskSidebar.Focused = false
				m.widgetViewsTabs.Focused = false

			case TasksStateViewsTabs:
				m.state = TasksStateTasksTable
				m.widgetTasksTable.Focused = true
				m.widgetTaskSidebar.Focused = false
				m.widgetViewsTabs.Focused = false

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
					m.widgetTasksTable.Focused = false
					m.widgetTaskSidebar.Focused = false
					m.widgetViewsTabs.Focused = true
				}

			case "j", "down", "k", "up":
				if m.state == TasksStateViewsTabs {
					m.state = TasksStateTasksTable
					m.widgetTasksTable.Focused = true
					m.widgetTaskSidebar.Focused = false
					m.widgetViewsTabs.Focused = false
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

	case tea.WindowSizeMsg:
		m.log.Info("Received: tea.WindowSizeMsg",
			"width", msg.Width,
			"height", msg.Height)

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
			m.widgetTasksTable.Focused = false
			m.widgetTaskSidebar.Focused = true
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

func (m Model) View() string {
	if m.showSpinner {
		return lipgloss.Place(
			m.ctx.WindowSize.Width, m.ctx.WindowSize.Height,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("%s Loading tasks...", m.spinner.View()),
		)
	}

	tableAndSidebar := ""
	if m.widgetTasksTable.Hidden {
		tableAndSidebar = lipgloss.Place(
			int(0.9*float32(m.ctx.WindowSize.Width)),
			int(0.7*float32(m.ctx.WindowSize.Height)),
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
		BorderRight(true).
		BorderBottom(true).
		BorderTop(true).
		BorderLeft(true).
		Render(lipgloss.JoinVertical(
			lipgloss.Top,
			m.widgetViewsTabs.View(),
			tableAndSidebar,
		))
}

func (m Model) Init() tea.Cmd {
	m.log.Infof("Initializing...")
	return tea.Batch(
		m.widgetViewsTabs.Init(),
		m.widgetTasksTable.Init(),
		m.widgetTaskSidebar.Init(),
		m.spinner.Tick,
	)
}
