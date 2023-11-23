package taskstabs

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/ui/common"
	listitem "github.com/prgrs/clickup/ui/components/list-item"
	"github.com/prgrs/clickup/ui/context"
)

const WidgetId = "widgetTasksTabs"

type Tab struct {
	Name   string
	Active bool
	Type   string
	Id     string
}

type Model struct {
	ctx          *context.UserContext
	tabs         map[string][]Tab
	SelectedTab  int
	SelectedList string
	Focused      bool
	log          *log.Logger
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	log := logger.WithPrefix(logger.GetPrefix() + "/" + WidgetId)

	return Model{
		ctx:  ctx,
		tabs: map[string][]Tab{},
		log:  log,
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
			index := prevTab(m.tabs[m.SelectedList], m.SelectedTab)
			m.tabs[m.SelectedList][m.SelectedTab].Active = false
			m.SelectedTab = index
			m.tabs[m.SelectedList][index].Active = true
			return m, TabChangedCmd(m.tabs[m.SelectedList][index])
		case "l", "right":
			index := nextTab(m.tabs[m.SelectedList], m.SelectedTab)
			m.tabs[m.SelectedList][m.SelectedTab].Active = false
			m.SelectedTab = index
			m.tabs[m.SelectedList][index].Active = true
			return m, TabChangedCmd(m.tabs[m.SelectedList][index])

		default:
			return m, nil
		}

	case common.ListChangeMsg:
		var cmds []tea.Cmd
		var tabs []Tab

		list := listitem.Item(msg)
		m.log.Infof("Received: ListChangeMsg: %s", list.Description())

		tabList := Tab{
			Id:     list.Description(),
			Name:   list.Title(),
			Type:   "list",
			Active: true,
		}
		tabs = append(tabs, tabList)

		m.SelectedList = list.Description()
		m.SelectedTab = 0

		views, err := m.ctx.Api.GetViewsFromList(list.Description())
		if err != nil {
			return m, common.ErrCmd(err)
		}
		m.log.Infof("Found %d views found for the list", len(views))

		if len(views) != 0 {
			for _, view := range views {
				tabView := Tab{
					Name:   view.Name,
					Type:   "view",
					Id:     view.Id,
					Active: false,
				}
				tabs = append(tabs, tabView)
			}
		}

		m.tabs[m.SelectedList] = tabs

		cmds = append(cmds,
			FetchTasksForTabsCmd(tabs),
			TabChangedCmd(tabList),
		)

		return m, tea.Batch(cmds...)

	case common.FocusMsg:
		m.log.Info("Received: FocusMsg")
		return m, nil

	case TabChangedMsg:
		m.log.Info("Received: TabChangedMsg")
		tab := Tab(msg)
		tabs := m.tabs[m.SelectedList]
		for _, t := range tabs {
			if t.Id == tab.Id {
				return m, tea.Batch(TabLoadedCmd(t))
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := strings.Builder{}
	s.WriteString(" Views |")

	if len(m.tabs[m.SelectedList]) == 0 {
		s.WriteString(" ")
		return s.String()
	}
	m.log.Debugf("Rendering %d tabs", len(m.tabs[m.SelectedList]))

	for i, tab := range m.tabs[m.SelectedList] {
		m.log.Debugf("Rendering tab: %s %s", tab.Name, tab.Id)
		t := ""
		tabContent := " " + tab.Name + " "
		if tab.Active {
			t = activeTabStyle.Render(tabContent)
		} else {
			t = inactiveTabStyle.Render(tabContent)
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
	m.log.Info("Initializing...")
	return nil
}
