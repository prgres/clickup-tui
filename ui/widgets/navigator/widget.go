package navigator

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/ui/common"
	folderslist "github.com/prgrs/clickup/ui/components/folders-list"
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

	componentWorkspacesList workspaceslist.Model
	componentFoldersList    folderslist.Model
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	log := logger.WithPrefix(logger.GetPrefix() + "/" + WidgetId)

	size := common.Size{
		Width:  0,
		Height: 0,
	}

	var (
		componentWorkspacesList = workspaceslist.InitialModel(ctx, log).SetFocused(true)
		componentFoldersList    = folderslist.InitialModel(ctx, log)
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

		state: componentWorkspacesList.ComponentId,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			switch m.state {
			case workspaceslist.ComponentId:
				m.state = folderslist.ComponentId
			}
		case "b":
			m.state = workspaceslist.ComponentId
		}

		switch m.state {
		case workspaceslist.ComponentId:
			m.componentWorkspacesList, cmd = m.componentWorkspacesList.Update(msg)
		case folderslist.ComponentId:
			m.componentFoldersList, cmd = m.componentFoldersList.Update(msg)
		}

		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	m.componentWorkspacesList, cmd = m.componentWorkspacesList.Update(msg)
	cmds = append(cmds, cmd)

	m.componentFoldersList, cmd = m.componentFoldersList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	bColor := lipgloss.Color("#FFF")
	if m.Focused {
		bColor = lipgloss.Color("#8909FF")
	}

	var content string
	switch m.state {
	case workspaceslist.ComponentId:
		content = m.componentWorkspacesList.View()
	case folderslist.ComponentId:
		content = m.componentFoldersList.View()
	}

	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(bColor).
		BorderBottom(m.ifBorders).
		BorderRight(m.ifBorders).
		BorderTop(m.ifBorders).
		BorderLeft(m.ifBorders).
		Render(
			content,
		)
}

func (m Model) SetFocused(f bool) Model {
	m.Focused = f
	return m
}

func (m *Model) SetSize(s common.Size) {
	m.componentWorkspacesList.SetSize(s)
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
