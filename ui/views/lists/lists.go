package lists

import (
	tea "github.com/charmbracelet/bubbletea"
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
}

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		ctx:                ctx,
		componentListsList: lists.InitialModel(ctx),
		state:              ListsStateList,
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
	}

	m.componentListsList, cmd = m.componentListsList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.componentListsList.View()
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Info("Initializing view: Lists")
	return m.componentListsList.Init()
}
