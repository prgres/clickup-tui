package tabletasks

import (
	"github.com/evertras/bubble-table/table"
	"github.com/prgrs/clickup/pkg/clickup"
)

func taskListToRows(tasks []clickup.Task, columns []string) []table.Row {
	rows := make([]table.Row, len(tasks))
	for i, task := range tasks {
		rows[i] = taskToRow(task, columns)
	}
	return rows
}

func taskToRow(task clickup.Task, columns []string) table.Row {
	values := map[string]interface{}{}
	for _, column := range columns {
		switch column {
		case "status":
			values[column] = task.Status.Status
		case "name":
			values[column] = task.Name
		case "url":
			values[column] = task.Url
			// After migration from charm to evertras/bubble-table I temporary removed all columns
			// except "status" and "name" since they are not supported yet. See autoColumns feature
			// case "assignee":
			// 	values = append(values, task.GetAssignees())
			// case "list":
			// 	values = append(values, task.List.String())
			// case "tags":
			// 	values = append(values, task.GetTags())
			// case "folder":
			// 	values = append(values, task.Folder.String())
			// case "space":
			// 	values = append(values, task.Space.Id)
			// case "id":
			// 	values = append(values, task.Id)
		}
	}

	return table.NewRow(table.RowData(values))
}

func rowToTask(row table.Row, columns []string) clickup.Task {
	var task clickup.Task
	data := row.Data

	for _, column := range columns {
		switch column {
		case "status":
			task.Status.Status = data[column].(string)
		case "name":
			task.Name = data[column].(string)
		case "url":
			task.Url = data[column].(string)
			// After migration from charm to evertras/bubble-table I temporary removed all columns
			// except "status" and "name" since they are not supported yet. See autoColumns feature
			// case "assignee":
			// 	values = append(values, task.GetAssignees())
			// case "list":
			// 	values = append(values, task.List.String())
			// case "tags":
			// 	values = append(values, task.GetTags())
			// case "folder":
			// 	values = append(values, task.Folder.String())
			// case "url":
			// 	values = append(values, task.Url)
			// case "space":
			// 	values = append(values, task.Space.Id)
			// case "id":
			// 	values = append(values, task.Id)
		}
	}

	return task
}
