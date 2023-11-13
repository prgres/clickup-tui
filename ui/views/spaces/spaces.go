package spaces

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/components/spaces"
	"github.com/prgrs/clickup/ui/context"
)

type SpacesState uint

const (
	SpacesStateLoading SpacesState = iota
	SpacesStateList
)

type Model struct {
	ctx   *context.UserContext
	state SpacesState

	componentSpaceList spaces.Model

	spinner     spinner.Model
	showSpinner bool
}

func InitialModel(ctx *context.UserContext) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	return Model{
		ctx:                ctx,
		componentSpaceList: spaces.InitialModel(ctx),
		state:              SpacesStateList,

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
			m.ctx.Logger.Info("Hiding space view")
			cmds = append(cmds, HideSpaceViewCmd())
		}

	case spinner.TickMsg:
		m.ctx.Logger.Info("ViewSpaces receive spinner.TickMsg")
		if m.showSpinner {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case common.TeamChangeMsg:
		m.ctx.Logger.Infof("ViewSpaces receive TeamChangeMsg")
		m.showSpinner = true

	case spaces.SpaceListReadyMsg:
		m.ctx.Logger.Infof("ViewSpaces receive SpaceListReadyMsg")
		m.showSpinner = false
	}

	m.componentSpaceList, cmd = m.componentSpaceList.Update(msg)
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

	return m.componentSpaceList.View()
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Info("Initializing view: Spaces")
	return m.componentSpaceList.Init()
}
