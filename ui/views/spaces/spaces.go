package spaces

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
	"github.com/prgrs/clickup/ui/widgets/spaces"
)

type SpacesState uint

const (
	SpacesStateLoading SpacesState = iota
	SpacesStateList
)

type Model struct {
	ViewId          common.ViewId
	ctx             *context.UserContext
	state           SpacesState
	widgetSpaceList spaces.Model
	spinner         spinner.Model
	showSpinner     bool
}

func (m Model) Ready() bool {
	return !m.showSpinner
}

func InitialModel(ctx *context.UserContext) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	return Model{
		ViewId:          "viewSpaces",
		ctx:             ctx,
		widgetSpaceList: spaces.InitialModel(ctx),
		state:           SpacesStateList,
		spinner:         s,
		showSpinner:     true,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.ctx.Logger.Info("ViewSpaces: Go to previous view")
			cmds = append(cmds, common.BackToPreviousViewCmd(m.ViewId))
		}

	case spinner.TickMsg:
		// m.ctx.Logger.Info("ViewSpaces receive spinner.TickMsg")
		if m.showSpinner {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case common.WorkspaceChangeMsg:
		m.ctx.Logger.Infof("ViewSpaces receive WorkspaceChangeMsg")
		m.showSpinner = true
		cmds = append(cmds, m.spinner.Tick)

	case spaces.SpaceListReadyMsg:
		m.ctx.Logger.Infof("ViewSpaces receive SpaceListReadyMsg")
		m.showSpinner = false
	}

	m.widgetSpaceList, cmd = m.widgetSpaceList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.showSpinner {
		return lipgloss.Place(
			m.ctx.WindowSize.Width, m.ctx.WindowSize.Height,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("%s Loading spaces...", m.spinner.View()),
		)
	}

	return m.widgetSpaceList.View()
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Info("Initializing view: Spaces")
	return tea.Batch(
		m.spinner.Tick,
		m.widgetSpaceList.Init(),
	)
}
