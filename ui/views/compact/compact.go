package compact

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	viewstabs "github.com/prgrs/clickup/ui/components/views-tabs"
	"github.com/prgrs/clickup/ui/context"
	"github.com/prgrs/clickup/ui/widgets/navigator"
	"github.com/prgrs/clickup/ui/widgets/tasks"
)

const ViewId = "viewCompact"

type Model struct {
	ctx    *context.UserContext
	log    *log.Logger
	ViewId common.ViewId
	state  common.WidgetId
	size   common.Size

	spinner     spinner.Model
	showSpinner bool

	widgetNavigator navigator.Model
	widgetViewsTabs viewstabs.Model
	widgetTasks     tasks.Model
}

func (m Model) GetSize() common.Size {
	return m.size
}

func (m Model) GetViewId() common.ViewId {
	return m.ViewId
}

func (m Model) Init() tea.Cmd {
	m.log.Info("Initializing...")
	return tea.Batch(
		m.spinner.Tick,
		InitCompactCmd(),
	)
}

func (Model) KeyMap() help.KeyMap {
	return common.NewEmptyKeyMap()
}

func (m Model) Ready() bool {
	return !m.showSpinner
}

func (m Model) SetSize(size common.Size) common.View {
	m.size = size
	return m
}

func (m Model) Update(msg tea.Msg) (common.View, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "tab":
			switch m.state {
			case navigator.WidgetId:
				m.state = tasks.WidgetId
				m.widgetTasks = m.widgetTasks.SetFocused(true)
				m.widgetViewsTabs = m.widgetViewsTabs.SetFocused(false)
				m.widgetNavigator = m.widgetNavigator.SetFocused(false)
			case viewstabs.WidgetId:
				m.state = navigator.WidgetId
				m.widgetTasks = m.widgetTasks.SetFocused(false)
				m.widgetViewsTabs = m.widgetViewsTabs.SetFocused(false)
				m.widgetNavigator = m.widgetNavigator.SetFocused(true)
			case tasks.WidgetId:
				m.state = viewstabs.WidgetId
				m.widgetTasks = m.widgetTasks.SetFocused(false)
				m.widgetViewsTabs = m.widgetViewsTabs.SetFocused(true)
				m.widgetNavigator = m.widgetNavigator.SetFocused(false)
			}
		}

		switch m.state {
		case navigator.WidgetId:
			m.widgetNavigator, cmd = m.widgetNavigator.Update(msg)
		case viewstabs.WidgetId:
			m.widgetViewsTabs, cmd = m.widgetViewsTabs.Update(msg)
		case tasks.WidgetId:
			m.widgetTasks, cmd = m.widgetTasks.Update(msg)
		}

		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case InitCompactMsg:
		m.log.Info("Received: InitCompactMsg")

		if err := m.widgetNavigator.Init(); err != nil {
			cmds = append(cmds, common.ErrCmd(err))
			return m, tea.Batch(cmds...)
		}
		initWorkspace := m.widgetNavigator.GetWorkspace()

		views, err := m.ctx.Api.GetViewsFromWorkspace(initWorkspace)
		if err != nil {
			cmds = append(cmds, common.ErrCmd(err))
			return m, tea.Batch(cmds...)
		}

		if len(views) == 0 {
			m.widgetTasks.SetTasks(nil)
			m.widgetViewsTabs.SetTabs(nil)
		} else {
			tabs := viewsToTabs(views)
			m.widgetViewsTabs.SetTabs(tabs)

			initTab := m.widgetViewsTabs.SelectedTab

			if err := m.reloadTasks(initTab); err != nil {
				cmds = append(cmds, common.ErrCmd(err))
				return m, tea.Batch(cmds...)
			}
		}

		m.showSpinner = false

	case spinner.TickMsg:
		// m.log.Info("Received: spinner.TickMsg")
		if m.showSpinner {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case common.WorkspacePreviewMsg:
		id := string(msg)
		m.log.Infof("Received: WorkspacePreviewMsg: %s", id)
		return m, m.handleFolderChangePreview(id)

	case common.WorkspaceChangeMsg:
		id := string(msg)
		m.log.Infof("Received: WorkspaceChangeMsg: %s", id)
		cmds = append(cmds, m.handleFolderChangePreview(id))

	case common.SpacePreviewMsg:
		id := string(msg)
		m.log.Infof("Received: received SpacePreviewMsg: %s", id)
		return m, m.handleSpaceChangePreview(id)

	case common.SpaceChangeMsg:
		id := string(msg)
		m.log.Infof("Received: received SpaceChangeMsg: %s", id)
		cmds = append(cmds, m.handleSpaceChangePreview(id))

	case common.FolderPreviewMsg:
		id := string(msg)
		m.log.Infof("Received: FolderPreviewMsg: %s", id)
		return m, m.handleFolderChangePreview(id)

	case common.FolderChangeMsg:
		id := string(msg)
		m.log.Infof("Received: FolderChangeMsg: %s", id)
		cmds = append(cmds, m.handleFolderChangePreview(id))

	case common.ListPreviewMsg:
		id := string(msg)
		m.log.Infof("Received: ListPreviewMsg: %s", id)
		cmds = append(cmds, m.handleListChangePreview(id))

	case common.ListChangeMsg:
		id := string(msg)
		m.log.Info("Received: ListChangeMsg", "id", id)
		// TODO: make state change as func
		m.state = m.widgetTasks.WidgetId
		m.widgetTasks = m.widgetTasks.SetFocused(true)
		m.widgetViewsTabs = m.widgetViewsTabs.SetFocused(false)
		m.widgetNavigator = m.widgetNavigator.SetFocused(false)
		return m, m.handleListChangePreview(id)

	case viewstabs.TabChangedMsg:
		idx := string(msg)
		initTab := m.widgetViewsTabs.SelectedTab
		m.log.Info("Received: TabChangedMsg", "idx", idx, "id", initTab)

		m.widgetTasks.SetSpinner(true)
		cmds = append(cmds, LoadingTasksFromViewCmd(initTab))

	case LoadingTasksFromViewMsg:
		id := string(msg)
		m.widgetTasks.SetSpinner(false)

		if id == "" {
			m.log.Info("Received: LoadingTasksFromViewMsg empty")
			m.widgetTasks.SetTasks(nil)
			break
		}

		m.log.Info("Received: LoadingTasksFromViewMsg", "id", id)
		if err := m.reloadTasks(id); err != nil {
			cmds = append(cmds, common.ErrCmd(err))
			return m, tea.Batch(cmds...)
		}
	}

	m.widgetViewsTabs, cmd = m.widgetViewsTabs.Update(msg)
	cmds = append(cmds, cmd)

	m.widgetNavigator, cmd = m.widgetNavigator.Update(msg)
	cmds = append(cmds, cmd)

	m.widgetTasks, cmd = m.widgetTasks.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.showSpinner {
		return lipgloss.Place(
			m.ctx.WindowSize.Width, m.ctx.WindowSize.Height,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("%s Loading...", m.spinner.View()),
		)
	}

	m.widgetViewsTabs.SetSize(common.Size{
		Width: m.ctx.WindowSize.Width,
		// Height: m.ctx.WindowSize.Height - m.ctx.WindowSize.MetaHeight - lipgloss.Height(widgetViewsTabsRendered) - 2,
	})
	widgetViewsTabsRendered := m.widgetViewsTabs.View()

	m.widgetNavigator.SetSize(common.Size{
		Width:  25,
		Height: m.ctx.WindowSize.Height - lipgloss.Height(widgetViewsTabsRendered), // - 2,
		// Height: m.ctx.WindowSize.Height - m.ctx.WindowSize.MetaHeight - lipgloss.Height(widgetViewsTabsRendered) - 2,
	})
	widgetNavigatorRendered := m.widgetNavigator.View()

	m.widgetTasks.SetSize(common.Size{
		Width:  m.ctx.WindowSize.Width - lipgloss.Width(widgetNavigatorRendered),
		Height: m.ctx.WindowSize.Height - lipgloss.Height(widgetViewsTabsRendered), // - 1,
		// Height: m.ctx.WindowSize.Height - m.ctx.WindowSize.MetaHeight - lipgloss.Height(widgetViewsTabsRendered), // - 1,
	})
	widgetTasksRendered := m.widgetTasks.View()

	return lipgloss.JoinVertical(
		lipgloss.Top,
		widgetViewsTabsRendered,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			widgetNavigatorRendered,
			widgetTasksRendered,
		),
	)
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) common.View {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	log := logger.WithPrefix(logger.GetPrefix() + "/" + ViewId)

	var (
		widgetViewsTabs = viewstabs.InitialModel(ctx, log)
		widgetTasks     = tasks.InitialModel(ctx, log)
		widgetNavigator = navigator.InitialModel(ctx, log).SetFocused(true)
	)

	return Model{
		ViewId:          ViewId,
		ctx:             ctx,
		spinner:         s,
		showSpinner:     true,
		log:             log,
		widgetViewsTabs: widgetViewsTabs,
		widgetNavigator: widgetNavigator,
		widgetTasks:     widgetTasks,
		state:           widgetNavigator.WidgetId,
	}
}

