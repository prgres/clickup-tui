package workspaceslist

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

const id = "workspaces-list"

type Model struct {
	id                common.Id
	list              list.Model
	ctx               *context.UserContext
	log               *log.Logger
	SelectedWorkspace clickup.Workspace
	workspaces        []clickup.Workspace
	ifBorders         bool
	Focused           bool
	keyMap            KeyMap
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

func (m *Model) SetFocused(f bool) {
	m.Focused = f
}

func (m Model) WithFocused(f bool) Model {
	m.Focused = f
	return m
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
	l := list.New(
		[]list.Item{},
		list.NewDefaultDelegate(),
		0, 0,
	)

	l.KeyMap.Quit.Unbind()
	l.SetShowHelp(false)
	l.Title = "Workspaces"

	log := logger.WithPrefix(logger.GetPrefix() + "/component/" + id)

	return Model{
		id:                id,
		list:              l,
		ctx:               ctx,
		SelectedWorkspace: clickup.Workspace{},
		workspaces:        []clickup.Workspace{},
		log:               log,
		ifBorders:         true,
		Focused:           false,
		keyMap:            DefaultKeyMap(),
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

	m.SelectedWorkspace = items[index].(listitem.Item).Data().(clickup.Workspace)
	m.list.SetItems(items)
	m.list.Select(index)
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
			selectedWorkspace := m.list.SelectedItem().(listitem.Item).Data().(clickup.Workspace)
			m.log.Info("Selected workspace", "id", selectedWorkspace.Id, "name", selectedWorkspace.Name)
			m.SelectedWorkspace = selectedWorkspace
			return common.WorkspaceChangeCmd(selectedWorkspace.Id)

		case key.Matches(msg, m.keyMap.CursorDown):
			m.list.CursorDown()

		case key.Matches(msg, m.keyMap.CursorDownAndSelect):
			m.list.CursorDown()
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedWorkspace := m.list.SelectedItem().(listitem.Item).Data().(clickup.Workspace)
			m.log.Info("Selected workspace", "id", selectedWorkspace.Id, "name", selectedWorkspace.Name)
			m.SelectedWorkspace = selectedWorkspace
			return common.WorkspacePreviewCmd(selectedWorkspace.Id)

		case key.Matches(msg, m.keyMap.CursorUp):
			m.list.CursorUp()

		case key.Matches(msg, m.keyMap.CursorUpAndSelect):
			m.list.CursorUp()
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedWorkspace := m.list.SelectedItem().(listitem.Item).Data().(clickup.Workspace)
			m.log.Info("Selected workspace", "id", selectedWorkspace.Id, "name", selectedWorkspace.Name)
			m.SelectedWorkspace = selectedWorkspace
			return common.WorkspacePreviewCmd(selectedWorkspace.Id)
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

func (m *Model) InitWorkspaces() error {
	m.log.Info("Received: InitWorkspacesMsg")
	workspaces, err := m.ctx.Api.GetWorkspaces()
	if err != nil {
		return err
	}

	m.SelectedWorkspace = workspaces[0]
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
