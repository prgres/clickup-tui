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

type Model struct {
	ctx         *context.UserContext
	state       common.ViewId
	viewSpaces  spaces.Model
	viewTasks   tasks.Model
	viewLists   lists.Model
	viewFolders folders.Model
}

func InitialModel(ctx *context.UserContext) Model {
	viewSpaces := spaces.InitialModel(ctx)
	viewTasks := tasks.InitialModel(ctx)
	viewLists := lists.InitialModel(ctx)
	viewFolders := folders.InitialModel(ctx)

	return Model{
		ctx:         ctx,
		state:       viewSpaces.ViewId,
		viewSpaces:  viewSpaces,
		viewTasks:   viewTasks,
		viewLists:   viewLists,
		viewFolders: viewFolders,
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

		case "2":
			if m.viewSpaces.Ready() {
				m.state = m.viewSpaces.ViewId
			}
		case "3":
			if m.viewFolders.Ready() {
				m.state = m.viewFolders.ViewId
			}
		case "4":
			if m.viewLists.Ready() {
				m.state = m.viewLists.ViewId
			}
		case "5":
			if m.viewTasks.Ready() {
				m.state = m.viewTasks.ViewId
			}

		default:
			switch m.state {
			case m.viewSpaces.ViewId:
				m.viewSpaces, cmd = m.viewSpaces.Update(msg)
				return m, cmd

			case m.viewFolders.ViewId:
				m.viewFolders, cmd = m.viewFolders.Update(msg)
				return m, cmd

			case m.viewLists.ViewId:
				m.viewLists, cmd = m.viewLists.Update(msg)
				return m, cmd

			case m.viewTasks.ViewId:
				m.viewTasks, cmd = m.viewTasks.Update(msg)
				return m, cmd

			default:
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.ctx.Logger.Infof("UI received tea.WindowSizeMsg. Width: %d Height %d", msg.Width, msg.Height)
		m.ctx.WindowSize.Width = msg.Width
		m.ctx.WindowSize.Height = msg.Height

	case common.SpaceChangeMsg:
		m.ctx.Logger.Infof("UI received SpaceChangeMsg: %s", string(msg))
		m.state = m.viewFolders.ViewId

	case common.FolderChangeMsg:
		m.ctx.Logger.Infof("UI received FolderChangeMsg: %s", string(msg))
		m.state = m.viewLists.ViewId

	case common.ListChangeMsg:
		m.ctx.Logger.Infof("UI received ListChangeMsg: %s", string(msg))
		m.state = m.viewTasks.ViewId

	case common.BackToPreviousViewMsg:
		m.ctx.Logger.Infof("UI received BackToPreviousViewMsg")
		switch m.state {
		case m.viewFolders.ViewId:
			m.state = m.viewSpaces.ViewId
		case m.viewLists.ViewId:
			m.state = m.viewFolders.ViewId
		case m.viewTasks.ViewId:
			m.state = m.viewLists.ViewId
		}
	}

	m.viewSpaces, cmd = m.viewSpaces.Update(msg)
	cmds = append(cmds, cmd)

	m.viewFolders, cmd = m.viewFolders.Update(msg)
	cmds = append(cmds, cmd)

	m.viewLists, cmd = m.viewLists.Update(msg)
	cmds = append(cmds, cmd)

	m.viewTasks, cmd = m.viewTasks.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.state {
	case m.viewSpaces.ViewId:
		return m.viewSpaces.View()
	case m.viewFolders.ViewId:
		return m.viewFolders.View()
	case m.viewLists.ViewId:
		return m.viewLists.View()
	case m.viewTasks.ViewId:
		return m.viewTasks.View()
	default:
		return m.viewSpaces.View()
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
