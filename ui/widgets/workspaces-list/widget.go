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

const WidgetId = "workspacesList"

type Model struct {
	WidgetId          common.WidgetId
	ctx               *context.UserContext
	list              list.Model
	SelectedWorkspace string
	workspaces        []clickup.Workspace
	log               *log.Logger
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	l := list.New([]list.Item{},
		list.NewDefaultDelegate(),
		0, 0)
	l.KeyMap.Quit.Unbind()

	log := logger.WithPrefix(logger.GetPrefix() + "/" + WidgetId)

	return Model{
		WidgetId:          WidgetId,
		list:              l,
		ctx:               ctx,
		SelectedWorkspace: "",
		workspaces:        []clickup.Workspace{},
		log:               log,
	}
}

func (m *Model) syncList(workspaces []clickup.Workspace) {
	m.log.Info("Synchronizing list...")
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
	m.log.Info("List synchronized")
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case WorkspaceListReloadedMsg:
		m.log.Info("Received: WorkspaceListReloadedMsg")
		m.syncList(msg)
		cmds = append(cmds, WorkspaceListReadyCmd())

	case tea.WindowSizeMsg:
		m.log.Debug("Received: tea.WindowSizeMsg",
			"width", msg.Width,
			"height", msg.Height)
		m.list.SetSize(msg.Width, msg.Height)

	case common.WorkspaceChangeMsg:
		m.log.Info("Received: WorkspaceChangeMsg")
		m.SelectedWorkspace = string(msg)
		cmds = append(cmds, m.getWorkspacesCmd())

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedWorkspace := listitem.BubblesItemToItem(m.list.SelectedItem()).Description()
			m.log.Info("Selected workspace", "workspace", selectedWorkspace)
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
	m.log.Infof("Initializing...")
	return m.getWorkspacesCmd()
}
