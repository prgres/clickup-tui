package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

type ViewsListLoadedMsg []clickup.View

type ViewChangedMsg string

func ViewChangedCmd(view string) tea.Cmd {
	return func() tea.Msg {
		return ViewChangedMsg(view)
	}
}

type SpaceChangedMsg string

func SpaceChangedCmd(space string) tea.Cmd {
	return func() tea.Msg {
		return SpaceChangedMsg(space)
	}
}

type Model struct {
	ctx           *context.UserContext
	views         map[string][]clickup.View
	SelectedView  string
	SelectedSpace string
}

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		ctx:           ctx,
		views:         map[string][]clickup.View{},
		SelectedView:  SPACE_SRE_LIST_COOL,
		SelectedSpace: SPACE_SRE,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case SpaceChangedMsg:
		m.ctx.Logger.Infof("ViewsView received SpaceChangedMsg")
		m.SelectedSpace = string(msg)
		views,err := m.getViews(string(msg))
		if err != nil {
		  return m, nil // TODO: handle error
		}
		m.views[m.SelectedSpace]= views
		m.SelectedView = m.views[m.SelectedSpace][0].Id

		return m, tea.Batch(
			ViewChangedCmd(m.SelectedView),
		)

	case ViewsListLoadedMsg:
		m.ctx.Logger.Infof("ViewsView received ViewsListLoadedMsg")
		m.views[m.SelectedSpace] = []clickup.View(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	viewsNames := []string{}
	for _, view := range m.views[m.SelectedSpace] {
		viewsNames = append(viewsNames, view.Name)
	}

	return strings.Join(viewsNames, " | ")
}

func (m Model) Init() tea.Msg {
	return SpaceChangedMsg(SPACE_SRE)
}

func (m Model) getViewsCmd(space string) tea.Cmd {
	return func() tea.Msg {
		views, err := m.getViews(space)
		if err != nil {
			return common.ErrMsg(err)
		}
		return ViewsListLoadedMsg(views)
	}
}

func (m Model) getViews(space string) ([]clickup.View, error) {
	m.ctx.Logger.Infof("Getting views for space: %s", space)
	if m.views[space] != nil {
		m.ctx.Logger.Info("Tasks found in cache")
		return m.views[space], nil
	}

	m.ctx.Logger.Info("Fetching tasks from API")
	client := m.ctx.Clickup

	m.ctx.Logger.Infof("Getting views from space: %s", space)
	views, err := client.GetViewsFromSpace(space)
	if err != nil {
		return nil, err
	}
	m.ctx.Logger.Infof("Found %d views in space %s", len(views), space)

	return views, nil
}
