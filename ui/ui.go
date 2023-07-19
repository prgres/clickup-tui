package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
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
	case common.ErrMsg:
		m.ctx.Logger.Fatal(msg.Error())
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1":
			m.state = sessionSpacesView
		case "esc":
			m.state = sessionTasksView
		}

		switch m.state {
		case sessionSpacesView:
			m.spaces, cmd = m.spaces.Update(msg)
		case sessionTasksView:
			m.tickets, cmd = m.tickets.Update(msg)
		}
		cmds = append(cmds, cmd)
	case spaces.HideSpaceMsg:
		m.ctx.Logger.Info("hide space")
		m.state = sessionTasksView
		m.tickets.SelectedSpace = m.spaces.SelectedSpace
		m.ctx.Logger.Info(m.tickets.SelectedSpace, m.tickets.PrevSelectedSpace, m.spaces.SelectedSpace)

		m.tickets, cmd = m.tickets.Update(msg)
		cmds = append(cmds, cmd)
	default:
		m.spaces, cmd = m.spaces.Update(msg)
		cmds = append(cmds, cmd)
		m.tickets, cmd = m.tickets.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.state == sessionSpacesView {
		return m.spaces.View()
	} else {
		return m.tickets.View()
	}

	// return lipgloss.JoinHorizontal(
	// 	lipgloss.Left,
	// 	m.spaces.View(),
	// 	m.tickets.View(),
	// )
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spaces.Init, m.tickets.Init)
}

type ticketsMsg []clickup.Task
