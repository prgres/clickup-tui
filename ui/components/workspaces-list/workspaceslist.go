package workspaceslist

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

const ComponentId = "workspacesList"

type Model struct {
	list              list.Model
	ctx               *context.UserContext
	log               *log.Logger
	ComponentId       common.ComponentId
	SelectedWorkspace clickup.Workspace
	workspaces        []clickup.Workspace
	ifBorders         bool
	Focused           bool
}

func (m Model) SetFocused(f bool) Model {
	m.Focused = f
	return m
}

func (m Model) KeyMap() help.KeyMap {
	return common.NewKeyMap(
		m.list.FullHelp,
		m.list.ShortHelp,
	)
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

	log := logger.WithPrefix(logger.GetPrefix() + "/" + ComponentId)

	return Model{
		ComponentId:       ComponentId,
		list:              l,
		ctx:               ctx,
		SelectedWorkspace: clickup.Workspace{},
		workspaces:        []clickup.Workspace{},
		log:               log,
		ifBorders:         true,
		Focused:           false,
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
			selectedWorkspace := m.list.SelectedItem().(listitem.Item).Data().(clickup.Workspace)
			m.log.Info("Selected workspace", "id", selectedWorkspace.Id, "name", selectedWorkspace.Name)
			m.SelectedWorkspace = selectedWorkspace
			return m, common.WorkspaceChangeCmd(selectedWorkspace.Id)

		case "J", "shift+down":
			m.list.CursorDown()
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedWorkspace := m.list.SelectedItem().(listitem.Item).Data().(clickup.Workspace)
			m.log.Info("Selected workspace", "id", selectedWorkspace.Id, "name", selectedWorkspace.Name)
			m.SelectedWorkspace = selectedWorkspace
			return m, common.WorkspacePreviewCmd(selectedWorkspace.Id)

		case "K", "shift+up":
			m.list.CursorUp()
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedWorkspace := m.list.SelectedItem().(listitem.Item).Data().(clickup.Workspace)
			m.log.Info("Selected workspace", "id", selectedWorkspace.Id, "name", selectedWorkspace.Name)
			m.SelectedWorkspace = selectedWorkspace
			return m, common.WorkspacePreviewCmd(selectedWorkspace.Id)
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
