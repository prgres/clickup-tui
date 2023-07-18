package spaces

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/context"
)

const (
	TEAM_RAMP_NETWORK   = "24301226"
	SPACE_SRE_LIST_COOL = "q5kna-61288"
	SPACE_SRE           = "48458830"
)

type Model struct {
	ctx           *context.UserContext
	list          list.Model
	SelectedSpace string
	spaces        []clickup.Space
	hidden        bool
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
		list:          list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		ctx:           ctx,
		SelectedSpace: "",
		spaces:        []clickup.Space{},
		hidden:        true,
	}
}

func (m Model) syncItems() Model {
	m.ctx.Logger.Info(fmt.Sprintf("sync items %d", len(m.spaces)))
	items := make([]list.Item, len(m.spaces))
	sre_index := 0
	for i := range m.spaces {
		if m.spaces[i].Id == SPACE_SRE {
			sre_index = i
		}
		items[i] = item{
			title: m.spaces[i].Name,
			desc:  m.spaces[i].Id,
		}
	}

	m.list.SetItems(items)
	m.list.Select(sre_index)
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case spacesMsg:
		m.ctx.Logger.Info("spacesMsg")
		m.spaces = msg
		m = m.syncItems()
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			m.SelectedSpace = m.list.SelectedItem().(item).desc
			m.ctx.Logger.Info(fmt.Sprintf("selected space %s", m.SelectedSpace))
			m.hidden = true
			cmds = append(cmds, hideSpaceView())
		}
	}

	m.list, cmd = m.list.Update(msg)

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func hideSpaceView() tea.Cmd {
	return func() tea.Msg {
		return HideSpaceMsg(true)
	}
}

type HideSpaceMsg bool

func (m Model) View() string {
	// if len(m.list.Items()) != len(m.spaces) {
	// items := make([]list.Item, 3)
	// items := make([]list.Item, len(m.spaces))
	// for i := range m.spaces {
	// 	items[i] = item{
	// 		title: "test",
	// 		desc:  "desc",
	// 	}
	// }
	// m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	// return m.ctx.Style.Tabs.Tab.Render(m.list.View())
	return m.list.View()

	// return docStyle.Render(m.list.View())
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m Model) Init() tea.Msg {
	client := m.ctx.Clickup
	m.ctx.Logger.Info("fetching spaces")

	spaces, err := client.GetSpaces(TEAM_RAMP_NETWORK)
	if err != nil {
		return tea.Quit
	}
	return spacesMsg(spaces)

	// return fetchSpaces(m.ctx.Clickup)
}

func fetchSpaces(clickup *clickup.Client) tea.Cmd {
	spaces, err := clickup.GetSpaces(TEAM_RAMP_NETWORK)
	if err != nil {
		return tea.Quit
	}
	return tea.Quit
	return func() tea.Msg {
		return spacesMsg(spaces)
	}

	// names := make([]string, len(spaces))
	// for i := range spaces {
	// 	names[i] = spaces[i].Name
	// }

	// return func() tea.Msg {
	// 	// return someMsg{id: id}
	// 	return spacesNames(names)
	// }
}

type spacesMsg []clickup.Space
