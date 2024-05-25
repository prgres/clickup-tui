package spaceslist

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
	"github.com/prgrs/clickup/ui/context"
)

const ComponentId = "spacesList"

type Model struct {
	list          list.Model
	ctx           *context.UserContext
	log           *log.Logger
	ComponentId   common.ComponentId
	SelectedSpace clickup.Space
	spaces        []clickup.Space
}

func (m Model) KeyMap() help.KeyMap {
	return common.NewKeyMap(
		m.list.FullHelp,
		m.list.ShortHelp,
	).With(common.KeyBindingBack)
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	l := list.New([]list.Item{},
		list.NewDefaultDelegate(),
		0, 0)
	l.KeyMap.Quit.Unbind()
	l.SetShowHelp(false)
	l.Title = "Spaces"

	log := logger.WithPrefix(logger.GetPrefix() + "/" + ComponentId)

	return Model{
		ComponentId:   ComponentId,
		list:          l,
		ctx:           ctx,
		SelectedSpace: clickup.Space{},
		spaces:        []clickup.Space{},
		log:           log,
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

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedSpace := m.list.SelectedItem().(listitem.Item).Data().(clickup.Space)
			m.log.Info("Selected space", "id", selectedSpace.Id, "name", selectedSpace.Name)
			m.SelectedSpace = selectedSpace
			return m, SpaceChangedCmd(selectedSpace.Id)

		case "J", "shift+down":
			m.list.CursorDown()
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedSpace := m.list.SelectedItem().(listitem.Item).Data().(clickup.Space)
			m.log.Info("Selected space", "id", selectedSpace.Id, "name", selectedSpace.Name)
			m.SelectedSpace = selectedSpace
			return m, common.SpacePreviewCmd(selectedSpace.Id)

		case "K", "shift+up":
			m.list.CursorUp()
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedSpace := m.list.SelectedItem().(listitem.Item).Data().(clickup.Space)
			m.log.Info("Selected space", "id", selectedSpace.Id, "name", selectedSpace.Name)
			m.SelectedSpace = selectedSpace
			return m, common.SpacePreviewCmd(selectedSpace.Id)
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
