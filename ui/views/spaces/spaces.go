package spaces

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
	"github.com/prgrs/clickup/ui/widgets/spaces"
)

const ViewId = "viewSpaces"

type SpacesState uint

const (
	SpacesStateLoading SpacesState = iota
	SpacesStateList
)

type Model struct {
	widgetSpaceList spaces.Model
	ctx             *context.UserContext
	log             *log.Logger
	ViewId          common.ViewId
	spinner         spinner.Model
	size            common.Size
	state           SpacesState
	showSpinner     bool
}

func (m Model) Ready() bool {
	return !m.showSpinner
}

func (m Model) KeyMap() help.KeyMap {
	return m.widgetSpaceList.KeyMap()
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) common.View {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	log := logger.WithPrefix(logger.GetPrefix() + "/" + ViewId)

	return Model{
		ViewId:          ViewId,
		ctx:             ctx,
		widgetSpaceList: spaces.InitialModel(ctx, log),
		state:           SpacesStateList,
		spinner:         s,
		showSpinner:     true,
		log:             log,
	}
}

func (m Model) Update(msg tea.Msg) (common.View, tea.Cmd) {
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

	case spaces.SpaceListReadyMsg:
		m.log.Infof("Received: SpaceListReadyMsg")
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
	m.log.Info("Initializing...")
	return tea.Batch(
		m.spinner.Tick,
		m.widgetSpaceList.Init(),
	)
}

func (m Model) SetSize(size common.Size) common.View {
	m.size = size
	m.widgetSpaceList = m.widgetSpaceList.SetSize(size)
	return m
}

func (m Model) GetSize() common.Size {
	return m.size
}

func (m Model) GetViewId() common.ViewId {
	return m.ViewId
}
