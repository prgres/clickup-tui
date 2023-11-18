package tasktable

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/prgrs/clickup/pkg/clickup"
)

func taskListToRows(tasks []clickup.Task, columns []table.Column) []table.Row {
	rows := make([]table.Row, len(tasks))
	for i, task := range tasks {
		rows[i] = taskToRow(task, columns)
	}
	return rows
}

func (m Model) getSelectedViewTaskIdByIndex(index int) string {
	return m.getSelectedViewTasks()[index].Id
}

func (m Model) getSelectedViewTasks() []clickup.Task {
	return m.tasks[m.SelectedTab.Id]
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
