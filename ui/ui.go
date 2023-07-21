package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/components/tabs"
	"github.com/prgrs/clickup/ui/components/tickets"
	"github.com/prgrs/clickup/ui/components/views"
	"github.com/prgrs/clickup/ui/context"
	"github.com/prgrs/clickup/ui/views/spaces"
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
	sessionViewsView
)

type Model struct {
	ctx   *context.UserContext
	state sessionState

	tabs    tabs.Model
	views   views.Model
	tickets tickets.Model

	viewSpaces spaces.Model
}

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		ctx:  ctx,
		tabs: tabs.InitialModel(ctx),

		views:   views.InitialModel(ctx),
		tickets: tickets.InitialModel(ctx),
		state:   sessionTasksView,

		viewSpaces: spaces.InitialModel(ctx),
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

		case "tab":
			if m.state == sessionTasksView {
				return m, ChangeViewCmd(sessionViewsView)
			}

			if m.state == sessionViewsView {
				return m, ChangeViewCmd(sessionTasksView)
			}

		default:
			switch m.state {
			case sessionSpacesView:
				m.viewSpaces, cmd = m.viewSpaces.Update(msg)
				return m, cmd

			case sessionTasksView:
				m.tickets, cmd = m.tickets.Update(msg)
				return m, cmd

			case sessionViewsView:
				m.views, cmd = m.views.Update(msg)
				return m, cmd

			default:
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.ctx.Logger.Info("UI received tea.WindowSizeMsg")
		m.ctx.WindowSize.Width = msg.Width
		m.ctx.WindowSize.Height = msg.Height
		// return here is disable on purpose to allow the msg
		// to be passed to the other components

	case ChangeViewMsg: // maybe ChangeScreenMsg
		m.ctx.Logger.Infof("UI received ChangeViewMsg: %d", msg)

		switch sessionState(msg) {
		case sessionSpacesView:
			m.state = sessionSpacesView
			m.viewSpaces, cmd = m.viewSpaces.Update(common.FocusMsg(true))
			return m, cmd

		case sessionTasksView:
			m.state = sessionTasksView
			m.tickets, cmd = m.tickets.Update(common.FocusMsg(true))
			return m, cmd

		case sessionViewsView:
			m.state = sessionViewsView
			m.views, cmd = m.views.Update(common.FocusMsg(true))
			return m, cmd
		}

	case spaces.HideSpaceViewMsg:
		m.ctx.Logger.Info("UI received HideSpaceViewMsg")
		return m, ChangeViewCmd(sessionTasksView)

	case spaces.SpaceChangeMsg:
		m.ctx.Logger.Infof("UI received SpaceChangeMsg: %s", string(msg))
		return m, tea.Batch(
			views.SpaceChangedCmd(string(msg)),
			ChangeViewCmd(sessionTasksView))

	case views.ViewChangedMsg:
		m.ctx.Logger.Infof("UI received ViewChangedMsg: %s", string(msg))

		return m, tea.Batch(
			tickets.ViewChangedCmd(string(msg)),
		)

	case views.FetchViewsMsg:
		m.ctx.Logger.Infof(
			"UI received FetchedViewsMsg: %s",
			strings.Join(msg, ", "))

		var cmds []tea.Cmd
		for _, viewID := range msg {
			cmds = append(cmds, tickets.FetchTasksForViewCmd(viewID))
		}

		return m, tea.Batch(cmds...)
	}

	m.viewSpaces, cmd = m.viewSpaces.Update(msg)
	cmds = append(cmds, cmd)

	m.tickets, cmd = m.tickets.Update(msg)
	cmds = append(cmds, cmd)

	m.views, cmd = m.views.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.state {
	case sessionSpacesView:
		return m.viewSpaces.View()
	case sessionTasksView:
		return lipgloss.JoinVertical(
			lipgloss.Top,
			m.views.View(),
			m.tickets.View(),
		)

	default:
		return lipgloss.JoinVertical(
			lipgloss.Top,
			m.views.View(),
			m.tickets.View(),
		)
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.viewSpaces.Init,
		m.views.Init,
		m.tickets.Init,
	)
}
