package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
	"github.com/prgrs/clickup/ui/views/compact"
	"github.com/prgrs/clickup/ui/widgets/help"
)

type Model struct {
	state common.ViewId

	ctx         *context.UserContext
	viewCompact common.View
	log         *log.Logger

	// viewWorkspaces common.View
	// viewSpaces     common.View
	// viewFolders    common.View
	// viewLists      common.View
	//
	//
	//
	//
	// viewTasks      common.View

	dialogHelp help.Model

	KeyMap KeyMap
}

type KeyMap struct {
	// GoToViewWorkspaces key.Binding
	// GoToViewSpaces     key.Binding
	// GoToViewFolders    key.Binding
	// GoToViewLists      key.Binding
	// GoToViewTasks      key.Binding
	Refresh   key.Binding
	ForceQuit key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		// GoToViewWorkspaces: key.NewBinding(
		// 	key.WithKeys("1"),
		// 	key.WithHelp("1", "go to workspaces"),
		// ),
		// GoToViewSpaces: key.NewBinding(
		// 	key.WithKeys("2"),
		// 	key.WithHelp("2", "go to spaces"),
		// ),
		// GoToViewFolders: key.NewBinding(
		// 	key.WithKeys("3"),
		// 	key.WithHelp("3", "go to folders"),
		// ),
		// GoToViewLists: key.NewBinding(
		// 	key.WithKeys("4"),
		// 	key.WithHelp("4", "go to lists"),
		// ),
		// GoToViewTasks: key.NewBinding(
		// 	key.WithKeys("5"),
		// 	key.WithHelp("5", "go to tasks"),
		// ),
		Refresh: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "go to refresh"),
		),
		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("ctrl+c/q", "quit"),
		),
	}
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	log := logger.WithPrefix("UI")

	// var state common.ViewId = workspaces.ViewId
	// if ctx.Config.DefaultWorkspace != "" {
	// 	state = spaces.ViewId
	// }

	return Model{
		ctx: ctx,
		// state: state,
		log: log,

		// viewWorkspaces: workspaces.InitialModel(ctx, log),
		// viewSpaces:     spaces.InitialModel(ctx, log),
		// viewTasks:      tasks.InitialModel(ctx, log),
		// viewLists:      lists.InitialModel(ctx, log),
		// viewFolders:    folders.InitialModel(ctx, log),

		viewCompact: compact.InitialModel(ctx, log),

		dialogHelp: help.InitialModel(ctx, log),
		KeyMap:     DefaultKeyMap(),
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
		switch {
		case key.Matches(msg, m.KeyMap.ForceQuit):
			return m, tea.Quit

		// case key.Matches(msg, m.KeyMap.GoToViewWorkspaces):
		// 	if m.viewWorkspaces.Ready() {
		// 		m.state = m.viewWorkspaces.GetViewId()
		// 	}
		// case key.Matches(msg, m.KeyMap.GoToViewSpaces):
		// 	if m.viewSpaces.Ready() {
		// 		m.state = m.viewSpaces.GetViewId()
		// 	}
		// case key.Matches(msg, m.KeyMap.GoToViewFolders):
		// 	if m.viewFolders.Ready() {
		// 		m.state = m.viewFolders.GetViewId()
		// 	}
		// case key.Matches(msg, m.KeyMap.GoToViewLists):
		// 	if m.viewLists.Ready() {
		// 		m.state = m.viewLists.GetViewId()
		// 	}
		// case key.Matches(msg, m.KeyMap.GoToViewTasks):
		// 	if m.viewTasks.Ready() {
		// 		m.state = m.viewTasks.GetViewId()
		// 	}
		case key.Matches(msg, m.KeyMap.Refresh):
			m.log.Info("Refreshing...")
			if err := m.ctx.Api.InvalidateCache(); err != nil {
				m.log.Error("Failed to invalidate cache", "error", err)
			}
			m.log.Debug("Cache invalidated")

			// default:
			// switch m.state {
			// case m.viewSpaces.GetViewId():
			// 	m.viewSpaces, cmd = m.viewSpaces.Update(msg)
			// case m.viewFolders.GetViewId():
			// 	m.viewFolders, cmd = m.viewFolders.Update(msg)
			// case m.viewLists.GetViewId():
			// 	m.viewLists, cmd = m.viewLists.Update(msg)
			// case m.viewTasks.GetViewId():
			// 	m.viewTasks, cmd = m.viewTasks.Update(msg)
			// case m.viewWorkspaces.GetViewId():
			// 	m.viewWorkspaces, cmd = m.viewWorkspaces.Update(msg)
			// }
			//
			// cmds = append(cmds, cmd)
			//
			// m.dialogHelp, cmd = m.dialogHelp.Update(msg)
			// cmds = append(cmds, cmd)
			//
			// return m, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		m.log.Debug("Received: tea.WindowSizeMsg",
			"width", msg.Width,
			"height", msg.Height)
		m.ctx.WindowSize.Width = msg.Width
		m.ctx.WindowSize.Height = msg.Height
		m.ctx.WindowSize.Height = msg.Height
		m.ctx.WindowSize.Height = msg.Height
		m.ctx.WindowSize.Height = msg.Height

		return m, nil

		// case common.WorkspaceChangeMsg:
		// 	workspace := string(msg)
		// 	m.log.Info("Received: WorkspaceChangeMsg", "workspace", workspace)
		// 	m.state = m.viewSpaces.GetViewId()
		//
		// case common.SpaceChangeMsg:
		// 	m.log.Info("Received: SpaceChangeMsg", "space", string(msg))
		// 	m.state = m.viewFolders.GetViewId()
		//
		// case common.FolderChangeMsg:
		// 	m.log.Info("Received: FolderChangeMsg", "folder", string(msg))
		// 	m.state = m.viewLists.GetViewId()
		//
		// case common.ListChangeMsg:
		// 	id := string(msg)
		// 	m.log.Info("Received: ListChangeMsg", "list", id)
		// 	m.state = m.viewTasks.GetViewId()

		// case common.BackToPreviousViewMsg:
		// 	m.log.Info("Received: BackToPreviousViewMsg")
		// 	switch m.state {
		// 	case m.viewSpaces.GetViewId():
		// 		m.state = m.viewWorkspaces.GetViewId()
		// 	case m.viewFolders.GetViewId():
		// 		m.state = m.viewSpaces.GetViewId()
		// 	case m.viewLists.GetViewId():
		// 		m.state = m.viewFolders.GetViewId()
		// 	case m.viewTasks.GetViewId():
		// 		m.state = m.viewLists.GetViewId()
		// }

		// m.dialogHelp.ShowHelp = false
	}

	// m.viewWorkspaces, cmd = m.viewWorkspaces.Update(msg)
	// cmds = append(cmds, cmd)
	//
	// m.viewSpaces, cmd = m.viewSpaces.Update(msg)
	// cmds = append(cmds, cmd)
	//
	// m.viewFolders, cmd = m.viewFolders.Update(msg)
	// cmds = append(cmds, cmd)
	//
	// m.viewLists, cmd = m.viewLists.Update(msg)
	// cmds = append(cmds, cmd)
	//
	// m.viewTasks, cmd = m.viewTasks.Update(msg)
	// cmds = append(cmds, cmd)

	m.viewCompact, cmd = m.viewCompact.Update(msg)
	cmds = append(cmds, cmd)

	m.dialogHelp, cmd = m.dialogHelp.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	// m.log.Info("Rendering...")
	var viewToRender common.View

	// switch m.state {
	// case m.viewWorkspaces.GetViewId():
	// 	viewToRender = m.viewWorkspaces
	// case m.viewSpaces.GetViewId():
	// 	viewToRender = m.viewSpaces
	// case m.viewFolders.GetViewId():
	// 	viewToRender = m.viewFolders
	// case m.viewLists.GetViewId():
	// 	viewToRender = m.viewLists
	// case m.viewTasks.GetViewId():
	// 	viewToRender = m.viewTasks
	// default:
	// 	panic("Unknown view")
	// }

	viewToRender = m.viewCompact

	viewKm := viewToRender.KeyMap()
	km := common.NewKeyMap(
		func() [][]key.Binding {
			return append(viewKm.FullHelp(), [][]key.Binding{
				// {
				// 	m.KeyMap.GoToViewWorkspaces,
				// 	m.KeyMap.GoToViewSpaces,
				// 	m.KeyMap.GoToViewFolders,
				// 	m.KeyMap.GoToViewLists,
				// 	m.KeyMap.GoToViewTasks,
				// },
				{
					m.KeyMap.Refresh,
				},
			}...)
		},
		viewKm.ShortHelp,
	)

	footer := m.dialogHelp.View(km)
	footerHeight := lipgloss.Height(footer)

	physicalHeight := m.ctx.WindowSize.Height
	physicalWidth := m.ctx.WindowSize.Width

	viewHeight := physicalHeight - footerHeight
	viewToRender = viewToRender.SetSize(common.Size{
		Width:  physicalWidth,
		Height: viewHeight - m.ctx.WindowSize.MetaHeight,
	})

	dividerHeight := physicalHeight - viewHeight - footerHeight

	if dividerHeight < 0 {
		dividerHeight = 0
		m.log.Info("dividerHeight", "dividerHeight", dividerHeight)
	}

	divider := strings.Repeat("\n", dividerHeight)

	m.ctx.WindowSize.MetaHeight = lipgloss.Height(divider) + footerHeight

	return lipgloss.JoinVertical(
		lipgloss.Left,
		viewToRender.View(),
		divider,
		footer,
	)
}

func (m Model) Init() tea.Cmd {
	m.log.Info("Initializing...")
	return tea.Batch(
		// m.viewWorkspaces.Init(),
		// m.viewSpaces.Init(),
		// m.viewFolders.Init(),
		// m.viewLists.Init(),
		// m.viewTasks.Init(),
		m.viewCompact.Init(),
		m.dialogHelp.Init(),
	)
}
