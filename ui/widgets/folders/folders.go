package folders

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
	folders        []clickup.Folder
	SelectedSpace  string
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
		SelectedSpace:  SPACE_SRE,
		folders:        []clickup.Folder{},
	}
}

func (m *Model) syncList(folders []clickup.Folder) {
	m.ctx.Logger.Info("Synchronizing list")
	m.folders = folders

	sre_index := 0
	items := folderListToItems(folders)
	itemsList := itemListToItems(items)

	for i, item := range items {
		if item.desc == SPACE_SRE {
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
		m.ctx.Logger.Info("FolderView received FoldersListReloadedMsg")
		m.syncList(msg)
		cmds = append(cmds, FoldersListReadyCmd())

	case tea.WindowSizeMsg:
		m.ctx.Logger.Info("FolderView received tea.WindowSizeMsg")
		m.list.SetSize(msg.Width, msg.Height)

	case common.SpaceChangeMsg:
		m.ctx.Logger.Infof("FolderView received SpaceChangeMsg: %s", string(msg))
		m.SelectedSpace = string(msg)
		cmds = append(cmds, m.getFoldersCmd(m.SelectedSpace))

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			selectedFolder := m.list.SelectedItem().(item).desc
			m.ctx.Logger.Infof("Selected folder %s", selectedFolder)
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
	m.ctx.Logger.Infof("Initializing component: foldersList")
	return nil
}
