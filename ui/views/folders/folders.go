package folders

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
	"github.com/prgrs/clickup/ui/widgets/folders"
)

type FoldersState uint

const (
	FoldersStateLoading FoldersState = iota
	FoldersStateList
)

type Model struct {
	ViewId            common.ViewId
	ctx               *context.UserContext
	state             FoldersState
	widgetFoldersList folders.Model
	spinner           spinner.Model
	showSpinner       bool
}

func InitialModel(ctx *context.UserContext) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	return Model{
		ViewId:            "viewFolders",
		ctx:               ctx,
		widgetFoldersList: folders.InitialModel(ctx),
		state:             FoldersStateList,
		spinner:           s,
		showSpinner:       true,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.ctx.Logger.Info("ViewFolders: Go to previous view")
			cmds = append(cmds, common.BackToPreviousViewCmd(m.ViewId))
		}

	case spinner.TickMsg:
		// m.ctx.Logger.Info("ViewFolders receive spinner.TickMsg")
		if m.showSpinner {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case common.SpaceChangeMsg:
		m.ctx.Logger.Infof("ViewFolders received SpaceChangeMsg: %s", string(msg))
		m.showSpinner = true
		cmds = append(cmds, m.spinner.Tick)

	case folders.FoldersListReadyMsg:
		m.ctx.Logger.Infof("ViewFolders receive FoldersListReadyMsg")
		m.showSpinner = false
	}

	m.widgetFoldersList, cmd = m.widgetFoldersList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.showSpinner {
		return lipgloss.Place(
			m.ctx.WindowSize.Width, m.ctx.WindowSize.Height,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("%s Loading folders...", m.spinner.View()),
		)
	}

	return m.widgetFoldersList.View()
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Info("Initializing view: Folders")
	return tea.Batch(
		m.spinner.Tick,
		m.widgetFoldersList.Init(),
	)
}
