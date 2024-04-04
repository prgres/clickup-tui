package viewstabs

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
	ctx            *context.UserContext
	log            *log.Logger
	SelectedTab    string
	keyMap         KeyMap
	tabs           []Tab
	size           common.Size
	SelectedTabIdx int
	Focused        bool
	Hidden         bool
	ifBorders      bool
}

func (m *Model) SetSize(s common.Size) {
	m.size = s
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
		ctx:       ctx,
		tabs:      []Tab{},
		log:       log,
		keyMap:    DefaultKeyMap(),
		ifBorders: true,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "H", "shift+left":
			index := prevTab(m.tabs, m.SelectedTabIdx)
			m.SelectedTabIdx = index
			m.SelectedTab = m.tabs[index].Id
			return m, TabChangedCmd(m.SelectedTab)

		case "L", "shift+right":
			index := nextTab(m.tabs, m.SelectedTabIdx)
			m.SelectedTabIdx = index
			m.SelectedTab = m.tabs[index].Id
			return m, TabChangedCmd(m.SelectedTab)

		case "h", "left":
			index := prevTab(m.tabs, m.SelectedTabIdx)
			m.SelectedTabIdx = index
			m.SelectedTab = m.tabs[index].Id
			return m, nil

		case "l", "right":
			index := nextTab(m.tabs, m.SelectedTabIdx)
			m.SelectedTabIdx = index
			m.SelectedTab = m.tabs[index].Id
			return m, nil

		case "enter":
			return m, TabChangedCmd(m.SelectedTab)

		default:
			return m, nil
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	bColor := lipgloss.Color("#FFF")
	if m.Focused {
		bColor = lipgloss.Color("#8909FF")
	}

	borderMargin := 0
	if m.ifBorders {
		borderMargin = 2
	}

	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(bColor).
		BorderBottom(m.ifBorders).
		BorderRight(m.ifBorders).
		BorderTop(m.ifBorders).
		BorderLeft(m.ifBorders).
		Height(1).
		MaxHeight(1 + borderMargin).
		Width(m.size.Width - borderMargin).
		MaxWidth(m.size.Width + borderMargin)

	s := strings.Builder{}
	s.WriteString(" Views |")

	if len(m.tabs) == 0 {
		s.WriteString(" ")
		return style.Render(s.String())
	}
	m.log.Debugf("Rendering %d tabs", len(m.tabs))

	moreTabsIcon := " + "
	for i, tab := range m.tabs {
		m.log.Debugf("Rendering tab: %s %s", tab.Name, tab.Id)
		t := ""
		tabContent := " " + tab.Name + " "
		if m.SelectedTab == tab.Id {
			t = activeTabStyle.Render(tabContent)
		} else {
			t = inactiveTabStyle.Render(tabContent)
		}

		content := " " + t + " "

		if lipgloss.Width(s.String()+content+moreTabsIcon) >= m.size.Width-borderMargin {
			s.WriteString(moreTabsIcon)
			break
		}
		s.WriteString(content)

		if i != len(m.tabs)-1 {
			s.WriteString("|")
		}

	}

	return style.Render(s.String())
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

func (m Model) SetFocused(f bool) Model {
	m.Focused = f
	return m
}

func (m Model) GetHidden() bool {
	return m.Hidden
}

func (m Model) SetHidden(h bool) Model {
	m.Hidden = h
	return m
}

func (m *Model) SetTabs(tabs []Tab) {
	m.SelectedTabIdx = 0
	m.tabs = tabs
	if len(tabs) > 0 {
		m.SelectedTab = tabs[0].Id
	}
}
