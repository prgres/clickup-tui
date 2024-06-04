package help

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

const id = "help"

type Model struct {
	id         common.Id
	ShowHelp   bool
	lastKey    string
	ctx        *context.UserContext
	log        *log.Logger
	help       help.Model
	inputStyle lipgloss.Style
	size       common.Size
	keyMap     KeyMap
}

func (m Model) Id() common.Id {
	return m.id
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	log := common.NewLogger(logger, common.ResourceTypeRegistry.WIDGET, id)

	return Model{
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
		id:         id,
		ctx:        ctx,
		log:        log,
		help:       help.New(),
		ShowHelp:   false,
		keyMap:     DefaultKeyMap(),
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeys(msg)
	}

	m.help, cmd = m.help.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m Model) View(keyMap help.KeyMap) string {
	availableWidth := m.size.Width

	status := " Waiting for input... "
	if m.lastKey != "" {
		status = " You chose: " + m.inputStyle.Render(m.lastKey) + " "
	}

	availableWidth -= lipgloss.Width(status)

	m.help.Width = availableWidth
	m.help.ShowAll = m.ShowHelp
	helpView := m.help.View(keyMap)

	availableWidth -= lipgloss.Width(helpView)

	dividerWidth := availableWidth
	if dividerWidth < 0 {
		dividerWidth = 0
	}

	divider := strings.Repeat(" ", dividerWidth)

	return lipgloss.JoinHorizontal(
		lipgloss.Bottom,
		helpView,
		divider,
		status,
	)
}

func (m Model) Init() tea.Cmd {
	m.log.Info("Initializing...")
	return nil
}

func (m *Model) SetSize(s common.Size) {
	m.size = s
	m.help.Width = s.Width
}
