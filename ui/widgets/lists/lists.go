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
	list           list.Model
	ctx            *context.UserContext
	log            *log.Logger
	WidgetId       common.WidgetId
	SelectedList   string
	SelectedFolder string
	lists          []clickup.List
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
		SelectedList:   "",
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
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedList := listitem.BubblesItemToItem(m.list.SelectedItem()).Description()
			m.log.Infof("Selected list %s", selectedList)
			m.SelectedList = selectedList
			cmds = append(cmds, ListChangedCmd(m.SelectedList))
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

func (m Model) SetSize(size common.Size) Model {
	m.list.SetSize(size.Width, size.Height)
	return m
}

func (m *Model) SpaceChanged(id string) error {
	folders, err := m.ctx.Api.GetLists(id)
	if err != nil {
		return err
	}
	m.syncList(folders)

	return nil
}
