package taskssidebar

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/mattn/go-runewidth"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

const ComponentId = "widgetTaskSidebar"

type Model struct {
	ctx          *context.UserContext
	ComponentId  common.ComponentId
	log          *log.Logger
	SelectedTask clickup.Task
	viewport     viewport.Model
	size         common.Size
	Focused      bool
	Hidden       bool
	Ready        bool
	ifBorders    bool
}

func (m *Model) SetSize(s common.Size) {
	if m.ifBorders {
		s.Width -= 2  // two borders
		s.Height -= 2 // two borders
	}

	m.size = s
	m.viewport.Width = m.size.Width
	m.viewport.Height = m.size.Height

	task := lipgloss.NewStyle().Width(m.size.Width).
		Render(m.renderTask(m.SelectedTask))
	m.viewport.SetContent(task)
}

func (m Model) KeyMap() help.KeyMap {
	km := m.viewport.KeyMap

	return common.NewKeyMap(
		func() [][]key.Binding {
			return [][]key.Binding{
				{
					km.Down,
					km.Up,
				},
				{
					km.PageDown,
					km.PageUp,
				},
				{
					km.HalfPageUp,
					km.HalfPageDown,
				},
			}
		},
		func() []key.Binding {
			return []key.Binding{
				km.Down,
				km.Up,
				km.PageDown,
				km.PageUp,
			}
		},
	)
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	v := viewport.New(0, 0)
	v.Style = lipgloss.NewStyle().
		Height(0)
	v.SetContent("Loading...")

	size := common.Size{
		Width:  0,
		Height: 0,
	}

	log := logger.WithPrefix(logger.GetPrefix() + "/" + ComponentId)

	return Model{
		ctx:          ctx,
		viewport:     v,
		Focused:      false,
		Hidden:       false,
		SelectedTask: clickup.Task{},
		Ready:        false,
		log:          log,
		ifBorders:    true,
		size:         size,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) renderTask(task clickup.Task) string {
	s := strings.Builder{}

	header := fmt.Sprintf("[#%s] %s\n", task.Id, task.Name)
	s.WriteString(header)

	divider := strings.Repeat("-", runewidth.StringWidth(header))
	s.WriteString(divider)

	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(m.viewport.Width),
	)

	out, err := r.Render(task.MarkdownDescription)
	if err != nil {
		return err.Error()
	}
	s.WriteString(out)

	return s.String()
}

func (m Model) View() string {
	bColor := m.ctx.Theme.BordersColorInactive
	if m.Focused {
		bColor = m.ctx.Theme.BordersColorActive
	}

	styleBorders := m.ctx.Style.Borders.Copy().
		BorderForeground(bColor)

	return lipgloss.NewStyle().
		Inherit(styleBorders).
		Render(
			m.viewport.View(),
		)
}

func (m Model) Init() tea.Cmd {
	m.log.Info("Initializing...")
	return nil
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

func (m *Model) SetTask(id string) error {
	m.log.Infof("Received: TaskSelectedMsg: %s", id)
	m.Ready = false

	task, err := m.ctx.Api.GetTask(id)
	if err != nil {
		return err
	}

	m.SelectedTask = task
	m.viewport.SetContent(m.renderTask(task))

	_ = m.viewport.GotoTop()
	m.Ready = true

	return nil
}
