// <<<<<<< Updated upstream
// package navigator
//
// import (
//
//	tea "github.com/charmbracelet/bubbletea"
//	"github.com/charmbracelet/lipgloss"
//	"github.com/charmbracelet/log"
//	"github.com/prgrs/clickup/ui/common"
//	folderslist "github.com/prgrs/clickup/ui/components/folders-list"
//	workspaceslist "github.com/prgrs/clickup/ui/components/workspaces-list"
//	"github.com/prgrs/clickup/ui/context"
//
// )
//
//	func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
//		var cmd tea.Cmd
//		var cmds []tea.Cmd
//
//		switch msg := msg.(type) {
//		case tea.KeyMsg:
//			switch keypress := msg.String(); keypress {
//			case "enter":
//				switch m.state {
//				case workspaceslist.ComponentId:
//					m.state = folderslist.ComponentId
//				}
//			case "b":
//				m.state = workspaceslist.ComponentId
//			}
//
//			switch m.state {
//			case workspaceslist.ComponentId:
//				m.componentWorkspacesList, cmd = m.componentWorkspacesList.Update(msg)
//			case folderslist.ComponentId:
//				m.componentFoldersList, cmd = m.componentFoldersList.Update(msg)
//			}
//
//			cmds = append(cmds, cmd)
//			return m, tea.Batch(cmds...)
//		}
//
//		m.componentWorkspacesList, cmd = m.componentWorkspacesList.Update(msg)
//		cmds = append(cmds, cmd)
//
//		m.componentFoldersList, cmd = m.componentFoldersList.Update(msg)
//		cmds = append(cmds, cmd)
//
//		return m, tea.Batch(cmds...)
//	}
//
//	func (m Model) View() string {
//		bColor := lipgloss.Color("#FFF")
//		if m.Focused {
//			bColor = lipgloss.Color("#8909FF")
//		}
//
//		var content string
//		switch m.state {
//		case workspaceslist.ComponentId:
//			content = m.componentWorkspacesList.View()
//		case folderslist.ComponentId:
//			content = m.componentFoldersList.View()
//		}
//
//		return lipgloss.NewStyle().
//			BorderStyle(lipgloss.RoundedBorder()).
//			BorderForeground(bColor).
//			BorderBottom(m.ifBorders).
//			BorderRight(m.ifBorders).
//			BorderTop(m.ifBorders).
//			BorderLeft(m.ifBorders).
//			Render(
//				content,
//			)
//	}
//
//	func (m Model) SetFocused(f bool) Model {
//		m.Focused = f
//		return m
//	}
//
//	func (m *Model) SetSize(s common.Size) {
//		m.componentWorkspacesList.SetSize(s)
//	}
//
//	func (m Model) GetWorkspace() string {
//		return m.componentWorkspacesList.SelectedWorkspace
//	}
//
//	func (m *Model) Init() error {
//		if err := m.componentWorkspacesList.InitWorkspaces(); err != nil {
//			return err
//		}
//
//		return nil
//	}
//
// ||||||| Stash base
// =======
package navigator

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
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

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.log.Info("Received: Go to previous view")

			switch m.state {
			// case workspaceslist.ComponentId:
			case spaceslist.ComponentId:
				m.state = workspaceslist.ComponentId
			// 	cmd = common.WorkspaceChangeCmd(m.componentWorkspacesList.SelectedWorkspace)
			case folderslist.ComponentId:
				m.state = spaceslist.ComponentId
				// cmd = common.SpaceChangeCmd(m.componentSpacesList.SelectedSpace)
			case listslist.ComponentId:
				m.state = folderslist.ComponentId
				// cmd = common.FolderChangeCmd(m.componentFoldersList.SelectedFolder)
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

		switch keypress := msg.String(); keypress {
		// case "enter":
		// 	switch m.state {
		// 	case workspaceslist.ComponentId:
		// 		m.state = folderslist.ComponentId
		// 	}
		// case "b":
		// 	m.state = workspaceslist.ComponentId
		}

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

	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(bColor).
		BorderBottom(m.ifBorders).
		BorderRight(m.ifBorders).
		BorderTop(m.ifBorders).
		BorderLeft(m.ifBorders).
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

func (m Model) GetWorkspace() string {
	return m.componentWorkspacesList.SelectedWorkspace
}

func (m *Model) Init() error {
	if err := m.componentWorkspacesList.InitWorkspaces(); err != nil {
		return err
	}

	return nil
}
