package workspaceslist

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
	"github.com/prgrs/clickup/ui/context"
)

type Model struct {
	ctx               *context.UserContext
	list              list.Model
	SelectedWorkspace string
	workspaces        []clickup.Workspace
}

func InitialModel(ctx *context.UserContext) Model {
	l := list.New([]list.Item{},
		list.NewDefaultDelegate(),
		0, 0)
	l.KeyMap.Quit.Unbind()

	return Model{
		list:              l,
		ctx:               ctx,
		SelectedWorkspace: "",
		workspaces:        []clickup.Workspace{},
	}
}

func (m *Model) syncList(workspaces []clickup.Workspace) {
	m.ctx.Logger.Info("Synchronizing list")
	m.workspaces = workspaces

	sre_index := 0
	items := workspaceListToItems(workspaces)
	itemsList := listitem.ItemListToBubblesItems(items)

	for i, item := range items {
		if item.Description() == m.ctx.Config.DefaultWorkspace {
			sre_index = i
		}
	}

	m.list.SetItems(itemsList)
	m.list.Select(sre_index)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case WorkspaceListReloadedMsg:
		m.ctx.Logger.Info("WorkspaceView received WorkspaceListReloadedMsg")
		m.syncList(msg)
		cmds = append(cmds, WorkspaceListReadyCmd())

	case tea.WindowSizeMsg:
		m.ctx.Logger.Info("WorkspaceView received tea.WindowSizeMsg")
		m.list.SetSize(msg.Width, msg.Height)

	case common.WorkspaceChangeMsg:
		m.ctx.Logger.Info("WorkspaceView received WorkspaceChangeMsg")
		m.SelectedWorkspace = string(msg)
		cmds = append(cmds, m.getWorkspacesCmd())

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if m.list.SelectedItem() == nil {
				m.ctx.Logger.Info("WorkspaceView: list is empty")
				break
			}
			selectedWorkspace := listitem.BubblesItemToItem(m.list.SelectedItem()).Description()
			m.ctx.Logger.Infof("WorkspaceView: Selected workspace %s", selectedWorkspace)
			m.SelectedWorkspace = selectedWorkspace
			cmds = append(cmds, common.WorkspaceChangeCmd(selectedWorkspace))
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
	m.ctx.Logger.Infof("Initializing component: workspacesList")
	return m.getWorkspacesCmd()
}
