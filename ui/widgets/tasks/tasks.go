package tasks

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	tabletasks "github.com/prgrs/clickup/ui/components/table-tasks"
	"github.com/prgrs/clickup/ui/context"
)

const WidgetId = "tasks"

type Model struct {
	log       *log.Logger
	ctx       *context.UserContext
	WidgetId  common.WidgetId
	size      common.Size
	Focused   bool
	Hidden    bool
	ifBorders bool

	state common.ComponentId

	spinner     spinner.Model
	showSpinner bool

	componenetTasksTable tabletasks.Model
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	log := logger.WithPrefix(logger.GetPrefix() + "/" + WidgetId)

	size := common.Size{
		Width:  0,
		Height: 0,
	}

	componenetTasksTable := tabletasks.InitialModel(ctx, log).SetFocused(true)

	return Model{
		WidgetId:  WidgetId,
		ctx:       ctx,
		size:      size,
		Focused:   false,
		Hidden:    false,
		log:       log,
		ifBorders: true,

		componenetTasksTable: componenetTasksTable,

		state: componenetTasksTable.ComponentId,

		spinner:     s,
		showSpinner: false,
	}
}

func (m *Model) SetSpinner(f bool) {
	m.showSpinner = f
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.componenetTasksTable, cmd = m.componenetTasksTable.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) SetTasks(tasks []clickup.Task) {
	m.componenetTasksTable.SetTasks(tasks)
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
		Width(m.size.Width - borderMargin).
		MaxWidth(m.size.Width + borderMargin).
		Height(m.size.Height - borderMargin).
		MaxHeight(m.size.Height + borderMargin)

	size := common.Size{
		Width:  m.size.Width - borderMargin,
		Height: m.size.Height - borderMargin,
	}

	if m.showSpinner {
		return style.Render(
			lipgloss.Place(
				size.Width, size.Height,
				lipgloss.Center,
				lipgloss.Center,
				fmt.Sprintf("%s Loading lists...", m.spinner.View()),
			),
		)
	}

	m.componenetTasksTable.SetSize(size)

	var content string
	switch m.state {
	case tabletasks.ComponentId:
		content = m.componenetTasksTable.View()
	default:
		content = "Unknown state"
	}

	return style.Render(content)
}

func (m Model) SetFocused(f bool) Model {
	m.Focused = f
	return m
}

func (m *Model) SetSize(s common.Size) {
	m.size = s
}

func (m *Model) Init() error {
	m.log.Info("Initializing...")
	return nil
}
