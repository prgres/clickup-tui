package folderslist

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

const id = "folders-list"

type Model struct {
	id             common.Id
	list           list.Model
	ctx            *context.UserContext
	log            *log.Logger
	SelectedFolder clickup.Folder
	folders        []clickup.Folder
}

func (m Model) Id() common.Id {
	return m.id
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
	l.Title = "Folders"

	log := logger.WithPrefix(logger.GetPrefix() + "/component/" + id)

	return Model{
		id:             id,
		list:           l,
		ctx:            ctx,
		SelectedFolder: clickup.Folder{},
		folders:        []clickup.Folder{},
		log:            log,
	}
}

func (m *Model) syncList(folders []clickup.Folder) {
	m.log.Info("Synchronizing list")
	m.folders = folders

	items := NewListItem(folders)

	for _, item := range items {
		i := item.(listitem.Item)
		if i.Title() == m.ctx.Config.DefaultFolder {
			m.SelectedFolder = i.Data().(clickup.Folder)
		}
	}

	m.list.SetItems(items)
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
			selectedFolder := m.list.SelectedItem().(listitem.Item).Data().(clickup.Folder)
			m.log.Info("Selected folder", "id", selectedFolder.Id, "name", selectedFolder.Name)
			m.SelectedFolder = selectedFolder
			return m, FolderChangeCmd(selectedFolder.Id)

		case "J", "shift+down":
			m.list.CursorDown()
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedFolder := m.list.SelectedItem().(listitem.Item).Data().(clickup.Folder)
			m.log.Info("Selected folder", "id", selectedFolder.Id, "name", selectedFolder.Name)
			m.SelectedFolder = selectedFolder
			return m, common.FolderPreviewCmd(selectedFolder.Id)

		case "K", "shift+up":
			m.list.CursorUp()
			if m.list.SelectedItem() == nil {
				m.log.Info("List is empty")
				break
			}
			selectedFolder := m.list.SelectedItem().(listitem.Item).Data().(clickup.Folder)
			m.log.Info("Selected folder", "id", selectedFolder.Id, "name", selectedFolder.Name)
			m.SelectedFolder = selectedFolder
			return m, common.FolderPreviewCmd(selectedFolder.Id)
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

func (m *Model) SetSize(s common.Size) {
	m.list.SetSize(s.Width, s.Height)
}

func (m *Model) SpaceChanged(id string) error {
	m.log.Infof("Received: SpaceChangedMsg: %s", id)

	folders, err := m.ctx.Api.GetFolders(id)
	if err != nil {
		return err
	}

	m.syncList(folders)
	return nil
}

func NewListItem(items []clickup.Folder) []list.Item {
	result := make([]list.Item, len(items))
	for i, v := range items {
		result[i] = listitem.NewItem(v.Name, v.Id, v)
	}
	return result
}
