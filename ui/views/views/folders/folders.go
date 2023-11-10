package folders

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/ui/components/folders"
	"github.com/prgrs/clickup/ui/context"
)

type FoldersState uint

const (
	FoldersStateLoading FoldersState = iota
	FoldersStateList
)

type Model struct {
	ctx   *context.UserContext
	state FoldersState

	componentFoldersList folders.Model
}

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		ctx:                  ctx,
		componentFoldersList: folders.InitialModel(ctx),
		state:                FoldersStateList,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "esc":
			m.ctx.Logger.Info("Hiding folders view")
			cmds = append(cmds, HideFolderViewCmd())
		}
	}

	m.componentFoldersList, cmd = m.componentFoldersList.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.componentFoldersList.View()
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Info("Initializing view: Folders")
	return m.componentFoldersList.Init()
}
