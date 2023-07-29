package tasktable

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

const (
	SPACE_SRE_LIST_COOL = "q5kna-61288"
	SPACE_SRE           = "48458830"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func taskListToRows(tasks []clickup.Task, columns []table.Column) []table.Row {
	rows := make([]table.Row, len(tasks))
	for i, task := range tasks {
		rows[i] = taskToRow(task, columns)
	}
	return rows
}

func taskToRow(task clickup.Task, columns []table.Column) table.Row {
	values := table.Row{}
	for _, column := range columns {
		switch column.Title {
		case "status":
			values = append(values, task.Status.Status)
		case "name":
			values = append(values, task.Name)
		case "assignee":
			values = append(values, task.GetAssignees())
		case "list":
			values = append(values, task.List.String())
		case "tags":
			values = append(values, task.GetTags())
		case "folder":
			values = append(values, task.Folder.String())
		case "url":
			values = append(values, task.Url)
		case "space":
			values = append(values, task.Space.Id)
		case "id":
			values = append(values, task.Id)
		default:
			values = append(values, "XXX")
		}
	}

	return values
}

func (m Model) getTasksCmd(view string) tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.getTasks(view)
		if err != nil {
			return common.ErrMsg(err)
		}

		return TasksListReloadedMsg(tasks)
	}
}

func (m Model) getTasks(view string) ([]clickup.Task, error) {
	m.ctx.Logger.Infof("Getting tasks for view: %s", view)

	data, ok := m.ctx.Cache.Get("tasks", view)
	if ok {
		m.ctx.Logger.Infof("Tasks found in cache")
		var tasks []clickup.Task
		if err := m.ctx.Cache.ParseData(data, &tasks); err != nil {
			return nil, err
		}

		return tasks, nil
	}
	m.ctx.Logger.Info("Tasks not found in cache")

	m.ctx.Logger.Info("Fetching tasks from API")
	client := m.ctx.Clickup

	tasks, err := client.GetTasksFromView(view)
	if err != nil {
		return nil, err
	}
	m.ctx.Logger.Infof("Found %d tasks in view %s", len(tasks), view)

	m.ctx.Logger.Info("Caching tasks")
	m.ctx.Cache.Set("tasks", view, tasks)

	return tasks, nil
}
