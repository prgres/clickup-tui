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
	ctx         *context.UserContext
	viewCompact common.View
	log         *log.Logger
	dialogHelp  help.Model
	keyMap      KeyMap
}

type KeyMap struct {
	Refresh   key.Binding
	ForceQuit key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Refresh: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "go to refresh"),
		),
		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("ctrl+c/q", "quit"),
		),
	}
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	log := logger.WithPrefix("UI")

	return Model{
		ctx:         ctx,
		log:         log,
		viewCompact: compact.InitialModel(ctx, log),
		dialogHelp:  help.InitialModel(ctx, log),
		keyMap:      DefaultKeyMap(),
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case common.ErrMsg:
		m.log.Fatal(msg.Error())
		return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.ForceQuit):
			return m, tea.Quit

		case key.Matches(msg, m.keyMap.Refresh):
			m.log.Info("Refreshing...")
			if err := m.ctx.Api.InvalidateCache(); err != nil {
				m.log.Error("Failed to invalidate cache", "error", err)
			}
			m.log.Debug("Cache invalidated")
		}

	case tea.WindowSizeMsg:
		m.log.Debug(
			"Received: tea.WindowSizeMsg",
			"width", msg.Width,
			"height", msg.Height)

		m.ctx.WindowSize.Width = msg.Width
		m.ctx.WindowSize.Height = msg.Height
		m.ctx.WindowSize.Height = msg.Height
		m.ctx.WindowSize.Height = msg.Height
		m.ctx.WindowSize.Height = msg.Height
	}

	m.viewCompact, cmd = m.viewCompact.Update(msg)
	cmds = append(cmds, cmd)

	m.dialogHelp, cmd = m.dialogHelp.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var viewToRender common.View

	viewToRender = m.viewCompact

	viewKm := viewToRender.KeyMap()

	km := common.NewKeyMap(
		func() [][]key.Binding {
			return append(
				viewKm.FullHelp(),
				[][]key.Binding{
					{
						m.keyMap.Refresh,
					},
				}...)
		},
		viewKm.ShortHelp,
	)

	footer := m.dialogHelp.View(km)
	footerHeight := lipgloss.Height(footer)

	physicalHeight := m.ctx.WindowSize.Height
	physicalWidth := m.ctx.WindowSize.Width

	viewHeight := physicalHeight - footerHeight
	viewToRender = viewToRender.SetSize(common.Size{
		Width:  physicalWidth,
		Height: viewHeight - m.ctx.WindowSize.MetaHeight,
	})

	dividerHeight := physicalHeight - viewHeight - footerHeight

	if dividerHeight < 0 {
		dividerHeight = 0
		m.log.Info("dividerHeight", "dividerHeight", dividerHeight)
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
