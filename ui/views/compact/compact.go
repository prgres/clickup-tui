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

const id = "Compact"

type Model struct {
	id    common.Id
	ctx   *context.UserContext
	log   *log.Logger
	state common.Id
	size  common.Size

	spinner     spinner.Model
	showSpinner bool

	widgetNavigator *navigator.Model
	widgetViewsTabs *viewstabs.Model
	widgetTasks     *tasks.Model
}

func (m Model) Size() common.Size {
	return m.size
}

func (m Model) Id() common.Id {
	return m.id
}

func (m Model) Init() tea.Cmd {
	m.log.Info("Initializing...")
	return tea.Batch(
		m.spinner.Tick,
		InitCompactCmd(),
	)
}

func (m Model) KeyMap() help.KeyMap {
	switch m.state {
	case m.widgetNavigator.Id():
		return m.widgetNavigator.KeyMap()
	case m.widgetViewsTabs.Id():
		return m.widgetViewsTabs.KeyMap()
	case m.widgetTasks.Id():
		return m.widgetTasks.KeyMap()
	default:
		return common.NewEmptyKeyMap()
	}
}

func (m Model) Ready() bool {
	return !m.showSpinner
}

func (m *Model) SetSize(size common.Size) {
	m.size = size
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.widgetViewsTabs.Path = m.widgetNavigator.GetPath()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "tab":
			switch m.state {
			case m.widgetNavigator.Id():
				m.state = m.widgetTasks.Id()
				m.widgetTasks.SetFocused(true)
				m.widgetViewsTabs.SetFocused(false)
				m.widgetNavigator.SetFocused(false)
			case m.widgetViewsTabs.Id():
				m.state = m.widgetNavigator.Id()
				m.widgetTasks.SetFocused(false)
				m.widgetViewsTabs.SetFocused(false)
				m.widgetNavigator.SetFocused(true)
			case m.widgetTasks.Id():
				m.state = m.widgetViewsTabs.Id()
				m.widgetTasks.SetFocused(false)
				m.widgetViewsTabs.SetFocused(true)
				m.widgetNavigator.SetFocused(false)
			}
		}

		// m.getActiveElement().
		switch m.state {
		case m.widgetNavigator.Id():
			cmd = m.widgetNavigator.Update(msg)
		case m.widgetViewsTabs.Id():
			cmd = m.widgetViewsTabs.Update(msg)
		case m.widgetTasks.Id():
			cmd = m.widgetTasks.Update(msg)
		}

		m.widgetViewsTabs.Path = m.widgetNavigator.GetPath()

		cmds = append(cmds, cmd)
		return tea.Batch(cmds...)

	case InitCompactMsg:
		m.showSpinner = false
		m.log.Info("Received: InitCompactMsg")

		if err := m.widgetNavigator.Init(); err != nil {
			cmds = append(cmds, common.ErrCmd(err))
			return tea.Batch(cmds...)
		}

		initWorkspace := m.widgetNavigator.GetWorkspace()
		m.widgetNavigator.SetWorksapce(initWorkspace)
		cmds = append(cmds, common.WorkspacePreviewCmd(initWorkspace.Id))

		views, err := m.ctx.Api.GetViewsFromWorkspace(initWorkspace.Id)
		if err != nil {
			cmds = append(cmds, common.ErrCmd(err))
			return tea.Batch(cmds...)
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
				return tea.Batch(cmds...)
			}
		}

	case spinner.TickMsg:
		if m.showSpinner {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case common.WorkspacePreviewMsg:
		id := string(msg)
		m.log.Infof("Received: WorkspacePreviewMsg: %s", id)
		return m.handleWorkspaceChangePreview(id)

	case common.WorkspaceChangeMsg:
		id := string(msg)
		m.log.Infof("Received: WorkspaceChangeMsg: %s", id)
		cmds = append(cmds, m.handleWorkspaceChangePreview(id))

	case common.SpacePreviewMsg:
		id := string(msg)
		m.log.Infof("Received: received SpacePreviewMsg: %s", id)
		return m.handleSpaceChangePreview(id)

	case common.SpaceChangeMsg:
		id := string(msg)
		m.log.Infof("Received: received SpaceChangeMsg: %s", id)
		cmds = append(cmds, m.handleSpaceChangePreview(id))

	case common.FolderPreviewMsg:
		id := string(msg)
		m.log.Infof("Received: FolderPreviewMsg: %s", id)
		return m.handleFolderChangePreview(id)

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
		m.state = m.widgetTasks.Id()
		m.widgetTasks.SetFocused(true)
		m.widgetViewsTabs.SetFocused(false)
		m.widgetNavigator.SetFocused(false)
		return m.handleListChangePreview(id)

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
			return tea.Batch(cmds...)
		}

	case tasks.LostFocusMsg:
		m.log.Info("Received: tasks.LostFocusMsg")
		m.state = m.widgetNavigator.Id()
		m.widgetTasks.SetFocused(false)
		m.widgetViewsTabs.SetFocused(false)
		m.widgetNavigator.SetFocused(true)
		m.widgetViewsTabs.Path = m.widgetNavigator.GetPath()
	}

	cmd = m.widgetNavigator.Update(msg)
	cmds = append(cmds, cmd)

	cmd = m.widgetViewsTabs.Update(msg)
	cmds = append(cmds, cmd)

	cmd = m.widgetTasks.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
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

	size := m.ctx.WindowSize
	size.Height -= size.MetaHeight

	m.widgetViewsTabs.SetSize(common.Size{
		Width: size.Width,
	})
	widgetViewsTabsRendered := m.widgetViewsTabs.View()

	m.widgetNavigator.SetSize(common.Size{
		Width:  25,
		Height: size.Height - lipgloss.Height(widgetViewsTabsRendered),
	})
	widgetNavigatorRendered := m.widgetNavigator.View()

	m.widgetTasks.SetSize(common.Size{
		Width:  size.Width - lipgloss.Width(widgetNavigatorRendered),
		Height: size.Height - lipgloss.Height(widgetViewsTabsRendered),
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

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	log := logger.WithPrefix(logger.GetPrefix() + "/view/" + id)

	var (
		widgetViewsTabs = viewstabs.InitialModel(ctx, log)
		widgetTasks     = tasks.InitialModel(ctx, log)
		widgetNavigator = navigator.InitialModel(ctx, log).WithFocused(true)
	)

	return Model{
		id:              id,
		ctx:             ctx,
		spinner:         s,
		showSpinner:     true,
		log:             log,
		widgetViewsTabs: &widgetViewsTabs,
		widgetNavigator: &widgetNavigator,
		widgetTasks:     &widgetTasks,
		state:           widgetNavigator.Id(),
	}
}

func viewsToTabs(views []clickup.View) []viewstabs.Tab {
	tabs := make([]viewstabs.Tab, len(views))
	for i, view := range views {
		tabs[i] = viewstabs.Tab{
			Name: view.Name,
			Id:   view.Id,
		}
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

// func (m *Model) getActiveElement() common.UIElement {
// 	switch m.state {
// 	case m.widgetNavigator.Id():
// 		return m.widgetNavigator
// 	case m.widgetViewsTabs.Id():
// 		return m.widgetViewsTabs
// 	case m.widgetTasks.Id():
// 		return m.widgetTasks
// 	default:
// 		return nil
// 	}
// }
