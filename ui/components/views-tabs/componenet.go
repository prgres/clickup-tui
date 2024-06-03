package viewstabs

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

const id = "tasks-tab"

type Tab struct {
	Name string
	Id   string
}

type Model struct {
	id        common.Id
	ctx       *context.UserContext
	log       *log.Logger
	keyMap    KeyMap
	tabs      []Tab
	size      common.Size
	Focused   bool
	Hidden    bool
	ifBorders bool
	Path      string
	StartIdx  int
	EndIdx    int

	SelectedIdx int
	Selected    string
}

func (m Model) Id() common.Id {
	return m.id
}

func (m *Model) SetSize(s common.Size) {
	m.size = s
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	log := common.NewLogger(logger, common.ResourceTypeRegistry.COMPONENT, id)

	return Model{
		id:          id,
		ctx:         ctx,
		tabs:        []Tab{},
		log:         log,
		keyMap:      DefaultKeyMap(),
		ifBorders:   true,
		Path:        "",
		StartIdx:    0,
		EndIdx:      0,
		SelectedIdx: 0,
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeys(msg)
	}

	return nil
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
		if m.Selected == tab.Id {
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
	m.SelectedIdx = 0
	m.tabs = tabs

	selectedTabId := ""
	if len(tabs) > 0 {
		selectedTabId = tabs[0].Id
	}
	m.Selected = selectedTabId
}

func (m Model) Size() common.Size {
	return m.size
}
