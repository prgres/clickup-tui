package workspaces

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
	"github.com/prgrs/clickup/ui/widgets/workspaces-list"
)

type Model struct {
	ViewId               common.ViewId
	ctx                  *context.UserContext
	widgetWorkspacesList workspaceslist.Model
	spinner              spinner.Model
	showSpinner          bool
}

func (m Model) Ready() bool {
	return !m.showSpinner
}

func InitialModel(ctx *context.UserContext) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	return Model{
		ViewId:               "viewWorkspaces",
		ctx:                  ctx,
		widgetWorkspacesList: workspaceslist.InitialModel(ctx),
		spinner:              s,
		showSpinner:          true,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.ctx.Logger.Info("ViewWorkspaces: Go to previous view")
			cmds = append(cmds, common.BackToPreviousViewCmd(m.ViewId))
		}

	case spinner.TickMsg:
		// m.ctx.Logger.Info("ViewWorkspaces receive spinner.TickMsg")
		if m.showSpinner {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case common.WorkspaceChangeMsg:
		m.ctx.Logger.Infof("ViewWorkspaces receive WorkspaceChangeMsg")
		m.showSpinner = true
		cmds = append(cmds, m.spinner.Tick)

	case workspaceslist.WorkspaceListReadyMsg:
		m.ctx.Logger.Infof("ViewWorkspaces receive WorkspaceListReadyMsg")
		m.showSpinner = false
	}

	m.widgetWorkspacesList, cmd = m.widgetWorkspacesList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.showSpinner {
		return lipgloss.Place(
			m.ctx.WindowSize.Width, m.ctx.WindowSize.Height,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("%s Loading workspaces...", m.spinner.View()),
		)
	}

	return m.widgetWorkspacesList.View()
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Info("Initializing view: Workspaces")
	return tea.Batch(
		m.spinner.Tick,
		m.widgetWorkspacesList.Init(),
	)

}
