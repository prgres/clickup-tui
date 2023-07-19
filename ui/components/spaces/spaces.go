package spaces

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

type HideSpaceViewMsg bool

func HideSpaceViewCmd() tea.Cmd {
	return func() tea.Msg {
		return HideSpaceViewMsg(true)
	}
}

type SpaceChangeMsg string

func SpaceChangeCmd(space string) tea.Cmd {
	return func() tea.Msg {
		return SpaceChangeMsg(space)
	}
}

type SpaceListReloadedMsg []clickup.Space

type Model struct {
	ctx           *context.UserContext
	list          list.Model
	SelectedSpace string
	spaces        []clickup.Space
}

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		list: list.New([]list.Item{},
			list.NewDefaultDelegate(),
			0, 0),
		ctx:           ctx,
		SelectedSpace: "",
		spaces:        []clickup.Space{},
	}
}

func (m Model) syncList(spaces []clickup.Space) Model {
	m.ctx.Logger.Info("Synchronizing list")
	m.spaces = spaces

	sre_index := 0
	items := spaceListToItems(spaces)
	itemsList := itemListToItems(items)

	for i, item := range items {
		if item.desc == SPACE_SRE {
			sre_index = i
		}
	}

	m.list.SetItems(itemsList)
	m.list.Select(sre_index)
	return m
}

func itemListToItems(items []item) []list.Item {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = itemToListItem(item)
	}
	return listItems
}

func itemToListItem(item item) list.Item {
	return list.Item(item)
}

func spaceListToItems(spaces []clickup.Space) []item {
	items := make([]item, len(spaces))
	for i, space := range spaces {
		items[i] = spaceToItem(space)
	}
	return items
}

func spaceToItem(space clickup.Space) item {
	return item{
		space.Name,
		space.Id,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case SpaceListReloadedMsg:
		m.ctx.Logger.Info("SpaceView received SpaceListReloadedMsg")
		m = m.syncList(msg)
		return m, nil

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			selectedSpace := m.list.SelectedItem().(item).desc
			m.ctx.Logger.Info("Selected space %s", selectedSpace)
			m.SelectedSpace = selectedSpace
			return m, SpaceChangeCmd(selectedSpace)

		case "esc":
			m.ctx.Logger.Info("Hiding space view")
			return m, HideSpaceViewCmd()
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.list.View()
}

func (m Model) Init() tea.Msg {
	spaces, err := m.getSpaces(TEAM_RAMP_NETWORK)
	if err != nil {
		return common.ErrMsg(err)
	}

	return SpaceListReloadedMsg(spaces)
}

func (m Model) getSpaces(team string) ([]clickup.Space, error) {
	m.ctx.Logger.Info("Getting spaces for team: %s", team)
	client := m.ctx.Clickup

	spaces, err := client.GetSpaces(team)
	if err != nil {
		return nil, err
	}
	m.ctx.Logger.Info("Found %d spaces for team: %d", len(spaces), team)

	return spaces, nil
}
