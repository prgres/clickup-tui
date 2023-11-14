package viewtabs

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

type Model struct {
	ctx                *context.UserContext
	views              map[string][]clickup.View
	SelectedView       string
	SelectedViewStruct clickup.View
	SelectedList       string
	Focused            bool
}

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		ctx:          ctx,
		views:        map[string][]clickup.View{},
		SelectedView: "",
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "h", "left":
			m.SelectedView = prevView(m.views[m.SelectedList], m.SelectedView)
			return m, ViewChangedCmd(m.SelectedView)
		case "l", "right":
			m.SelectedView = nextView(m.views[m.SelectedList], m.SelectedView)
			return m, ViewChangedCmd(m.SelectedView)

		default:
			return m, nil
		}

	case common.ListChangeMsg:
		m.ctx.Logger.Infof("ViewsView received ListChangeMsg: %s", string(msg))
		m.SelectedList = string(msg)

		views, err := m.ctx.Api.GetViewsFromList(string(msg))
		if err != nil {
			return m, common.ErrCmd(err)
		}

		if len(views) == 0 {
			return m, tea.Batch(
				ViewChangedCmd(m.SelectedView),
			)
		}

		m.views[m.SelectedList] = views
		m.SelectedView = m.views[m.SelectedList][0].Id
		viewsIds := []string{}
		for _, view := range views {
			if view.Id == m.SelectedView {
				continue
			}
			viewsIds = append(viewsIds, view.Id)
		}

		return m, tea.Batch(
			ViewChangedCmd(m.SelectedView),
			FetchViewsCmd(viewsToIdList(views)),
		)

	case common.FocusMsg:
		m.ctx.Logger.Info("ViewsView received FocusMsg")
		return m, nil

	case ViewsListLoadedMsg:
		m.ctx.Logger.Info("ViewsView received ViewsListLoadedMsg")
		m.views[m.SelectedList] = []clickup.View(msg)
		return m, nil

	case ViewChangedMsg:
		m.ctx.Logger.Info("ViewsView received ViewChangedMsg")
		viewName := string(msg)
		allViews := m.views[m.SelectedList]
		for _, view := range allViews {
			if view.Id == viewName {
				m.SelectedViewStruct = view
				return m, tea.Batch(common.ViewLoadedCmd(view))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := strings.Builder{}
	if len(m.views[m.SelectedList]) == 0 {
		return s.String()
	}

	s.WriteString("Views")

	for _, view := range m.views[m.SelectedList] {
		t := ""
		if m.SelectedView == view.Id {
			t = activeTabStyle.Render(view.Name)
		} else {
			t = inactiveTabStyle.Render(view.Name)
		}
		s.WriteString(" | ")
		s.WriteString(t)
	}

	bColor := lipgloss.Color("#FFF")
	if m.Focused {
		bColor = lipgloss.Color("#8909FF")
	}
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(bColor).
		BorderRight(true).
		BorderBottom(true).
		BorderTop(true).
		BorderLeft(true).
		Render(
			s.String(),
		)
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Info("Initializing component: TabsView")
	return nil
}
