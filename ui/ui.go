package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/ui/common"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
	"github.com/prgrs/clickup/ui/context"
	"github.com/prgrs/clickup/ui/views/folders"
	"github.com/prgrs/clickup/ui/views/lists"
	"github.com/prgrs/clickup/ui/views/spaces"
	"github.com/prgrs/clickup/ui/views/tasks"
	"github.com/prgrs/clickup/ui/views/workspaces"
)

type Model struct {
	ctx   *context.UserContext
	state common.ViewId
	log   *log.Logger

	viewWorkspaces workspaces.Model
	viewSpaces     spaces.Model
	viewTasks      tasks.Model
	viewLists      lists.Model
	viewFolders    folders.Model
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	log := logger.WithPrefix("UI")

	return Model{
		ctx:   ctx,
		state: spaces.ViewId,
		log:   log,

		viewWorkspaces: workspaces.InitialModel(ctx, log),
		viewSpaces:     spaces.InitialModel(ctx, log),
		viewTasks:      tasks.InitialModel(ctx, log),
		viewLists:      lists.InitialModel(ctx, log),
		viewFolders:    folders.InitialModel(ctx, log),
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case common.ErrMsg:
		m.log.Fatal(msg.Error())
		return m, tea.Quit

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "1":
			if m.viewWorkspaces.Ready() {
				m.state = m.viewWorkspaces.ViewId
			}
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

			case m.viewWorkspaces.ViewId:
				m.viewWorkspaces, cmd = m.viewWorkspaces.Update(msg)
				return m, cmd

			default:
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		m.log.Info("Received: tea.WindowSizeMsg",
			"width", msg.Width,
			"height", msg.Height)
		m.ctx.WindowSize.Width = msg.Width
		m.ctx.WindowSize.Height = msg.Height

	case common.WorkspaceChangeMsg:
		workspace := string(msg)
		m.log.Info("Received: WorkspaceChangeMsg", "workspace", workspace)
		m.state = m.viewSpaces.ViewId

	case common.SpaceChangeMsg:
		m.log.Info("Received: SpaceChangeMsg", "space", string(msg))
		m.state = m.viewFolders.ViewId

	case common.FolderChangeMsg:
		m.log.Info("Received: FolderChangeMsg", "folder", string(msg))
		m.state = m.viewLists.ViewId

	case common.ListChangeMsg:
		m.log.Info("Received: ListChangeMsg", "list", listitem.Item(msg).Description())
		m.state = m.viewTasks.ViewId

	case common.BackToPreviousViewMsg:
		m.log.Info("Received: BackToPreviousViewMsg")
		switch m.state {
		case m.viewFolders.ViewId:
			m.state = m.viewSpaces.ViewId
		case m.viewLists.ViewId:
			m.state = m.viewFolders.ViewId
		case m.viewTasks.ViewId:
			m.state = m.viewLists.ViewId
		}
	}

	m.viewWorkspaces, cmd = m.viewWorkspaces.Update(msg)
	cmds = append(cmds, cmd)

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
	case m.viewWorkspaces.ViewId:
		return m.viewWorkspaces.View()
	default:
		return m.viewSpaces.View()
	}
}

func (m Model) Init() tea.Cmd {
	m.log.Info("Initializing...")
	return tea.Batch(
		m.viewWorkspaces.Init(),
		m.viewSpaces.Init(),
		m.viewTasks.Init(),
		m.viewLists.Init(),
		m.viewFolders.Init(),
	)
}
