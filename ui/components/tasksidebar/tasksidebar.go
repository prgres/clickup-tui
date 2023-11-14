package tasksidebar

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

type Model struct {
	ctx      *context.UserContext
	viewport viewport.Model
	Focused  bool
}

func InitialModel(ctx *context.UserContext) Model {
	v := viewport.New(0, 0)
	v.Style = lipgloss.NewStyle().
		Height(0)

	return Model{
		ctx:      ctx,
		viewport: v,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case InitMsg:
		m.ctx.Logger.Info("TaskSidebar receive InitMsg")
		m.viewport.SetContent("Loading...")

	case tea.WindowSizeMsg:
		m.ctx.Logger.Info("TaskSidebar receive tea.WindowSizeMsg")
		m.viewport.Width = int(0.6 * float32(m.ctx.WindowSize.Width))
		m.viewport.Height = int(0.7 * float32(m.ctx.WindowSize.Height))

	case TaskSelectedMsg:
		id := string(msg)
		m.ctx.Logger.Infof("TaskSidebar receive TaskSelectedMsg: %s", id)

		task, err := m.getTask(id)
		if err != nil {
			return m, common.ErrCmd(err)
		}

		m.viewport.SetContent(m.renderTask(task))

		_ = m.viewport.GotoTop()
	}

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
	bColor := lipgloss.Color("#FFF")
	if m.Focused {
		bColor = lipgloss.Color("#8909FF")
	}
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(bColor).
		BorderRight(true).
		BorderBottom(true).
		BorderTop(true).
		BorderLeft(true).
		Width(m.viewport.Width).
		Height(m.viewport.Height).
		Render(
			m.viewport.View(),
			// m.viewport.View() + "\n", // a extra newline due to table height do not consider headerView
		)
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Info("Initializing component: TaskSidebar")
	return InitCmd()
}

func (m Model) getTask(id string) (clickup.Task, error) {
	m.ctx.Logger.Infof("Getting task: %s", id)

	data, ok := m.ctx.Cache.Get("task", id)
	if ok {
		m.ctx.Logger.Infof("Task found in cache")
		var task clickup.Task
		if err := m.ctx.Cache.ParseData(data, &task); err != nil {
			return clickup.Task{}, err
		}

		return task, nil
	}
	m.ctx.Logger.Info("Task not found in cache")

	m.ctx.Logger.Info("Fetching task from API")
	client := m.ctx.Clickup

	task, err := client.GetTask(id)
	if err != nil {
		return clickup.Task{}, err
	}
	m.ctx.Logger.Infof("Found tasks %s", id)

	m.ctx.Logger.Info("Caching tasks")
	m.ctx.Cache.Set("task", id, task)

	return task, nil
}
