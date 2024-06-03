package workspaceslist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
	"github.com/prgrs/clickup/ui/context"
)

const id = "workspaces-list"

type Model struct {
	id         common.Id
	list       list.Model
	ctx        *context.UserContext
	log        *log.Logger
	workspaces []clickup.Workspace
	ifBorders  bool
	Focused    bool
	keyMap     KeyMap

	Selected clickup.Workspace
}

func (m Model) Id() common.Id {
	return m.id
}

func (m Model) KeyMap() KeyMap {
	return m.keyMap
}

func (m *Model) SetFocused(f bool) {
	m.Focused = f
}

func (m Model) WithFocused(f bool) Model {
	m.Focused = f
	return m
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	l := list.New(
		[]list.Item{},
		list.NewDefaultDelegate(),
		0, 0,
	)

	l.KeyMap.Quit.Unbind()
	l.SetShowHelp(false)
	l.Title = "Workspaces"

	log := common.NewLogger(logger, common.ResourceTypeRegistry.COMPONENT, id)

	return Model{
		id:         id,
		list:       l,
		ctx:        ctx,
		Selected:   clickup.Workspace{},
		workspaces: []clickup.Workspace{},
		log:        log,
		ifBorders:  true,
		Focused:    false,
		keyMap:     DefaultKeyMap(),
	}
}

func (m *Model) syncList(workspaces []clickup.Workspace) {
	m.log.Info("Synchronizing list...")
	m.workspaces = workspaces

	items := NewListItem(workspaces)
	if len(items) == 0 {
		panic("list is empty")
	}

	index := 0
	for i, item := range items {
		it := item.(listitem.Item)
		if it.Title() == m.ctx.Config.DefaultWorkspace {
			index = i
		}
	}

	m.Selected = items[index].(listitem.Item).Data().(clickup.Workspace)
	m.list.SetItems(items)
	m.list.Select(index)
	m.log.Info("List synchronized")
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

func (m *Model) InitWorkspaces() error {
	m.log.Info("Received: InitWorkspacesMsg")
	workspaces, err := m.ctx.Api.GetWorkspaces()
	if err != nil {
		return err
	}

	m.Selected = workspaces[0]
	m.syncList(workspaces)

	return nil
}

func NewListItem(items []clickup.Workspace) []list.Item {
	result := make([]list.Item, len(items))
	for i, v := range items {
		result[i] = listitem.NewItem(v.Name, v.Id, v)
	}
	return result
}
