package tasksidebar

import (
	"encoding/json"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

type Model struct {
	ctx      *context.UserContext
	viewport viewport.Model
}

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		ctx:      ctx,
		viewport: viewport.New(30, 30),
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

	case common.WindowSizeMsg:
		m.ctx.Logger.Info("TaskSidebar receive tea.WindowSizeMsg")

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 2

	case TaskSelectedMsg:
		id := string(msg)
		m.ctx.Logger.Infof("TaskSidebar receive TaskSelectedMsg: %s", id)

		task, err := m.getTask(id)
		if err != nil {
			return m, common.ErrCmd(err)
		}

		taskJson, err := json.MarshalIndent(task, "", "  ")
		if err != nil {
			return m, common.ErrCmd(err)
			return m, nil
		}

		m.viewport.SetContent(string(taskJson))
		m.viewport.GotoTop()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.viewport.View()
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
