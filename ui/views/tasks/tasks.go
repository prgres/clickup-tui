package tasks

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
	"github.com/prgrs/clickup/ui/widgets/tasksidebar"
	"github.com/prgrs/clickup/ui/widgets/tasktable"
	"github.com/prgrs/clickup/ui/widgets/viewtabs"
)

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
	widgetViewsTabs   viewtabs.Model
	widgetTasksTable  tasktable.Model
	widgetTaskSidebar tasksidebar.Model
	spinner           spinner.Model
	showSpinner       bool
}

func InitialModel(ctx *context.UserContext) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	return Model{
		ViewId:            "viewTasks",
		ctx:               ctx,
		state:             TasksStateTasksTable,
		widgetViewsTabs:   viewtabs.InitialModel(ctx),
		widgetTasksTable:  tasktable.InitialModel(ctx),
		widgetTaskSidebar: tasksidebar.InitialModel(ctx),
		spinner:           s,
		showSpinner:       true,
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
				m.ctx.Logger.Info("ViewTasks: Go to previous view")
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
				return m, tea.Batch(cmds...)

			case TasksStateTaskSidebar:
				m.widgetTaskSidebar, cmd = m.widgetTaskSidebar.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)

			case TasksStateViewsTabs:
				m.widgetViewsTabs, cmd = m.widgetViewsTabs.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)

			default:
				m.ctx.Logger.Infof("ViewTasks received unhandled keypress: %s", keypress)
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.ctx.Logger.Info("ViewTasks receive tea.WindowSizeMsg")

	case viewtabs.TabChangedMsg:
		tab := viewtabs.Tab(msg)
		m.ctx.Logger.Infof("ViewTasks received TabChangedMsg: name=%s id=%s", tab.Name, tab.Id)
		m.showSpinner = true

		cmds = append(cmds, tasktable.TabChangedCmd(tab))

		m.widgetViewsTabs, cmd = m.widgetViewsTabs.Update(msg)
		cmds = append(cmds, cmd)

	case tasktable.TasksListReadyMsg:
		m.ctx.Logger.Info("ViewTasks received TasksListReady")
		m.showSpinner = false
		cmds = append(cmds,
			m.spinner.Tick,
		)

	case spinner.TickMsg:
		// m.ctx.Logger.Info("ViewTask receive spinner.TickMsg")
		if m.showSpinner {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case tasktable.TaskSelectedMsg:
		id := string(msg)
		m.ctx.Logger.Infof("ViewTask receive tasktable.TaskSelectedMsg: %s", id)
		m.state = TasksStateTaskSidebar
		m.widgetTasksTable.Focused = false
		m.widgetTaskSidebar.Focused = true
		cmds = append(cmds, tasksidebar.TaskSelectedCmd(id))
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
	m.ctx.Logger.Infof("Initializing view: Tasks")
	return tea.Batch(
		m.widgetViewsTabs.Init(),
		m.widgetTasksTable.Init(),
		m.widgetTaskSidebar.Init(),
		m.spinner.Tick,
	)
}
