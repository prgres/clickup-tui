package navigator

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	folderslist "github.com/prgrs/clickup/ui/components/folders-list"
	listslist "github.com/prgrs/clickup/ui/components/lists-list"
	spaceslist "github.com/prgrs/clickup/ui/components/spaces-list"
	workspaceslist "github.com/prgrs/clickup/ui/components/workspaces-list"
	"github.com/prgrs/clickup/ui/context"
)

const WidgetId = "navigator"

type Model struct {
	log       *log.Logger
	ctx       *context.UserContext
	WidgetId  common.WidgetId
	size      common.Size
	Focused   bool
	Hidden    bool
	ifBorders bool

	state common.ComponentId

	spinner     spinner.Model
	showSpinner bool

	componentWorkspacesList workspaceslist.Model
	componentSpacesList     spaceslist.Model
	componentFoldersList    folderslist.Model
	componentListsList      listslist.Model
}

// TODO: refactor
func (m *Model) SetWorksapce(workspace clickup.Workspace) {
	m.componentWorkspacesList.SelectedWorkspace = workspace
}

func (m Model) GetPath() string {
	switch m.state {
	case workspaceslist.ComponentId:
		return "/"
	case spaceslist.ComponentId:
		return "/" + m.componentWorkspacesList.SelectedWorkspace.Name
	case folderslist.ComponentId:
		return "/" + m.componentWorkspacesList.SelectedWorkspace.Name + "/" + m.componentSpacesList.SelectedSpace.Name
	case listslist.ComponentId:
		if m.Focused {
			return "/" + m.componentWorkspacesList.SelectedWorkspace.Name + "/" + m.componentSpacesList.SelectedSpace.Name + "/" + m.componentFoldersList.SelectedFolder.Name
		}
		return "/" + m.componentWorkspacesList.SelectedWorkspace.Name + "/" + m.componentSpacesList.SelectedSpace.Name + "/" + m.componentFoldersList.SelectedFolder.Name + "/" + m.componentListsList.SelectedList.Name
	default:
		return ""
	}
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	log := logger.WithPrefix(logger.GetPrefix() + "/" + WidgetId)

	size := common.Size{
		Width:  0,
		Height: 0,
	}

	var (
		componentWorkspacesList = workspaceslist.InitialModel(ctx, log).SetFocused(true)
		componentFoldersList    = folderslist.InitialModel(ctx, log)
		componentSpacesList     = spaceslist.InitialModel(ctx, log)
		cpomponentListsList     = listslist.InitialModel(ctx, log)
	)

	return Model{
		WidgetId:  WidgetId,
		ctx:       ctx,
		size:      size,
		Focused:   false,
		Hidden:    false,
		log:       log,
		ifBorders: true,

		componentWorkspacesList: componentWorkspacesList,
		componentFoldersList:    componentFoldersList,
		componentSpacesList:     componentSpacesList,
		componentListsList:      cpomponentListsList,

		state: componentWorkspacesList.ComponentId,

		spinner:     s,
		showSpinner: false,
	}
}

