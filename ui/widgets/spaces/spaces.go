package spaces

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
	"github.com/prgrs/clickup/ui/context"
)

type Model struct {
	ctx           *context.UserContext
	list          list.Model
	SelectedSpace string
	spaces        []clickup.Space
	SelectedTeam  string
}

func InitialModel(ctx *context.UserContext) Model {
	l := list.New([]list.Item{},
		list.NewDefaultDelegate(),
		0, 0)
	l.KeyMap.Quit.Unbind()

	return Model{
		list:          l,
		ctx:           ctx,
		SelectedSpace: "",
		spaces:        []clickup.Space{},
	}
}

func (m *Model) syncList(spaces []clickup.Space) {
	m.ctx.Logger.Info("Synchronizing list")
	m.spaces = spaces

	sre_index := 0
	items := spaceListToItems(spaces)
	itemsList := listitem.ItemListToBubblesItems(items)

	for i, item := range items {
		if item.Description() == m.ctx.Config.DefaultSpace {
			sre_index = i
		}
	}

	m.list.SetItems(itemsList)
	m.list.Select(sre_index)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case SpaceListReloadedMsg:
		m.ctx.Logger.Info("SpaceView received SpaceListReloadedMsg")
		m.syncList(msg)
		cmds = append(cmds, SpaceListReadyCmd())

	case tea.WindowSizeMsg:
		m.ctx.Logger.Info("SpaceView received tea.WindowSizeMsg")
		m.list.SetSize(msg.Width, msg.Height)

	case common.TeamChangeMsg:
		m.ctx.Logger.Info("SpaceView received TeamChangeMsg")
		m.SelectedTeam = string(msg)
		cmds = append(cmds, m.getSpacesCmd())

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if m.list.SelectedItem() == nil {
				m.ctx.Logger.Info("SpaceView: list is empty")
				break
			}
			selectedSpace := listitem.BubblesItemToItem(m.list.SelectedItem()).Description()
			m.ctx.Logger.Infof("SpaceView: Selected space %s", selectedSpace)
			m.SelectedSpace = selectedSpace
			cmds = append(cmds, common.SpaceChangeCmd(selectedSpace))
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.list.View()
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Infof("Initializing component: spacesList")
	return m.getSpacesCmd()
}
