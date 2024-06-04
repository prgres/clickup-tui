package listslist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
	"github.com/prgrs/clickup/ui/context"
)

const id = "lists-list"

type Model struct {
	id     common.Id
	list   list.Model
	ctx    *context.UserContext
	log    *log.Logger
	lists  []clickup.List
	keyMap KeyMap

	Selected clickup.List
}

func (m Model) Id() common.Id {
	return m.id
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	l := list.New([]list.Item{},
		list.NewDefaultDelegate(),
		0, 0)

	l.KeyMap.Quit.Unbind()
	l.KeyMap.CursorUp.Unbind()
	l.KeyMap.CursorDown.Unbind()

	l.SetShowHelp(false)
	l.Title = "Lists"

	log := common.NewLogger(logger, common.ResourceTypeRegistry.COMPONENT, id)

	return Model{
		id:       id,
		list:     l,
		ctx:      ctx,
		Selected: clickup.List{},
		lists:    []clickup.List{},
		keyMap:   DefaultKeyMap(),
		log:      log,
	}
}

func (m Model) KeyMap() KeyMap {
	return m.keyMap
}

func (m *Model) SetList(lists []clickup.List) {
	m.log.Info("Synchronizing list")
	m.lists = lists
	items := NewListItem(lists)
	m.list.SetItems(items)
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeys(msg)
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

func (m *Model) SetSize(size common.Size) {
	m.list.SetSize(size.Width, size.Height)
}

func (m *Model) FolderChanged(id string) error {
	lists, err := m.ctx.Api.GetLists(id)
	if err != nil {
		return err
	}
	m.SetList(lists)

	return nil
}

func NewListItem(items []clickup.List) []list.Item {
	result := make([]list.Item, len(items))
	for i, v := range items {
		result[i] = listitem.NewItem(v.Name, v.Id, v)
	}
	return result
}
