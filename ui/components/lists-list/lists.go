package listslist

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

const ComponentId = "viewLists"

type Model struct {
	list         list.Model
	ctx          *context.UserContext
	log          *log.Logger
	ComponentId  common.ComponentId
	SelectedList clickup.List
	lists        []clickup.List
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	l := list.New([]list.Item{},
		list.NewDefaultDelegate(),
		0, 0)
	l.KeyMap.Quit.Unbind()
	l.SetShowHelp(false)
	l.Title = "Lists"

	log := logger.WithPrefix(logger.GetPrefix() + "/" + ComponentId)

	return Model{
		ComponentId:  ComponentId,
		list:         l,
		ctx:          ctx,
		SelectedList: clickup.List{},
		lists:        []clickup.List{},
		log:          log,
	}
}

func (m Model) KeyMap() help.KeyMap {
	return common.NewKeyMap(
		m.list.FullHelp,
		m.list.ShortHelp,
	).With(common.KeyBindingBack)
}

func (m *Model) syncList(lists []clickup.List) {
	m.log.Info("Synchronizing list")
	m.lists = lists

	items := NewListItem(lists)

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
			selectedList := m.list.SelectedItem().(listitem.Item).Data().(clickup.List)
			m.log.Info("Selected list", "id", selectedList.Id, "name", selectedList.Name)
			m.SelectedList = selectedList
			return m, ListChangedCmd(m.SelectedList.Id)

		case "J", "shift+down":
			m.list.CursorDown()
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedList := m.list.SelectedItem().(listitem.Item).Data().(clickup.List)
			m.log.Info("Selected list", "id", selectedList.Id, "name", selectedList.Name)
			m.SelectedList = selectedList
			return m, common.ListPreviewCmd(m.SelectedList.Id)

		case "K", "shift+up":
			m.list.CursorUp()
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedList := m.list.SelectedItem().(listitem.Item).Data().(clickup.List)
			m.log.Info("Selected list", "id", selectedList.Id, "name", selectedList.Name)
			m.SelectedList = selectedList
			return m, common.ListPreviewCmd(m.SelectedList.Id)
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

func (m *Model) SetSize(size common.Size) {
	m.list.SetSize(size.Width, size.Height)
}

func (m *Model) SpaceChanged(id string) error {
	folders, err := m.ctx.Api.GetLists(id)
	if err != nil {
		return err
	}
	m.syncList(folders)

	return nil
}

func NewListItem(items []clickup.List) []list.Item {
	result := make([]list.Item, len(items))
	for i, v := range items {
		result[i] = listitem.NewItem(v.Name, v.Id, v)
	}
	return result
}
