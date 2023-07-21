package spaces

import (
	tea "github.com/charmbracelet/bubbletea"
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
}

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		ctx:                ctx,
		componentSpaceList: spaces.InitialModel(ctx),
		state:              SpacesStateList,
		// state:              SpacesStateLoading,
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
			return m, HideSpaceViewCmd()

		default:
			switch m.state {
			case SpacesStateList:
				m.componentSpaceList, cmd = m.componentSpaceList.Update(msg)
				return m, cmd
			}
		}

	case spaces.SpaceChangeMsg:
		m.ctx.Logger.Infof("UI received SpaceChangeMsg: %s", string(msg))
		return m, tea.Batch(
			SpaceChangeCmd(string(msg)),
		)
	}

	m.componentSpaceList, cmd = m.componentSpaceList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.componentSpaceList.View()
}

func (m Model) Init() tea.Msg {
	m.ctx.Logger.Infof("Initializing view: spaces")
	return m.componentSpaceList.Init()
}
