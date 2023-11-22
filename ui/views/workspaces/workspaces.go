package workspaces

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
	workspaceslist "github.com/prgrs/clickup/ui/widgets/workspaces-list"
)

const ViewId = "viewWorkspaces"

type Model struct {
	ViewId               common.ViewId
	ctx                  *context.UserContext
	widgetWorkspacesList workspaceslist.Model
	spinner              spinner.Model
	showSpinner          bool
	log                  *log.Logger
}

func (m Model) Ready() bool {
	return !m.showSpinner
}

func (m Model) KeyMap() help.KeyMap {
	return m.widgetWorkspacesList.KeyMap()
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	log := logger.WithPrefix(logger.GetPrefix() + "/" + ViewId)

	return Model{
		ViewId:               ViewId,
		ctx:                  ctx,
		widgetWorkspacesList: workspaceslist.InitialModel(ctx, log),
		spinner:              s,
		showSpinner:          true,
		log:                  log,
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
			cmds = append(cmds, common.BackToPreviousViewCmd(m.ViewId))
		}

	case spinner.TickMsg:
		// m.log.Info("Received: spinner.TickMsg")
		if m.showSpinner {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case common.WorkspaceChangeMsg:
		m.log.Infof("Received: WorkspaceChangeMsg")
		m.showSpinner = true
		cmds = append(cmds, m.spinner.Tick)

	case workspaceslist.WorkspaceListReadyMsg:
		m.log.Infof("Received: WorkspaceListReadyMsg")
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
	m.log.Info("Initializing...")
	return tea.Batch(
		m.spinner.Tick,
		m.widgetWorkspacesList.Init(),
	)
}
