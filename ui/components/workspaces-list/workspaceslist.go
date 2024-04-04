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
	SelectedWorkspace string
	workspaces        []clickup.Workspace
	ifBorders         bool
	Focused           bool
}

func (m Model) GetFocused() bool {
	return m.Focused
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
	l := list.New([]list.Item{},
		list.NewDefaultDelegate(),
		0, 0)
	l.KeyMap.Quit.Unbind()
	l.SetShowHelp(false)

	log := logger.WithPrefix(logger.GetPrefix() + "/" + ComponentId)

	return Model{
		ComponentId:       ComponentId,
		list:              l,
		ctx:               ctx,
		SelectedWorkspace: "",
		workspaces:        []clickup.Workspace{},
		log:               log,
		ifBorders:         true,
		Focused:           false,
	}
}

func (m *Model) syncList(workspaces []clickup.Workspace) {
	m.log.Info("Synchronizing list...")
	m.workspaces = workspaces

	sre_index := 0
	items := workspaceListToItems(workspaces)
	itemsList := listitem.ItemListToBubblesItems(items)
	if len(items) == 0 {
		panic("list is empty")
	}

	for i, item := range items {
		if item.Description() == m.ctx.Config.DefaultWorkspace {
			sre_index = i
			m.SelectedWorkspace = item.Description()
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
			return m, common.WorkspaceChangeCmd(selectedWorkspace)

		case "J", "shift+down":
			m.list.CursorDown()
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedWorkspace := listitem.BubblesItemToItem(m.list.SelectedItem()).Description()
			m.log.Info("Selected workspace", "workspace", selectedWorkspace)
			m.SelectedWorkspace = selectedWorkspace
			return m, common.WorkspacePreviewCmd(selectedWorkspace)

		case "K", "shift+up":
			m.list.CursorUp()
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedWorkspace := listitem.BubblesItemToItem(m.list.SelectedItem()).Description()
			m.log.Info("Selected workspace", "workspace", selectedWorkspace)
			m.SelectedWorkspace = selectedWorkspace
			return m, common.WorkspacePreviewCmd(selectedWorkspace)
		}

		// switch {
		// case key.Matches(msg, m.list.KeyMap.CursorDown):
		// 	m.list.CursorUp()
		// 	return m, nil
		// }
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// func (m Model) View() string {
func (m Model) View() string {
	return m.list.View()
}

// 	bColor := lipgloss.Color("#FFF")
// 	if m.Focused {
// 		bColor = lipgloss.Color("#8909FF")
// 	}
//
// 	return lipgloss.NewStyle().
// 		BorderStyle(lipgloss.RoundedBorder()).
// 		BorderForeground(bColor).
// 		BorderBottom(m.ifBorders).
// 		BorderRight(m.ifBorders).
// 		BorderTop(m.ifBorders).
// 		BorderLeft(m.ifBorders).
// 		Render(
// 			m.list.View(),
// 		)
// }

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

	m.SelectedWorkspace = workspaces[0].Id
	m.syncList(workspaces)

	return nil
}
