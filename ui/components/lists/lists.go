package lists

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

type Model struct {
	ctx            *context.UserContext
	list           list.Model
	lists          []clickup.List
	SelectedList   string
	SelectedFolder string
}

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

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
	itemsList := itemListToItems(items)

	m.list.SetItems(itemsList)
	m.list.Select(0)
}

func itemListToItems(items []item) []list.Item {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = itemToListItem(item)
	}
	return listItems
}

func itemToListItem(item item) list.Item {
	return list.Item(item)
}

func listsListToItems(lists []clickup.List) []item {
	items := make([]item, len(lists))
	for i, list := range lists {
		items[i] = listToItem(list)
	}
	return items
}

func listToItem(list clickup.List) item {
	return item{
		list.Name,
		list.Id,
	}
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
			selectedList := m.list.SelectedItem().(item).desc
			m.ctx.Logger.Infof("Selected list %s", selectedList)
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
	// return common.FolderChangeCmd(m.SelectedFolder)
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
