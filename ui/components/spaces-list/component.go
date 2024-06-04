package spaceslist

import (
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
	id     common.Id
	list   list.Model
	ctx    *context.UserContext
	log    *log.Logger
	spaces []clickup.Space
	keyMap KeyMap

	Selected clickup.Space
}

func (m Model) Id() common.Id {
	return m.id
}

func (m Model) KeyMap() KeyMap {
	return m.keyMap
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

	log := common.NewLogger(logger, common.ResourceTypeRegistry.COMPONENT, id)

	return Model{
		id:       id,
		list:     l,
		ctx:      ctx,
		Selected: clickup.Space{},
		spaces:   []clickup.Space{},
		log:      log,
		keyMap:   DefaultKeyMap(),
	}
}

func (m *Model) SetList(spaces []clickup.Space) {
	m.log.Info("Synchronizing list...")
	m.spaces = spaces
	items := NewListItem(spaces)
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

func (m *Model) SetSize(s common.Size) {
	m.list.SetSize(s.Width, s.Height)
}

func (m *Model) WorkspaceChanged(id string) error {
	m.log.Infof("Received: WorkspaceChangeMsg: %s", id)

	spaces, err := m.ctx.Api.GetSpaces(id)
	if err != nil {
		return err
	}

	m.SetList(spaces)
	return nil
}

func NewListItem(items []clickup.Space) []list.Item {
	result := make([]list.Item, len(items))
	for i, v := range items {
		result[i] = listitem.NewItem(v.Name, v.Id, v)
	}
	return result
}