func (m Model) KeyMap() help.KeyMap {
	switch m.state {
	case workspaceslist.ComponentId:
		return m.componentWorkspacesList.KeyMap()
	case spaceslist.ComponentId:
		return m.componentSpacesList.KeyMap()
	case folderslist.ComponentId:
		return m.componentFoldersList.KeyMap()
	case listslist.ComponentId:
		return m.componentListsList.KeyMap()
	default:
		return common.NewEmptyKeyMap()
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.log.Info("Received: Go to previous view")

			switch m.state {
			case spaceslist.ComponentId:
				m.state = workspaceslist.ComponentId
			case folderslist.ComponentId:
				m.state = spaceslist.ComponentId
			case listslist.ComponentId:
				m.state = folderslist.ComponentId
			}

			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		switch m.state {
		case workspaceslist.ComponentId:
			m.componentWorkspacesList, cmd = m.componentWorkspacesList.Update(msg)
		case spaceslist.ComponentId:
			m.componentSpacesList, cmd = m.componentSpacesList.Update(msg)
		case folderslist.ComponentId:
			m.componentFoldersList, cmd = m.componentFoldersList.Update(msg)
		case listslist.ComponentId:
			m.componentListsList, cmd = m.componentListsList.Update(msg)
		}

		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case common.WorkspaceChangeMsg:
		id := string(msg)
		m.log.Infof("Received: WorkspaceChangeMsg: %s", id)
		m.showSpinner = true
		cmds = append(cmds, m.spinner.Tick, LoadingSpacesFromWorkspaceCmd(id))

	case spinner.TickMsg:
		if m.showSpinner {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case LoadingSpacesFromWorkspaceMsg:
		id := string(msg)
		m.log.Infof("Received: LoadingSpacesFromWorkspaceMsg: %s", id)
		if err := m.componentSpacesList.WorkspaceChanged(id); err != nil {
			cmds = append(cmds, common.ErrCmd(err))
			return m, tea.Batch(cmds...)
		}
		m.showSpinner = false
		m.state = spaceslist.ComponentId

	case spaceslist.SpaceChangedMsg:
		id := string(msg)
		m.log.Infof("Received: SpaceChangedMsg: %s", id)
		cmds = append(cmds, common.SpaceChangeCmd(id))

	case common.SpaceChangeMsg:
		id := string(msg)
		m.log.Infof("Received: received SpaceChangeMsg: %s", id)
		m.showSpinner = true
		cmds = append(cmds, m.spinner.Tick, LoadingFoldersFromSpaceCmd(id))

	case LoadingFoldersFromSpaceMsg:
		id := string(msg)
		m.log.Infof("Received: LoadingFoldersFromSpaceMsg: %s", id)
		if err := m.componentFoldersList.SpaceChanged(id); err != nil {
			cmds = append(cmds, common.ErrCmd(err))
			return m, tea.Batch(cmds...)
		}
		m.showSpinner = false
		m.state = folderslist.ComponentId

	case folderslist.FolderChangeMsg:
		id := string(msg)
		m.log.Infof("Received: FolderChangeMsg: %s", id)
		cmds = append(cmds, common.FolderChangeCmd(id))

	case common.FolderChangeMsg:
		id := string(msg)
		m.log.Infof("Received: FolderChangeMsg: %s", id)
		m.showSpinner = true
		cmds = append(cmds, m.spinner.Tick, LoadingListsFromFolderCmd(id))

	case LoadingListsFromFolderMsg:
		id := string(msg)
		m.log.Infof("Received: LoadingListsFromFolderMsg: %s", id)
		if err := m.componentListsList.SpaceChanged(id); err != nil {
			cmds = append(cmds, common.ErrCmd(err))
			return m, tea.Batch(cmds...)
		}
		m.showSpinner = false
		m.state = listslist.ComponentId

	case listslist.ListChangedMsg:
		id := string(msg)
		m.log.Infof("Received: ListChangedMsg: %s", id)
		cmds = append(cmds, common.ListChangeCmd(id))
	}

	m.componentWorkspacesList, cmd = m.componentWorkspacesList.Update(msg)
	cmds = append(cmds, cmd)

	m.componentSpacesList, cmd = m.componentSpacesList.Update(msg)
	cmds = append(cmds, cmd)

	m.componentFoldersList, cmd = m.componentFoldersList.Update(msg)
	cmds = append(cmds, cmd)

	m.componentListsList, cmd = m.componentListsList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	bColor := m.ctx.Theme.BordersColorInactive
	if m.Focused {
		bColor = m.ctx.Theme.BordersColorActive
	}

	borderMargin := 0
	if m.ifBorders {
		borderMargin = 2
	}

	styleBorders := m.ctx.Style.Borders.Copy().
		BorderForeground(bColor)

	style := lipgloss.NewStyle().
		Inherit(styleBorders).
		Width(m.size.Width - borderMargin).
		MaxWidth(m.size.Width + borderMargin).
		Height(m.size.Height - borderMargin).
		MaxHeight(m.size.Height + borderMargin)

	if m.showSpinner {
		return style.Render(
			lipgloss.Place(
				m.size.Width-borderMargin, m.size.Height-borderMargin,
				lipgloss.Center,
				lipgloss.Center,
				fmt.Sprintf("%s Loading lists...", m.spinner.View()),
			),
		)
	}

	size := common.Size{
		Width:  m.size.Width - borderMargin,
		Height: m.size.Height - borderMargin,
	}

	var content string
	switch m.state {
	case workspaceslist.ComponentId:
		m.componentWorkspacesList.SetSize(size)
		content = m.componentWorkspacesList.View()
	case spaceslist.ComponentId:
		m.componentSpacesList.SetSize(size)
		content = m.componentSpacesList.View()
	case folderslist.ComponentId:
		m.componentFoldersList.SetSize(size)
		content = m.componentFoldersList.View()
	case listslist.ComponentId:
		m.componentListsList.SetSize(size)
		content = m.componentListsList.View()
	default:
		content = "Unknown state"
	}

	return style.Render(content)
}

func (m Model) SetFocused(f bool) Model {
	m.Focused = f
	return m
}

func (m *Model) SetSize(s common.Size) {
	m.size = s
}

func (m Model) GetWorkspace() clickup.Workspace {
	return m.componentWorkspacesList.SelectedWorkspace
}

func (m *Model) Init() error {
	if err := m.componentWorkspacesList.InitWorkspaces(); err != nil {
		return err
	}

	return nil
}
