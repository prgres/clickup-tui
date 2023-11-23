package help

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
	"golang.org/x/term"
)

const WidgetId = "widgetHelp"

type Model struct {
	WidgetId common.WidgetId
	ShowHelp bool

	lastKey    string
	ctx        *context.UserContext
	log        *log.Logger
	help       help.Model
	inputStyle lipgloss.Style
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	log := logger.WithPrefix(logger.GetPrefix() + "/" + WidgetId)

	return Model{
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),

		// keys:     keys,
		WidgetId: WidgetId,
		ctx:      ctx,
		log:      log,
		help:     help.New(),
		ShowHelp: false,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.log.Info("Received: tea.WindowSizeMsg",
			"width", msg.Width,
			"height", msg.Height)
		m.help.Width = msg.Width
	case tea.KeyMsg:
		m.lastKey = msg.String()
		switch keypress := msg.String(); keypress {
		// case key.Matches(msg, m.keys.Up):
		// 	m.lastKey = "↑"
		// case key.Matches(msg, m.keys.Down):
		// 	m.lastKey = "↓"
		// case key.Matches(msg, m.keys.Left):
		// 	m.lastKey = "←"
		// case key.Matches(msg, m.keys.Right):
		// 	m.lastKey = "→"
		case "?": // key.Matches(msg, m.keys.Help):
			m.ShowHelp = !m.ShowHelp
			m.help.ShowAll = !m.help.ShowAll
			// }
			// case tea.KeyMsg:
			// 	switch keypress := msg.String(); keypress {
			// 	case "enter":
		}
	}

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View(keyMap help.KeyMap) string {
	var status string
	if m.lastKey == "" {
		status = "Waiting for input..."
	} else {
		status = "You chose: " + m.inputStyle.Render(m.lastKey)
	}
	m.help.ShowAll = m.ShowHelp
	helpView := m.help.View(keyMap)

	physicalWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	dividerWidth := physicalWidth - lipgloss.Width(helpView) - lipgloss.Width(status)

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
