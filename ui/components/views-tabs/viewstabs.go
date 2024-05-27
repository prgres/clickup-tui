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

const id = "TasksTabs"

type Tab struct {
	Name string
	Id   string
}

type Model struct {
	id             common.Id
	ctx            *context.UserContext
	log            *log.Logger
	SelectedTab    string
	keyMap         KeyMap
	tabs           []Tab
	size           common.Size
	Focused        bool
	Hidden         bool
	ifBorders      bool
	Path           string
	StartIdx       int
	EndIdx         int
	SelectedTabIdx int
}

func (m Model) Id() common.Id {
	return m.id
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
	log := logger.WithPrefix(logger.GetPrefix() + "/componenet/" + id)

	return Model{
		id:             id,
		ctx:            ctx,
		tabs:           []Tab{},
		log:            log,
		keyMap:         DefaultKeyMap(),
		ifBorders:      true,
		Path:           "",
		StartIdx:       0,
		EndIdx:         0,
		SelectedTabIdx: 0,
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
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
			return TabChangedCmd(m.SelectedTab)

		case "L", "shift+right":
			index := nextTab(m.tabs, m.SelectedTabIdx)
			m.SelectedTabIdx = index
			m.SelectedTab = m.tabs[index].Id
			return TabChangedCmd(m.SelectedTab)

		case "h", "left":
			index := prevTab(m.tabs, m.SelectedTabIdx)
			m.SelectedTabIdx = index
			m.SelectedTab = m.tabs[index].Id
			return nil

		case "l", "right":
			index := nextTab(m.tabs, m.SelectedTabIdx)
			m.SelectedTabIdx = index
			m.SelectedTab = m.tabs[index].Id
			return nil

		case "enter":
			index := nextTab(m.tabs, m.SelectedTabIdx)
			m.SelectedTabIdx = index
			m.SelectedTab = m.tabs[index].Id
			return TabChangedCmd(m.SelectedTab)

		default:
			return nil
		}
	}

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	bColor := m.ctx.Theme.BordersColorInactive
	if m.Focused {
		bColor = m.ctx.Theme.BordersColorActive
	}

	borderMargin := 0
	if m.ifBorders {
		borderMargin = 2
	}

	styleBorders := m.ctx.Style.Borders.
		BorderForeground(bColor)

	style := lipgloss.NewStyle().
		Inherit(styleBorders).
		Height(1).
		MaxHeight(1 + borderMargin).
		Width(m.size.Width - borderMargin).
		MaxWidth(m.size.Width + borderMargin)

	var s []string

	moreTabsIcon := " + "
	tabPrefix := " Views |"
	suffix := ""

	if m.Path != "" {
		suffix += " | " + m.Path
	}

	for _, tab := range m.tabs {
		t := ""
		tabContent := " " + tab.Name + " "

		style := inactiveTabStyle
		if m.SelectedTab == tab.Id {
			style = activeTabStyle
		}
		t = style.Render(tabContent)

		content := " " + t + " "

		if lipgloss.Width(tabPrefix+strings.Join(s, "")+content+moreTabsIcon+suffix) >= m.size.Width-borderMargin {
			s = append(s, moreTabsIcon)
			break
		}
		s = append(s, content)
		s = append(s, "|")
	}
	content := strings.Join(s, "")

	dividerWidth := m.size.Width - borderMargin - lipgloss.Width(tabPrefix+content+moreTabsIcon+suffix)
	if dividerWidth < 0 {
		dividerWidth = 0
	}
	divider := strings.Repeat(" ", dividerWidth)

	return style.Render(
		tabPrefix + content + divider + suffix,
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

func (m *Model) SetFocused(f bool) {
	m.Focused = f
}

func (m Model) GetHidden() bool {
	return m.Hidden
}

func (m *Model) SetHidden(h bool) {
	m.Hidden = h
}

func (m *Model) SetTabs(tabs []Tab) {
	m.SelectedTabIdx = 0
	m.tabs = tabs

	selectedTabId := ""
	if len(tabs) > 0 {
		selectedTabId = tabs[0].Id
	}
	m.SelectedTab = selectedTabId
}

func (m Model) Size() common.Size {
	return m.size
}