func viewsToTabs(views []clickup.View) []viewstabs.Tab {
	tabs := make([]viewstabs.Tab, len(views))
	for i, view := range views {
		tabView := viewstabs.Tab{
			Name: view.Name,
			Type: "view",
			Id:   view.Id,
			// Active: false,
		}
		tabs[i] = tabView
	}

	return tabs
}

func (m *Model) reloadTasks(viewId string) error {
	tasks, err := m.ctx.Api.GetTasksFromView(viewId)
	if err != nil {
		return err
	}
	m.widgetTasks.SetTasks(tasks)
	return nil
}

func (m *Model) handleWorkspaceChangePreview(id string) tea.Cmd {
	views, err := m.ctx.Api.GetViewsFromWorkspace(id)
	if err != nil {
		return common.ErrCmd(err)
	}
	tabs := viewsToTabs(views)
	m.widgetViewsTabs.SetTabs(tabs)

	initTab := m.widgetViewsTabs.SelectedTab
	m.widgetTasks.SetSpinner(true)

	return LoadingTasksFromViewCmd(initTab)
}

func (m *Model) handleSpaceChangePreview(id string) tea.Cmd {
	views, err := m.ctx.Api.GetViewsFromSpace(id)
	if err != nil {
		return common.ErrCmd(err)
	}
	tabs := viewsToTabs(views)
	m.widgetViewsTabs.SetTabs(tabs)

	initTab := m.widgetViewsTabs.SelectedTab
	m.widgetTasks.SetSpinner(true)

	return LoadingTasksFromViewCmd(initTab)
}

func (m *Model) handleFolderChangePreview(id string) tea.Cmd {
	views, err := m.ctx.Api.GetViewsFromFolder(id)
	if err != nil {
		return common.ErrCmd(err)
	}
	tabs := viewsToTabs(views)
	m.widgetViewsTabs.SetTabs(tabs)

	initTab := m.widgetViewsTabs.SelectedTab
	m.widgetTasks.SetSpinner(true)

	return LoadingTasksFromViewCmd(initTab)
}

func (m *Model) handleListChangePreview(id string) tea.Cmd {
	views, err := m.ctx.Api.GetViewsFromList(id)
	if err != nil {
		return common.ErrCmd(err)
	}
	tabs := viewsToTabs(views)
	m.widgetViewsTabs.SetTabs(tabs)

	initTab := m.widgetViewsTabs.SelectedTab
	m.widgetTasks.SetSpinner(true)

	return LoadingTasksFromViewCmd(initTab)
}
