package taskstable

import (
	"github.com/mattn/go-runewidth"
	"github.com/prgrs/clickup/pkg/clickup"
)

func taskListToRows(tasks []clickup.Task, columns []string) [][]string {
	rows := make([][]string, len(tasks))
	for i, task := range tasks {
		rows[i] = taskToRow(task, columns)
	}
	return rows
}

func (m Model) getSelectedViewTaskIdByIndex(index int) string {
	return m.getSelectedViewTasks()[index-1].Id
}

func (m Model) getSelectedViewTasks() []clickup.Task {
	m.log.Infof("getSelectedViewTasks: %v", m.SelectedTab.Id)
	return m.tasks[m.SelectedTab.Id]
}

func taskToRow(task clickup.Task, columns []string) []string {
	values := []string{}
	for _, column := range columns {
		switch column {
		case "status":
			values = append(values, task.Status.Status)
		case "name":
			n := runewidth.Wrap(task.Name, 30)
			values = append(values, n)
			// values = append(values, task.Name)
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
