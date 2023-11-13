package folders

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/components/folders"
	"github.com/prgrs/clickup/ui/context"
)

type FoldersState uint

const (
	FoldersStateLoading FoldersState = iota
	FoldersStateList
)

type Model struct {
	ctx   *context.UserContext
	state FoldersState

	componentFoldersList folders.Model

	spinner     spinner.Model
	showSpinner bool
}

func InitialModel(ctx *context.UserContext) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	return Model{
		ctx:                  ctx,
		componentFoldersList: folders.InitialModel(ctx),
		state:                FoldersStateList,

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
			m.ctx.Logger.Info("Hiding folders view")
			cmds = append(cmds, HideFolderViewCmd())
		}

	case spinner.TickMsg:
		m.ctx.Logger.Info("ViewFolders receive spinner.TickMsg")
		if m.showSpinner {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case common.SpaceChangeMsg:
		m.ctx.Logger.Infof("ViewFolders received SpaceChangeMsg: %s", string(msg))
		m.showSpinner = true

	case folders.FoldersListReadyMsg:
		m.ctx.Logger.Infof("ViewFolders receive FoldersListReadyMsg")
		m.showSpinner = false
	}

	m.componentFoldersList, cmd = m.componentFoldersList.Update(msg)
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

	return m.componentFoldersList.View()
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Info("Initializing view: Folders")
	return m.componentFoldersList.Init()
}
