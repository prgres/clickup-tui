package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
	"github.com/prgrs/clickup/ui/views/compact"
	"github.com/prgrs/clickup/ui/widgets/help"
)

type Model struct {
	ctx    *context.UserContext
	log    *log.Logger
	keyMap KeyMap

	viewCompact *compact.Model
	dialogHelp  *help.Model
}

type KeyMap struct {
	ForceQuit key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("ctrl+c/q", "quit"),
		),
	}
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	log := logger.WithPrefix("UI")

	var (
		viewCompact = compact.InitialModel(ctx, log)
		dialogHelp  = help.InitialModel(ctx, log)
	)

	return Model{
		ctx:    ctx,
		log:    log,
		keyMap: DefaultKeyMap(),

		dialogHelp:  &dialogHelp,
		viewCompact: &viewCompact,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case common.ErrMsg:
		m.log.Error(msg.Error())
		return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.ForceQuit):
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.log.Debug(
			"Received: tea.WindowSizeMsg",
			"width", msg.Width,
			"height", msg.Height)
		m.ctx.WindowSize.Set(msg.Width, msg.Height)
		return m, tea.Batch(cmds...)
	}

	cmds = append(cmds,
		m.viewCompact.Update(msg),
		m.dialogHelp.Update(msg),
	)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var viewToRender common.UIElement = m.viewCompact

	viewKm := viewToRender.Help()

	km := common.NewHelp(
		viewKm.FullHelp,
		viewKm.ShortHelp,
	)

	physicalHeight := m.ctx.WindowSize.Height
	physicalWidth := m.ctx.WindowSize.Width

	m.dialogHelp.SetSize(common.Size{
		Width:  physicalWidth,
		Height: physicalHeight,
	})

	footer := m.dialogHelp.View(km)
	footerHeight := lipgloss.Height(footer)

	viewHeight := physicalHeight - footerHeight
	viewToRender.SetSize(common.Size{
		Width:  physicalWidth,
		Height: viewHeight - m.ctx.WindowSize.MetaHeight,
	})

	dividerHeight := physicalHeight - viewHeight - footerHeight

	if dividerHeight < 0 {
		dividerHeight = 0
	}

	divider := strings.Repeat("\n", dividerHeight)

	m.ctx.WindowSize.MetaHeight = lipgloss.Height(divider) + footerHeight

	return lipgloss.JoinVertical(
		lipgloss.Left,
		viewToRender.View(),
		divider,
		footer,
	)
}

func (m Model) Init() tea.Cmd {
	m.log.Info("Initializing...")
	return tea.Batch(
		m.viewCompact.Init(),
		m.dialogHelp.Init(),
	)
}
