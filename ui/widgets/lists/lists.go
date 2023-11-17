package lists

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
	"github.com/prgrs/clickup/ui/context"
)

type Model struct {
	ctx            *context.UserContext
	list           list.Model
	lists          []clickup.List
	SelectedList   string
	SelectedFolder string
}

func InitialModel(ctx *context.UserContext) Model {
	l := list.New([]list.Item{},
		list.NewDefaultDelegate(),
		0, 0)
	l.KeyMap.Quit.Unbind()

	return Model{
		list:           l,
		ctx:            ctx,
		SelectedFolder: "",
		SelectedList:   "",
		lists:          []clickup.List{},
	}
}

func (m *Model) syncList(lists []clickup.List) {
	m.ctx.Logger.Info("Synchronizing list")
	m.lists = lists

	items := listsListToItems(lists)
	itemsList := listitem.ItemListToBubblesItems(items)

	m.list.SetItems(itemsList)
	m.list.Select(0)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ListsListReloadedMsg:
		m.ctx.Logger.Info("ListsView received ListsListReloadedMsg")
		m.syncList(msg)
		cmds = append(cmds, ListsListReadyCmd())

	case tea.WindowSizeMsg:
		m.ctx.Logger.Info("ListsView received tea.WindowSizeMsg")
		m.list.SetSize(msg.Width, msg.Height)

	case common.FolderChangeMsg:
		m.ctx.Logger.Infof("ListsView received FolderChangeMsg: %s", string(msg))
		m.SelectedFolder = string(msg)
		cmds = append(cmds, m.getListsCmd(m.SelectedFolder))

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if m.list.SelectedItem() == nil {
				m.ctx.Logger.Info("ListsView: list is empty")
				break
			}
			selectedList := listitem.BubblesItemToItem(m.list.SelectedItem()).Description()
			m.ctx.Logger.Infof("ListsView: Selected list %s", selectedList)
			m.SelectedList = selectedList
			cmds = append(cmds, common.ListChangeCmd(m.SelectedList))
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
	m.ctx.Logger.Infof("Initializing component: listsList")
	return nil
}

func (m Model) getListsCmd(folderId string) tea.Cmd {
	return func() tea.Msg {
		folders, err := m.ctx.Api.GetLists(folderId)
		if err != nil {
			return common.ErrMsg(err)
		}

		return ListsListReloadedMsg(folders)
	}
}
