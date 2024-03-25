package taskstabs

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

const WidgetId = "widgetTasksTabs"

type Tab struct {
	Name   string
	Type   string
	Id     string
	Active bool
}

type Model struct {
	ctx         *context.UserContext
	tabs        []Tab
	log         *log.Logger
	keyMap      KeyMap
	size        common.Size
	SelectedTab int
	Focused     bool
	Hidden      bool
}

func (m Model) SetSize(s common.Size) common.Widget {
	m.size = s

	return m
}

func (m Model) KeyMap() help.KeyMap {
	return common.NewKeyMap(
		func() [][]key.Binding {
			return [][]key.Binding{
				{
					m.keyMap.CursorLeft,
					m.keyMap.CursorRight,
					m.keyMap.Select,
				},
			}
		},
		func() []key.Binding {
			return []key.Binding{
				m.keyMap.CursorLeft,
				m.keyMap.CursorRight,
				m.keyMap.Select,
			}
		},
	)
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	log := logger.WithPrefix(logger.GetPrefix() + "/" + WidgetId)

	return Model{
		ctx:    ctx,
		tabs:   []Tab{},
		log:    log,
		keyMap: DefaultKeyMap(),
	}
}

func (m Model) Update(msg tea.Msg) (common.Widget, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "h", "left":
			index := prevTab(m.tabs, m.SelectedTab)
			m.tabs[m.SelectedTab].Active = false
			m.SelectedTab = index
			m.tabs[index].Active = true
			return m, TabChangedCmd(m.tabs[index])
		case "l", "right":
			index := nextTab(m.tabs, m.SelectedTab)
			m.tabs[m.SelectedTab].Active = false
			m.SelectedTab = index
			m.tabs[index].Active = true
			return m, TabChangedCmd(m.tabs[index])

		default:
			return m, nil
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	s := strings.Builder{}
	s.WriteString(" Views |")

	if len(m.tabs) == 0 {
		s.WriteString(" ")
		return s.String()
	}
	m.log.Debugf("Rendering %d tabs", len(m.tabs))

	for i, tab := range m.tabs {
		m.log.Debugf("Rendering tab: %s %s", tab.Name, tab.Id)
		t := ""
		tabContent := " " + tab.Name + " "
		if tab.Active {
			t = activeTabStyle.Render(tabContent)
		} else {
			t = inactiveTabStyle.Render(tabContent)
		}
		s.WriteString(" " + t + " ")

		if i != len(m.tabs)-1 {
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

type KeyMap struct {
	CursorLeft         key.Binding
	CursorRight        key.Binding
	Select             key.Binding
	SwitchFocusToTasks key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		CursorLeft: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("h, left", "previous tab"),
		),
		CursorRight: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("l, right", "next tab"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		SwitchFocusToTasks: key.NewBinding(
			key.WithKeys("j", "k", "escape"),
			key.WithHelp("j/k/escape", "switch focus to tasks table"),
		),
	}
}

func (m Model) GetFocused() bool {
	return m.Focused
}

func (m Model) SetFocused(f bool) common.Widget {
	m.Focused = f
	return m
}

func (m Model) GetHidden() bool {
	return m.Hidden
}

func (m Model) SetHidden(h bool) common.Widget {
	m.Hidden = h
	return m
}

func (m Model) ListChanged(id string) (common.Widget, tea.Cmd) {
	var cmds []tea.Cmd
	var tabs []Tab

	m.log.Infof("Received: ListChangedMsg: %s", id)
	list, err := m.ctx.Api.GetList(id)
	if err != nil {
		return m, common.ErrCmd(err)
	}

	tabList := Tab{
		Id:     list.Id,
		Name:   list.Name,
		Type:   "list",
		Active: true,
	}
	tabs = append(tabs, tabList)

	m.SelectedTab = 0

	views, err := m.ctx.Api.GetViewsFromList(list.Id)
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

	m.tabs = tabs

	cmds = append(cmds,
		FetchTasksForTabsCmd(tabs),
		TabChangedCmd(tabList),
	)

	return m, tea.Batch(cmds...)
}
