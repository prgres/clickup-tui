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

	lastKey  string
	ctx      *context.UserContext
	log      *log.Logger
	help     help.Model
	ShowHelp bool
	// keys    keyMap

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

const (
	// In real life situations we'd adjust the document to fit the width we've
	// detected. In the case of this example we're hardcoding the width, and
	// later using the detected width only to truncate in order to avoid jaggy
	// wrapping.
	width = 96

	columnWidth = 30
)

func (m Model) View(keyMap help.KeyMap) string {

	// doc := strings.Builder{}
	// {
	// 	okButton := activeButtonStyle.Render("Yes")
	// 	cancelButton := buttonStyle.Render("Maybe")

	// 	question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("Are you sure you want to eat marmalade?")
	// 	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	// 	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	// 	dialog := lipgloss.Place(width, 9,
	// 		lipgloss.Center, lipgloss.Center,
	// 		dialogBoxStyle.Render(ui),
	// 		lipgloss.WithWhitespaceChars("猫咪"),
	// 		lipgloss.WithWhitespaceForeground(subtle),
	// 	)

	// 	doc.WriteString(dialog + "\n\n")
	// }

	// if physicalWidth > 0 {
	// docStyle = docStyle.MaxWidth(physicalWidth)
	// }
	// return docStyle.Render(doc.String())

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

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
// type keyMap struct {
// 	Up    key.Binding
// 	Down  key.Binding
// 	Left  key.Binding
// 	Right key.Binding
// 	Help  key.Binding
// 	Quit  key.Binding
// }

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
// func (k keyMap) ShortHelp() []key.Binding {
// 	return []key.Binding{k.Help, k.Quit}
// }

// // FullHelp returns keybindings for the expanded help view. It's part of the
// // key.Map interface.
// func (k keyMap) FullHelp() [][]key.Binding {
// 	return [][]key.Binding{
// 		{k.Up, k.Down, k.Left, k.Right}, // first column
// 		{k.Help, k.Quit},                // second column
// 	}
// }

// var keys = keyMap{
// 	Up: key.NewBinding(
// 		key.WithKeys("up", "k"),
// 		key.WithHelp("↑/k", "move up"),
// 	),
// 	Down: key.NewBinding(
// 		key.WithKeys("down", "j"),
// 		key.WithHelp("↓/j", "move down"),
// 	),
// 	Left: key.NewBinding(
// 		key.WithKeys("left", "h"),
// 		key.WithHelp("←/h", "move left"),
// 	),
// 	Right: key.NewBinding(
// 		key.WithKeys("right", "l"),
// 		key.WithHelp("→/l", "move right"),
// 	),
// 	Help: key.NewBinding(
// 		key.WithKeys("?"),
// 		key.WithHelp("?", "toggle help"),
// 	),
// 	Quit: key.NewBinding(
// 		key.WithKeys("q", "esc", "ctrl+c"),
// 		key.WithHelp("q", "quit"),
// 	),
// }
