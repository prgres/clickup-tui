package spaceslist

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
	"github.com/prgrs/clickup/ui/context"
)

const id = "spaces-list"

type Model struct {
	id            common.Id
	list          list.Model
	ctx           *context.UserContext
	log           *log.Logger
	SelectedSpace clickup.Space
	spaces        []clickup.Space
	keyMap        KeyMap
}

func (m Model) Id() common.Id {
	return m.id
}

type KeyMap struct {
	CursorUp            key.Binding
	CursorUpAndSelect   key.Binding
	CursorDown          key.Binding
	CursorDownAndSelect key.Binding
	Select              key.Binding
}

func (m Model) KeyMap() KeyMap {
	return m.keyMap
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		CursorUp: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k, up", "up"),
		),
		CursorUpAndSelect: key.NewBinding(
			key.WithKeys("K", "shift+up"),
			key.WithHelp("K, shift+up", "up and select"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j, down", "down"),
		),
		CursorDownAndSelect: key.NewBinding(
			key.WithKeys("J", "shift+down"),
			key.WithHelp("J, down", "down and select"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

func (m Model) Help() help.KeyMap {
	return common.NewHelp(
		func() [][]key.Binding {
			return append(
				m.list.FullHelp(),
				[]key.Binding{
					m.keyMap.CursorUp,
					m.keyMap.CursorUpAndSelect,
					m.keyMap.CursorDown,
					m.keyMap.CursorDownAndSelect,
					m.keyMap.Select,
				},
			)
		},
		func() []key.Binding {
			return append(
				m.list.ShortHelp(),
				m.keyMap.CursorUp,
				m.keyMap.CursorDown,
				m.keyMap.Select,
			)
		},
	).With(common.KeyBindingBack)
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	l := list.New([]list.Item{},
		list.NewDefaultDelegate(),
		0, 0)

	l.KeyMap.Quit.Unbind()
	l.KeyMap.CursorUp.Unbind()
	l.KeyMap.CursorDown.Unbind()

	l.SetShowHelp(false)
	l.Title = "Spaces"

	log := logger.WithPrefix(logger.GetPrefix() + "/component/" + id)

	return Model{
		id:            id,
		list:          l,
		ctx:           ctx,
		SelectedSpace: clickup.Space{},
		spaces:        []clickup.Space{},
		log:           log,
		keyMap:        DefaultKeyMap(),
	}
}

func (m *Model) syncList(spaces []clickup.Space) {
	m.log.Info("Synchronizing list...")
	m.spaces = spaces

	items := NewListItem(spaces)

	for _, item := range items {
		i := item.(listitem.Item)
		if i.Title() == m.ctx.Config.DefaultSpace {
			m.SelectedSpace = i.Data().(clickup.Space)
		}
	}

	m.list.SetItems(items)
	m.list.Select(0)
	m.log.Info("List synchronized")
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Select):
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedSpace := m.list.SelectedItem().(listitem.Item).Data().(clickup.Space)
			m.log.Info("Selected space", "id", selectedSpace.Id, "name", selectedSpace.Name)
			m.SelectedSpace = selectedSpace
			return SpaceChangedCmd(selectedSpace.Id)

		case key.Matches(msg, m.keyMap.CursorDown):
			m.list.CursorDown()

		case key.Matches(msg, m.keyMap.CursorDownAndSelect):
			m.list.CursorDown()
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedSpace := m.list.SelectedItem().(listitem.Item).Data().(clickup.Space)
			m.log.Info("Selected space", "id", selectedSpace.Id, "name", selectedSpace.Name)
			m.SelectedSpace = selectedSpace
			return common.SpacePreviewCmd(selectedSpace.Id)

		case key.Matches(msg, m.keyMap.CursorUp):
			m.list.CursorUp()

		case key.Matches(msg, m.keyMap.CursorUpAndSelect):
			m.list.CursorUp()
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedSpace := m.list.SelectedItem().(listitem.Item).Data().(clickup.Space)
			m.log.Info("Selected space", "id", selectedSpace.Id, "name", selectedSpace.Name)
			m.SelectedSpace = selectedSpace
			return common.SpacePreviewCmd(selectedSpace.Id)
		}
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.list.View()
}

func (m Model) Init() tea.Cmd {
	m.log.Infof("Initializing...")
	return nil
}

func (m *Model) SetSize(s common.Size) {
	m.list.SetSize(s.Width, s.Height)
}

func (m *Model) WorkspaceChanged(id string) error {
	m.log.Infof("Received: WorkspaceChangeMsg: %s", id)

	spaces, err := m.ctx.Api.GetSpaces(id)
	if err != nil {
		return err
	}

	m.syncList(spaces)
	return nil
}

func NewListItem(items []clickup.Space) []list.Item {
	result := make([]list.Item, len(items))
	for i, v := range items {
		result[i] = listitem.NewItem(v.Name, v.Id, v)
	}
	return result
}
