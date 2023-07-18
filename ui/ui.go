package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/components/spaces"
	"github.com/prgrs/clickup/ui/components/tabs"
	"github.com/prgrs/clickup/ui/components/tickets"
	"github.com/prgrs/clickup/ui/context"
)

type sessionState uint

const (
	sessionSpacesView sessionState = iota
	sessionTasksView
)

type Model struct {
	ctx     *context.UserContext
	tabs    tabs.Model
	tickets tickets.Model
	spaces  spaces.Model
	state   sessionState
}

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		ctx:     ctx,
		tabs:    tabs.InitialModel(ctx),
		spaces:  spaces.InitialModel(ctx),
		tickets: tickets.InitialModel(ctx),
		state:   sessionTasksView,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == sessionSpacesView {
				m.state = sessionTasksView
			} else {
				m.state = sessionSpacesView
			}
		}

		switch m.state {
		case sessionSpacesView:
			m.spaces, cmd = m.spaces.Update(msg)
		case sessionTasksView:
			m.tickets, cmd = m.tickets.Update(msg)
		}
	default:
		m.spaces, cmd = m.spaces.Update(msg)
		m.tickets, cmd = m.tickets.Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.spaces.View(),
		m.tickets.View(),
	)
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spaces.Init, m.tickets.Init)
}

type ticketsMsg []clickup.Task
