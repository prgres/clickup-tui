package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
	"github.com/prgrs/clickup/ui/views/folders"
	"github.com/prgrs/clickup/ui/views/lists"
	"github.com/prgrs/clickup/ui/views/spaces"
	"github.com/prgrs/clickup/ui/views/tasks"
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
	sessionFoldersView
	sessionListsView
	sessionTasksView
)

type Model struct {
	ctx   *context.UserContext
	state sessionState

	viewSpaces  spaces.Model
	viewTasks   tasks.Model
	viewLists   lists.Model
	viewFolders folders.Model
}

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		ctx:   ctx,
		state: sessionTasksView,

		viewSpaces:  spaces.InitialModel(ctx),
		viewTasks:   tasks.InitialModel(ctx),
		viewLists:   lists.InitialModel(ctx),
		viewFolders: folders.InitialModel(ctx),
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

		case "2":
			return m, ChangeViewCmd(sessionFoldersView)

		case "3":
			return m, ChangeViewCmd(sessionListsView)

		case "4":
			return m, ChangeViewCmd(sessionTasksView)

		case "esc":
			switch m.state {
			case sessionFoldersView:
				return m, ChangeViewCmd(sessionSpacesView)

			case sessionListsView:
				return m, ChangeViewCmd(sessionFoldersView)
			}

		default:
			switch m.state {
			case sessionSpacesView:
				m.viewSpaces, cmd = m.viewSpaces.Update(msg)
				return m, cmd

			case sessionTasksView:
				m.viewTasks, cmd = m.viewTasks.Update(msg)
				return m, cmd

			case sessionFoldersView:
				m.viewFolders, cmd = m.viewFolders.Update(msg)
				return m, cmd

			case sessionListsView:
				m.viewLists, cmd = m.viewLists.Update(msg)
				return m, cmd

			default:
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.ctx.Logger.Info("UI received tea.WindowSizeMsg")
		m.ctx.WindowSize.Width = msg.Width
		m.ctx.WindowSize.Height = msg.Height
		cmds = append(cmds, common.WindowSizeCmd(msg))

	case ChangeViewMsg: // maybe ChangeScreenMsg
		m.ctx.Logger.Infof("UI received ChangeViewMsg: %d", msg)

		switch sessionState(msg) {
		case sessionSpacesView:
			m.state = sessionSpacesView
			m.viewSpaces, cmd = m.viewSpaces.Update(common.FocusMsg(true))
			return m, cmd

		case sessionFoldersView:
			m.state = sessionFoldersView
			m.viewFolders, cmd = m.viewFolders.Update(common.FocusMsg(true))
			return m, cmd

		case sessionListsView:
			m.state = sessionListsView
			m.viewLists, cmd = m.viewLists.Update(common.FocusMsg(true))
			return m, cmd

		case sessionTasksView:
			m.state = sessionTasksView
			m.viewSpaces, cmd = m.viewSpaces.Update(common.FocusMsg(true))
			return m, cmd
		}

	case spaces.HideSpaceViewMsg:
		m.ctx.Logger.Info("UI received HideSpaceViewMsg")
		return m, ChangeViewCmd(sessionTasksView)

	case common.SpaceChangeMsg:
		m.ctx.Logger.Infof("UI received SpaceChangeMsg: %s", string(msg))
		cmds = append(cmds, ChangeViewCmd(sessionFoldersView))

	case common.FolderChangeMsg:
		m.ctx.Logger.Infof("UI received FolderChangeMsg: %s", string(msg))
		cmds = append(cmds, ChangeViewCmd(sessionListsView))

	case common.ListChangeMsg:
		m.ctx.Logger.Infof("UI received ListChangeMsg: %s", string(msg))
		cmds = append(cmds, ChangeViewCmd(sessionTasksView))
	}

	m.viewSpaces, cmd = m.viewSpaces.Update(msg)
	cmds = append(cmds, cmd)

	m.viewTasks, cmd = m.viewTasks.Update(msg)
	cmds = append(cmds, cmd)

	m.viewFolders, cmd = m.viewFolders.Update(msg)
	cmds = append(cmds, cmd)

	m.viewLists, cmd = m.viewLists.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.state {
	case sessionSpacesView:
		return m.viewSpaces.View()
	case sessionTasksView:
		return m.viewTasks.View()
	case sessionListsView:
		return m.viewLists.View()
	case sessionFoldersView:
		return m.viewFolders.View()
	default:
		return m.viewTasks.View()
	}
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Info("Initializing UI")
	return tea.Batch(
		m.viewSpaces.Init(),
		m.viewTasks.Init(),
		m.viewLists.Init(),
		m.viewFolders.Init(),
	)
}
