package lists

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

const WidgetId = "viewLists"

type Model struct {
	WidgetId       common.WidgetId
	ctx            *context.UserContext
	list           list.Model
	lists          []clickup.List
	SelectedList   listitem.Item
	SelectedFolder string
	log            *log.Logger
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	l := list.New([]list.Item{},
		list.NewDefaultDelegate(),
		0, 0)
	l.KeyMap.Quit.Unbind()
	l.SetShowHelp(false)

	log := logger.WithPrefix(logger.GetPrefix() + "/" + WidgetId)

	return Model{
		WidgetId:       WidgetId,
		list:           l,
		ctx:            ctx,
		SelectedFolder: "",
		SelectedList:   listitem.Item{},
		lists:          []clickup.List{},
		log:            log,
	}
}

func (m Model) KeyMap() help.KeyMap {
	return common.NewKeyMap(
		m.list.FullHelp,
		m.list.ShortHelp,
	).With(common.KeyBindingBack)
}

func (m *Model) syncList(lists []clickup.List) {
	m.log.Info("Synchronizing list")
	m.lists = lists

	items := listsListToItems(lists)
	itemsList := listitem.ItemListToBubblesItems(items)

	m.list.SetItems(itemsList)
	m.list.Select(0)
	m.log.Info("List synchronized")
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ListsListReloadedMsg:
		m.log.Info("Received: ListsListReloadedMsg")
		m.syncList(msg)
		cmds = append(cmds, ListsListReadyCmd())

	case tea.WindowSizeMsg:
		m.log.Debug("Received: tea.WindowSizeMsg",
			"width", msg.Width,
			"height", msg.Height)
		m.list.SetSize(msg.Width, msg.Height)

	case common.FolderChangeMsg:
		m.log.Infof("Received: FolderChangeMsg: %s", string(msg))
		m.SelectedFolder = string(msg)
		cmds = append(cmds, m.getListsCmd(m.SelectedFolder))

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedList := listitem.BubblesItemToItem(m.list.SelectedItem())
			m.log.Infof("Selected list %s", selectedList)
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
	m.log.Infof("Initializing...")
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
