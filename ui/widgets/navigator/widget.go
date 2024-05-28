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

const id = "navigator"

type Model struct {
	id          common.Id
	log         *log.Logger
	ctx         *context.UserContext
	size        common.Size
	Focused     bool
	Hidden      bool
	ifBorders   bool
	state       common.Id
	spinner     spinner.Model
	showSpinner bool

	componentWorkspacesList *workspaceslist.Model
	componentSpacesList     *spaceslist.Model
	componentFoldersList    *folderslist.Model
	componentListsList      *listslist.Model
}

func (m Model) Id() common.Id {
	return m.id
}

// TODO: refactor
func (m *Model) SetWorksapce(workspace clickup.Workspace) {
	m.componentWorkspacesList.SelectedWorkspace = workspace
}

func (m Model) GetPath() string {
	switch m.state {
	case m.componentWorkspacesList.Id():
		return "/"
	case m.componentSpacesList.Id():
		return "/" + m.componentWorkspacesList.SelectedWorkspace.Name
	case m.componentFoldersList.Id():
		return "/" + m.componentWorkspacesList.SelectedWorkspace.Name + "/" + m.componentSpacesList.SelectedSpace.Name
	case m.componentListsList.Id():
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

	log := logger.WithPrefix(logger.GetPrefix() + "/widget/" + id)

	size := common.Size{
		Width:  0,
		Height: 0,
	}

	var (
		componentWorkspacesList = workspaceslist.InitialModel(ctx, log).WithFocused(true)
		componentFoldersList    = folderslist.InitialModel(ctx, log)
		componentSpacesList     = spaceslist.InitialModel(ctx, log)
		cpomponentListsList     = listslist.InitialModel(ctx, log)
	)

	return Model{
		id:        id,
		ctx:       ctx,
		size:      size,
		Focused:   false,
		Hidden:    false,
		log:       log,
		ifBorders: true,

		componentWorkspacesList: &componentWorkspacesList,
		componentFoldersList:    &componentFoldersList,
		componentSpacesList:     &componentSpacesList,
		componentListsList:      &cpomponentListsList,

		state: componentWorkspacesList.Id(),

		spinner:     s,
		showSpinner: false,
	}
}

func (m Model) Help() help.KeyMap {
	switch m.state {
	case m.componentWorkspacesList.Id():
		return m.componentWorkspacesList.Help()
	case m.componentSpacesList.Id():
		return m.componentSpacesList.Help()
	case m.componentFoldersList.Id():
		return m.componentFoldersList.Help()
	case m.componentListsList.Id():
		return m.componentListsList.Help()
	default:
		return common.NewEmptyHelp()
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.log.Info("Received: Go to previous view")

			switch m.state {
			case m.componentSpacesList.Id():
				m.state = m.componentWorkspacesList.Id()
			case m.componentFoldersList.Id():
				m.state = m.componentSpacesList.Id()
			case m.componentListsList.Id():
				m.state = m.componentFoldersList.Id()
			}

			cmds = append(cmds, cmd)
			return tea.Batch(cmds...)
		}

		switch m.state {
		case m.componentWorkspacesList.Id():
			cmd = m.componentWorkspacesList.Update(msg)
		case m.componentSpacesList.Id():
			cmd = m.componentSpacesList.Update(msg)
		case m.componentFoldersList.Id():
			cmd = m.componentFoldersList.Update(msg)
		case m.componentListsList.Id():
			cmd = m.componentListsList.Update(msg)
		}

		cmds = append(cmds, cmd)
		return tea.Batch(cmds...)

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
			return tea.Batch(cmds...)
		}
		m.showSpinner = false
		m.state = m.componentSpacesList.Id()

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
			return tea.Batch(cmds...)
		}
		m.showSpinner = false
		m.state = m.componentFoldersList.Id()

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
		if err := m.componentListsList.FolderChanged(id); err != nil {
			cmds = append(cmds, common.ErrCmd(err))
			return tea.Batch(cmds...)
		}
		m.showSpinner = false
		m.state = m.componentListsList.Id()

	case listslist.ListChangedMsg:
		id := string(msg)
		m.log.Infof("Received: ListChangedMsg: %s", id)
		cmds = append(cmds, common.ListChangeCmd(id))
	}

	cmds = append(cmds,
		m.componentWorkspacesList.Update(msg),
		m.componentSpacesList.Update(msg),
		m.componentFoldersList.Update(msg),
		m.componentListsList.Update(msg),
	)

	return tea.Batch(cmds...)
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

	styleBorders := m.ctx.Style.Borders.
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
	case m.componentWorkspacesList.Id():
		m.componentWorkspacesList.SetSize(size)
		content = m.componentWorkspacesList.View()
	case m.componentSpacesList.Id():
		m.componentSpacesList.SetSize(size)
		content = m.componentSpacesList.View()
	case m.componentFoldersList.Id():
		m.componentFoldersList.SetSize(size)
		content = m.componentFoldersList.View()
	case m.componentListsList.Id():
		m.componentListsList.SetSize(size)
		content = m.componentListsList.View()
	default:
		content = "Unknown state"
	}

	return style.Render(content)
}

func (m *Model) SetFocused(f bool) {
	m.Focused = f
}

func (m Model) WithFocused(f bool) Model {
	m.Focused = f
	return m
}

func (m *Model) SetSize(s common.Size) {
	m.size = s
}

func (m Model) Size() common.Size {
	return m.size
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
