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

func folderListToItems(folders []clickup.Folder) []item {
	items := make([]item, len(folders))
	for i, folder := range folders {
		items[i] = folderToItem(folder)
	}
	return items
}

func folderToItem(folder clickup.Folder) item {
	return item{
		folder.Name,
		folder.Id,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case FoldersListReloadedMsg:
		m.ctx.Logger.Info("FolderView received FoldersListReloadedMsg")
		m.syncList(msg)

	case tea.WindowSizeMsg:
		m.ctx.Logger.Info("FolderView received tea.WindowSizeMsg")
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

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
	return common.SpaceChangeCmd(SPACE_SRE)
}

func (m Model) getFoldersCmd(space string) tea.Cmd {
	return func() tea.Msg {
		folders, err := m.getFolders(space)
		if err != nil {
			return common.ErrMsg(err)
		}

		return FoldersListReloadedMsg(folders)
	}
}

func (m Model) getFolders(space string) ([]clickup.Folder, error) {
	m.ctx.Logger.Infof("Getting folders for space: %s", space)
	client := m.ctx.Clickup

	data, ok := m.ctx.Cache.Get("folders", space)
	if ok {
		m.ctx.Logger.Infof("Folders found in cache")
		var folders []clickup.Folder
		if err := m.ctx.Cache.ParseData(data, &folders); err != nil {
			return nil, err
		}

		return folders, nil
	}
	m.ctx.Logger.Infof("Folders not found in cache")

	m.ctx.Logger.Infof("Fetching folders from API")
	folders, err := client.GetFolders(space)
	if err != nil {
		return nil, err
	}
	m.ctx.Logger.Infof("Found %d folders for space: %s", len(folders), space)

	m.ctx.Logger.Infof("Caching folders")
	m.ctx.Cache.Set("folders", space, folders)

	return folders, nil
}
