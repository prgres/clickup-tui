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

type ChangeViewMsg sessionState

func ChangeViewCmd(view sessionState) tea.Cmd {
	return func() tea.Msg {
		return ChangeViewMsg(view)
	}
}

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

	switch msg := msg.(type) {
	case common.ErrMsg:
		m.ctx.Logger.Fatal(msg.Error())
		return m, tea.Quit

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1":
			return m, ChangeViewCmd(sessionSpacesView)
		// case "esc":
		// 	return m, ChangeViewCmd(sessionTasksView)
		default:
			switch m.state {
			case sessionSpacesView:
				m.spaces, cmd = m.spaces.Update(msg)
				return m, cmd
			case sessionTasksView:
				m.tickets, cmd = m.tickets.Update(msg)
				return m, cmd
			}
		}
	case tea.WindowSizeMsg:
		m.ctx.Logger.Info("UI received tea.WindowSizeMsg")
		m.ctx.WindowSize.Width = msg.Width
		m.ctx.WindowSize.Height = msg.Height
		// return here is disable on purpose to allow the msg to be passed to the
		// other components

	case ChangeViewMsg:
		m.ctx.Logger.Info("UI received ChangeViewMsg")

		switch sessionState(msg) {
		case sessionSpacesView:
			m.state = sessionSpacesView
			m.spaces, cmd = m.spaces.Update(msg)
			return m, cmd

		case sessionTasksView:
			m.state = sessionTasksView
			m.tickets, cmd = m.tickets.Update(msg)
			return m, cmd
		}

	case spaces.HideSpaceViewMsg:
		m.ctx.Logger.Info("UI received HideSpaceViewMsg")
		return m, ChangeViewCmd(sessionTasksView)

	case spaces.SpaceChangeMsg:
	  m.ctx.Logger.Infof("UI received ChangeSpaceMsg: %d", string(msg))
		return m, tea.Batch(
			tickets.SpaceChangedCmd(string(msg)),
			ChangeViewCmd(sessionTasksView))
	}

	m.spaces, cmd = m.spaces.Update(msg)
	cmds = append(cmds, cmd)

	m.tickets, cmd = m.tickets.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.state {
	case sessionSpacesView:
		return m.spaces.View()
	case sessionTasksView:
		return m.tickets.View()
	default:
		return m.tickets.View()
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spaces.Init,
		m.tickets.Init,
	)
}

type ticketsMsg []clickup.Task
