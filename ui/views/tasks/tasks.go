package tasks

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	ctx   *context.UserContext
	state TasksState

	widgetViewsTabs   viewtabs.Model
	widgetTasksTable  tasktable.Model
	widgetTaskSidebar tasksidebar.Model

	spinner     spinner.Model
	showSpinner bool
}

func InitialModel(ctx *context.UserContext) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	return Model{
		ctx:   ctx,
		state: TasksStateTasksTable,

		widgetViewsTabs:   viewtabs.InitialModel(ctx),
		widgetTasksTable:  tasktable.InitialModel(ctx),
		widgetTaskSidebar: tasksidebar.InitialModel(ctx),

		spinner:     s,
		showSpinner: false,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// if m.showSpinner {
	// 	m.spinner, cmd = m.spinner.Update(msg)
	// 	cmds = append(cmds, cmd)
	// }

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "tab":
			switch m.state {
			case TasksStateTasksTable:
				m.state = TasksStateViewsTabs
				m.widgetTasksTable.Focused = false
				m.widgetViewsTabs.Focused = true
				return m, tea.Batch(cmds...)

			case TasksStateViewsTabs:
				m.state = TasksStateTasksTable
				m.widgetTasksTable.Focused = true
				m.widgetViewsTabs.Focused = false
				return m, tea.Batch(cmds...)
			}

		case "esc":
			m.state = TasksStateTasksTable
			m.widgetTasksTable.Focused = true
			m.widgetTaskSidebar.Focused = false
			m.widgetViewsTabs.Focused = false
			return m, tea.Batch(cmds...)

		default:
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
			}
		}

	case viewtabs.FetchViewsMsg:
		m.ctx.Logger.Infof("ViewTasks received FetchViewsMsg: %s",
			strings.Join(msg, ", "))

		var cmds []tea.Cmd
		for _, viewID := range msg {
			cmds = append(cmds, tasktable.FetchTasksForViewCmd(viewID))
		}

	case tea.WindowSizeMsg:
		m.ctx.Logger.Info("ViewTasks receive tea.WindowSizeMsg")

	case viewtabs.ViewChangedMsg:
		m.ctx.Logger.Info("ViewTasks received ViewChangedMsg")
		m.showSpinner = true

		cmd = tasktable.ViewChangedCmd(string(msg))
		cmds = append(cmds, cmd)

		m.widgetViewsTabs, cmd = m.widgetViewsTabs.Update(msg)
		cmds = append(cmds, cmd)

	case tasktable.TasksListReadyMsg:
		id := string(msg)
		m.ctx.Logger.Infof("ViewTasks received TasksListReady: %s", id)
		m.showSpinner = false
		cmds = append(cmds, tasksidebar.TaskSelectedCmd(id))

	case spinner.TickMsg:
		m.ctx.Logger.Info("ViewTask receive spinner.TickMsg")
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

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderRight(true).
		BorderBottom(true).
		BorderTop(true).
		BorderLeft(true).
		Render(lipgloss.JoinVertical(
			lipgloss.Top,
			m.widgetViewsTabs.View(),
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				m.widgetTasksTable.View(),
				m.widgetTaskSidebar.View(),
			),
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
