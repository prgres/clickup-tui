package lists

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/components/lists"
	"github.com/prgrs/clickup/ui/context"
)

type ListsState uint

const (
	ListsStateLoading ListsState = iota
	ListsStateList
)

type Model struct {
	ctx   *context.UserContext
	state ListsState

	componentListsList lists.Model

	spinner     spinner.Model
	showSpinner bool
}

func InitialModel(ctx *context.UserContext) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	return Model{
		ctx:                ctx,
		componentListsList: lists.InitialModel(ctx),
		state:              ListsStateList,

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
			m.ctx.Logger.Info("Hiding lists view")
			cmds = append(cmds, HideListsViewCmd())
		}

	case spinner.TickMsg:
		m.ctx.Logger.Info("ViewLists receive spinner.TickMsg")
		if m.showSpinner {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}

	case common.FolderChangeMsg:
		m.ctx.Logger.Infof("ViewLists receive FolderChangeMsg")
		m.showSpinner = true

	case lists.ListsListReadyMsg:
		m.ctx.Logger.Infof("ViewLists receive ListsListReadyMsg")
		m.showSpinner = false
	}

	m.componentListsList, cmd = m.componentListsList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.showSpinner {
		return lipgloss.Place(
			m.ctx.WindowSize.Width, m.ctx.WindowSize.Height,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("%s Loading lists...", m.spinner.View()),
		)
	}

	return m.componentListsList.View()
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Info("Initializing view: Lists")
	return m.componentListsList.Init()
}
