package spaces

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

const WidgetId = "spacesList"

type Model struct {
	list              list.Model
	ctx               *context.UserContext
	log               *log.Logger
	WidgetId          common.WidgetId
	SelectedSpace     string
	SelectedWorkspace string
	spaces            []clickup.Space
}

func (m Model) KeyMap() help.KeyMap {
	return common.NewKeyMap(
		m.list.FullHelp,
		m.list.ShortHelp,
	).With(common.KeyBindingBack)
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	l := list.New([]list.Item{},
		list.NewDefaultDelegate(),
		0, 0)
	l.KeyMap.Quit.Unbind()
	l.SetShowHelp(false)

	log := logger.WithPrefix(logger.GetPrefix() + "/" + WidgetId)

	return Model{
		WidgetId:      WidgetId,
		list:          l,
		ctx:           ctx,
		SelectedSpace: "",
		spaces:        []clickup.Space{},
		log:           log,
	}
}

func (m *Model) syncList(spaces []clickup.Space) {
	m.log.Info("Synchronizing list...")
	m.spaces = spaces

	sre_index := 0 // TODO: rename
	items := spaceListToItems(spaces)
	itemsList := listitem.ItemListToBubblesItems(items)

	for i, item := range items {
		if item.Description() == m.ctx.Config.DefaultSpace {
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
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedSpace := listitem.BubblesItemToItem(m.list.SelectedItem()).Description()
			m.log.Infof("Selected space %s", selectedSpace)
			m.SelectedSpace = selectedSpace
			cmds = append(cmds, SpaceChangedCmd(selectedSpace))
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

func (m Model) SetSize(s common.Size) Model {
	m.list.SetSize(s.Width, s.Height)
	return m
}

func (m *Model) WorkspaceChanged(id string) error {
	m.log.Infof("Received: WorkspaceChangeMsg: %s", id)
	m.SelectedWorkspace = id

	spaces, err := m.ctx.Api.GetSpaces(id)
	if err != nil {
		return err
	}

	m.syncList(spaces)
	return nil
}
