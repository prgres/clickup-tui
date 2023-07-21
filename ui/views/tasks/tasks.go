package tasks

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/ui/components/taskstable"
	"github.com/prgrs/clickup/ui/components/views"
	"github.com/prgrs/clickup/ui/context"
)

type TasksState uint

const (
	TasksStateLoading TasksState = iota
	TasksStateTasksTable
	TasksStateViewsTabs
)

type Model struct {
	ctx   *context.UserContext
	state TasksState

	componentViewsTabs  views.Model
	componentTasksTable taskstable.Model

	spinner     spinner.Model
	showSpinner bool
}

func InitialModel(ctx *context.UserContext) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse
	// spinner.Dot,
	// spinner.Line,
	// spinner.Pulse,
	// spinner.Points,
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		ctx:   ctx,
		state: TasksStateTasksTable,

		componentViewsTabs:  views.InitialModel(ctx),
		componentTasksTable: taskstable.InitialModel(ctx),

		spinner:     s,
		showSpinner: true,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.showSpinner {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "tab":
			switch m.state {
			case TasksStateTasksTable:
				m.state = TasksStateViewsTabs
				return m, tea.Batch(cmds...)

			case TasksStateViewsTabs:
				m.state = TasksStateTasksTable
				return m, tea.Batch(cmds...)
			}
		default:
			switch m.state {
			case TasksStateTasksTable:
				m.componentTasksTable, cmd = m.componentTasksTable.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)

			case TasksStateViewsTabs:
				m.componentViewsTabs, cmd = m.componentViewsTabs.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
		}

	case views.FetchViewsMsg:
		m.ctx.Logger.Infof("ViewTasks received FetchViewsMsg: %s",
			strings.Join(msg, ", "))

		var cmds []tea.Cmd
		for _, viewID := range msg {
			cmds = append(cmds, taskstable.FetchTasksForViewCmd(viewID))
		}

		return m, tea.Batch(cmds...)

	case views.ViewChangedMsg:
		m.ctx.Logger.Info("ViewTasks received ViewChangedMsg")
		m.showSpinner = true
		cmd = taskstable.ViewChangedCmd(string(msg))
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case taskstable.TasksListReady:
		m.ctx.Logger.Info("ViewTasks received TasksListReady")
		m.showSpinner = false
		return m, tea.Batch(cmds...)

	case SpaceChangedMsg:
		m.ctx.Logger.Info("ViewTasks received SpaceChangedMsg")
		cmd = views.SpaceChangedCmd(string(msg))
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case spinner.TickMsg:
		m.ctx.Logger.Info("TaskView receive spinner.TickMsg")
		if !m.showSpinner {
			return m, tea.Batch(cmds...)
		}
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	m.componentViewsTabs, cmd = m.componentViewsTabs.Update(msg)
	cmds = append(cmds, cmd)

	m.componentTasksTable, cmd = m.componentTasksTable.Update(msg)
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

	return lipgloss.JoinVertical(
		lipgloss.Top,
		m.componentViewsTabs.View(),
		m.componentTasksTable.View(),
	)
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Infof("Initializing view: Tasks")
	return tea.Batch(
		m.componentViewsTabs.Init(),
		m.componentTasksTable.Init(),
		m.spinner.Tick,
	)
}
