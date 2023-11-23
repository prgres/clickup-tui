package folders

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
	"github.com/prgrs/clickup/ui/context"
)

const WidgetId = "foldersList"

type Model struct {
	WidgetId       common.WidgetId
	ctx            *context.UserContext
	list           list.Model
	folders        []clickup.Folder
	SelectedSpace  string
	SelectedFolder string
	log            *log.Logger
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	l := list.New([]list.Item{},
		list.NewDefaultDelegate(),
		0, 0)
	l.KeyMap.Quit.Unbind()

	log := logger.WithPrefix(logger.GetPrefix() + "/" + WidgetId)

	return Model{
		list:           l,
		ctx:            ctx,
		SelectedFolder: "",
		SelectedSpace:  ctx.Config.DefaultSpace,
		folders:        []clickup.Folder{},
		log:            log, // <3
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
	case FoldersListReloadedMsg:
		m.log.Info("Received: FoldersListReloadedMsg")
		m.syncList(msg)
		cmds = append(cmds, FoldersListReadyCmd())

	case tea.WindowSizeMsg:
		m.log.Debug("Received: tea.WindowSizeMsg",
			"width", msg.Width,
			"height", msg.Height)
		m.list.SetSize(msg.Width, msg.Height)

	case common.SpaceChangeMsg:
		m.log.Infof("Received: SpaceChangeMsg: %s", string(msg))
		m.SelectedSpace = string(msg)
		cmds = append(cmds, m.getFoldersCmd(m.SelectedSpace))

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
			cmds = append(cmds, common.FolderChangeCmd(selectedFolder))
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
