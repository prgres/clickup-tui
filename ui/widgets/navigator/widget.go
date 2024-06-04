package navigator

import (
	"fmt"

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
	"golang.org/x/sync/errgroup"
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

func (m Model) State() common.Id {
	return m.state
}

// TODO: refactor
func (m *Model) SetWorksapce(workspace clickup.Workspace) {
	m.componentWorkspacesList.Selected = workspace
}

func (m Model) GetPath() string {
	switch m.state {
	case m.componentWorkspacesList.Id():
		return ""
	case m.componentSpacesList.Id():
		return "/" + m.componentWorkspacesList.Selected.Name
	case m.componentFoldersList.Id():
		return "/" + m.componentWorkspacesList.Selected.Name + "/" + m.componentSpacesList.Selected.Name
	case m.componentListsList.Id():
		if m.Focused {
			return "/" + m.componentWorkspacesList.Selected.Name + "/" + m.componentSpacesList.Selected.Name + "/" + m.componentFoldersList.Selected.Name
		}
		return "/" + m.componentWorkspacesList.Selected.Name + "/" + m.componentSpacesList.Selected.Name + "/" + m.componentFoldersList.Selected.Name + "/" + m.componentListsList.Selected.Name
	default:
		return ""
	}
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	log := common.NewLogger(logger, common.ResourceTypeRegistry.WIDGET, id)
	size := common.NewEmptySize()

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

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeys(msg)

	case workspaceslist.WorkspaceChangedMsg:
		id := string(msg)
		m.log.Infof("Received: WorkspaceChangeMsg: %s", id)
		m.showSpinner = true
		cmds = append(cmds,
			m.spinner.Tick,
			LoadingSpacesFromWorkspaceCmd(id),
			common.WorkspaceChangedCmd(id),
		)

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
		m.log.Infof("Received: spaceslist.SpaceChangedMsg: %s", id)
		m.showSpinner = true
		cmds = append(cmds,
			m.spinner.Tick,
			LoadingFoldersFromSpaceCmd(id),
			common.SpaceChangedCmd(id),
		)

	case LoadingFoldersFromSpaceMsg:
		id := string(msg)
		m.log.Infof("Received: LoadingFoldersFromSpaceMsg: %s", id)
		if err := m.componentFoldersList.SpaceChanged(id); err != nil {
			cmds = append(cmds, common.ErrCmd(err))
			return tea.Batch(cmds...)
		}
		m.showSpinner = false
		m.state = m.componentFoldersList.Id()

	case folderslist.FolderChangedMsg:
		id := string(msg)
		m.log.Infof("Received: FolderChangeMsg: %s", id)
		m.showSpinner = true
		cmds = append(cmds,
			m.spinner.Tick,
			LoadingListsFromFolderCmd(id),
			common.FolderChangedCmd(id),
		)

	case folderslist.FolderSelectedMsg:
		id := string(msg)
		m.log.Infof("Received: folderslist.FolderSelectedMsg: %s", id)
		m.state = "zxc" // m.componentListsList.Id()
		return nil

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
		// m.showSpinner = true
		cmds = append(cmds, m.spinner.Tick, ListChangedCmd(id))

	case listslist.ListSelectedMsg:
		id := string(msg)
		m.log.Infof("Received: ListSelectedMsg: %s", id)
		cmds = append(cmds, ListSelectedCmd(id))

	case common.RefreshMsg:
		m.log.Debug("Received: common.RefreshMsg")

		errgroup := new(errgroup.Group)

		errgroup.Go(func() error {
			t, err := m.ctx.Api.GetTeams()
			if err != nil {
				return err
			}
			m.componentWorkspacesList.SetList(t)
			return nil
		})

		errgroup.Go(func() error {
			id := m.componentWorkspacesList.Selected.Id
			if id != "" {
				t, err := m.ctx.Api.GetSpaces(id)
				if err != nil {
					return err
				}
				m.componentSpacesList.SetList(t)
			}
			return nil
		})

		errgroup.Go(func() error {
			id := m.componentSpacesList.Selected.Id
			if id != "" {
				t, err := m.ctx.Api.GetFolders(id)
				if err != nil {
					return err
				}
				m.componentFoldersList.SetList(t)
			}
			return nil
		})

		errgroup.Go(func() error {
			id := m.componentFoldersList.Selected.Id
			if id != "" {
				t, err := m.ctx.Api.GetLists(id)
				if err != nil {
					return err
				}
				m.componentListsList.SetList(t)
			}
			return nil
		})

		err := errgroup.Wait()
		if err != nil {
			return common.ErrCmd(err)
		}
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
	return m.componentWorkspacesList.Selected
}

func (m *Model) Init() error {
	if err := m.componentWorkspacesList.InitWorkspaces(); err != nil {
		return err
	}

	return nil
}
