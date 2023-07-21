package tasks

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/ui/components/tickets"
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
	componentTasksTable tickets.Model
}

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		ctx:                 ctx,
		componentViewsTabs:  views.InitialModel(ctx),
		componentTasksTable: tickets.InitialModel(ctx),
		state:               TasksStateTasksTable,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "tab":
			switch m.state {
			case TasksStateTasksTable:
			  m.state = TasksStateViewsTabs

			case TasksStateViewsTabs:
			  m.state = TasksStateTasksTable
			}
		default:
			switch m.state {
			case TasksStateTasksTable:
				m.componentTasksTable, cmd = m.componentTasksTable.Update(msg)
				return m, cmd

			case TasksStateViewsTabs:
				m.componentViewsTabs, cmd = m.componentViewsTabs.Update(msg)
				return m, cmd
			}
		}
	}

	m.componentViewsTabs, cmd = m.componentViewsTabs.Update(msg)
	cmds = append(cmds, cmd)

	m.componentTasksTable, cmd = m.componentTasksTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
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
	)
}
