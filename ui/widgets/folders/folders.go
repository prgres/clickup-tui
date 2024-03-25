package folders

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

const WidgetId = "widgetFoldersList"

type Model struct {
	list           list.Model
	ctx            *context.UserContext
	log            *log.Logger
	WidgetId       common.WidgetId
	SelectedSpace  string
	SelectedFolder string
	folders        []clickup.Folder
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
		WidgetId:       WidgetId,
		list:           l,
		ctx:            ctx,
		SelectedFolder: "",
		SelectedSpace:  ctx.Config.DefaultSpace,
		folders:        []clickup.Folder{},
		log:            log,
	}
}

func (m *Model) syncList(folders []clickup.Folder) {
	m.log.Info("Synchronizing list")
	m.folders = folders

	sre_index := 0
	items := folderListToItems(folders)
	itemsList := listitem.ItemListToBubblesItems(items)

	for i, item := range items {
		if item.Description() == m.ctx.Config.DefaultFolder {
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
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedFolder := listitem.BubblesItemToItem(m.list.SelectedItem()).Description()
			m.log.Infof("Selected folder %s", selectedFolder)
			m.SelectedFolder = selectedFolder
			cmds = append(cmds, FolderChangeCmd(selectedFolder))
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
	m.log.Info("Initializing...")
	return nil
}

func (m Model) SetSize(s common.Size) Model {
	m.list.SetSize(s.Width, s.Height)
	return m
}

func (m *Model) SpaceChange(id string) error {
	m.log.Infof("Received: SpaceChangeMsg: %s", id)
	m.SelectedSpace = id

	folders, err := m.ctx.Api.GetFolders(id)
	if err != nil {
		return err
	}

	m.syncList(folders)
	return nil
}
