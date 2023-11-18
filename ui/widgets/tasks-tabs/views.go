package taskstabs

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

type Tab struct {
	Name   string
	Active bool
	Type   string
	Id     string
}

type Model struct {
	ctx          *context.UserContext
	tabs         map[string][]Tab
	SelectedTab  Tab
	SelectedList string
	Focused      bool
}

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		ctx:  ctx,
		tabs: map[string][]Tab{},
		// SelectedTab: nil,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		// case "h", "left":
		// 	m.SelectedView = prevView(m.views[m.SelectedList], m.SelectedView)
		// 	return m, ViewChangedCmd(m.SelectedView)
		// case "l", "right":
		// 	m.SelectedView = nextView(m.views[m.SelectedList], m.SelectedView)
		// 	return m, ViewChangedCmd(m.SelectedView)

		default:
			return m, nil
		}

	case common.ListChangeMsg:
		var cmds []tea.Cmd
		var tabs []Tab

		listName := string(msg)
		m.ctx.Logger.Infof("ViewsView received ListChangeMsg: %s", listName)

		tabList := Tab{
			Name:   listName,
			Type:   "list",
			Id:     listName,
			Active: true,
		}
		tabs = append(tabs, tabList)

		m.SelectedList = string(msg)
		m.SelectedTab = tabList

		views, err := m.ctx.Api.GetViewsFromList(string(msg))
		if err != nil {
			return m, common.ErrCmd(err)
		}

		if len(views) != 0 {
			m.ctx.Logger.Infof("ViewTabs: Views found for this list: %d", len(views))
			for _, view := range views {
				tabView := Tab{
					Name:   view.Name,
					Type:   "view",
					Id:     view.Id,
					Active: false,
				}
				tabs = append(tabs, tabView)
			}
			m.tabs[m.SelectedList] = tabs

		}
		cmds = append(cmds,
			FetchTasksForTabsCmd(tabs),
			TabChangedCmd(tabList),
		)

		return m, tea.Batch(
			cmds...,
		)

	case common.FocusMsg:
		m.ctx.Logger.Info("ViewsView received FocusMsg")
		return m, nil

	case TabChangedMsg:
		m.ctx.Logger.Info("ViewsView received TabChangedMsg")
		tab := Tab(msg)
		tabs := m.tabs[m.SelectedList]
		for _, t := range tabs {
			if t.Id == tab.Id {
				// m.SelectedViewStruct = view
				return m, tea.Batch(TabLoadedCmd(t))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := strings.Builder{}
	if len(m.tabs[m.SelectedList]) == 0 {
		return s.String()
	}

	s.WriteString(" Views |")
	if len(m.tabs[m.SelectedList]) == 0 {
		s.WriteString(" ")
	}

	for i, tab := range m.tabs[m.SelectedList] {
		t := ""
		if tab.Active {
			t = activeTabStyle.Render(tab.Name)
		} else {
			t = inactiveTabStyle.Render(tab.Name)
		}
		s.WriteString(" " + t + " ")

		if i != len(m.tabs[m.SelectedList])-1 {
			s.WriteString("|")
		}
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

type TabLoadedMsg Tab

func TabLoadedCmd(tab Tab) tea.Cmd {
	return func() tea.Msg {
		return TabLoadedMsg(tab)
	}
}
