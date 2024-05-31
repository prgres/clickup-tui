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

const id = "task-sidebar"

type Model struct {
	ctx          *context.UserContext
	id           common.Id
	log          *log.Logger
	SelectedTask clickup.Task
	viewport     viewport.Model
	size         common.Size
	Focused      bool
	Hidden       bool
	Ready        bool
	ifBorders    bool
	keyMap       KeyMap
}

func (m Model) Id() common.Id {
	return m.id
}

type KeyMap struct {
	viewport.KeyMap
}

func (m Model) KeyMap() KeyMap {
	return m.keyMap
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		KeyMap: viewport.DefaultKeyMap(),
	}
}

func (m *Model) SetSize(s common.Size) {
	if m.ifBorders {
		s.Width -= 2  // two borders
		s.Height -= 2 // two borders
	}

	m.size = s
	m.viewport.Width = m.size.Width
	m.viewport.Height = m.size.Height
}

func (m Model) Help() help.KeyMap {
	km := m.keyMap

	return common.NewHelp(
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

	log := logger.WithPrefix(logger.GetPrefix() + "/component/" + id)

	return Model{
		id:           id,
		ctx:          ctx,
		viewport:     v,
		Focused:      false,
		Hidden:       false,
		SelectedTask: clickup.Task{},
		Ready:        false,
		log:          log,
		ifBorders:    true,
		size:         size,
		keyMap:       DefaultKeyMap(),
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m Model) renderTask(task clickup.Task) (string, error) {
	s := strings.Builder{}

	header := fmt.Sprintf("[#%s] %s\n", task.Id, task.Name)
	s.WriteString(header)

	divider := strings.Repeat("-", runewidth.StringWidth(header))
	s.WriteString(divider)

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(m.viewport.Width),
	)
	if err != nil {
		return "", err
	}

	out, err := r.Render(task.MarkdownDescription)
	if err != nil {
		return "", err
	}
	s.WriteString(out)

	return s.String(), nil
}

func (m Model) View() string {
	bColor := m.ctx.Theme.BordersColorInactive
	if m.Focused {
		bColor = m.ctx.Theme.BordersColorActive
	}

	styleBorders := m.ctx.Style.Borders.
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

func (m Model) WithFocused(f bool) Model {
	m.Focused = f
	return m
}

func (m *Model) SetFocused(f bool) *Model {
	m.Focused = f
	return m
}

func (m Model) GetHidden() bool {
	return m.Hidden
}

func (m *Model) SetHidden(h bool) *Model {
	m.Hidden = h
	return m
}

func (m Model) WithHidden(h bool) Model {
	m.Hidden = h
	return m
}

func (m *Model) SelectTask(id string) error {
	m.Ready = false

	task, err := m.ctx.Api.GetTask(id)
	if err != nil {
		return err
	}

	if err := m.SetTask(task); err != nil {
		return err
	}
	m.Ready = true

	return nil
}

func (m *Model) SetTask(task clickup.Task) error {
	m.SelectedTask = task
	renderedTask, err := m.renderTask(task)
	if err != nil {
		return err
	}

	m.viewport.SetContent(renderedTask)
	_ = m.viewport.GotoTop()

	return nil
}
