package clickup

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Task struct {
	Id                  string        `json:"id"`
	CustomId            string        `json:"custom_id"`
	Name                string        `json:"name"`
	TextContent         string        `json:"text_content"`
	Description         string        `json:"description"`
	MarkdownDescription string        `json:"markdown_description"`
	Status              Status        `json:"status"`
	Orderindex          string        `json:"orderindex"`
	DateCreated         string        `json:"date_created"`
	DateUpdated         string        `json:"date_updated"`
	DateClosed          string        `json:"date_closed"`
	DateDone            string        `json:"date_done"`
	Creator             Creator       `json:"creator"`
	Assignees           []Assignee    `json:"assignees"`
	Checklists          []interface{} `json:"checklists"`
	Tags                []TaskTag     `json:"tags"`
	Parent              interface{}   `json:"parent"`
	Priority            interface{}   `json:"priority"`
	Duedate             interface{}   `json:"duedate"`
	Startdate           interface{}   `json:"startdate"`
	Timeestimate        interface{}   `json:"timeestimate"`
	Timespent           interface{}   `json:"timespent"`
	List                TaskList      `json:"list"`
	Folder              TaskFolder    `json:"folder"`
	Space               TaskSpace     `json:"space"`
	Url                 string        `json:"url"`
}

type TaskTag struct {
	Created int64  `json:"created"`
	Name    string `json:"name"`
	Tag_bg  string `json:"tag_bg"`
	Tag_fg  string `json:"tag_fg"`
}

func (t Task) GetTags() string {
	tags := strings.Builder{}
	for _, tag := range t.Tags {
		tags.WriteString(tag.Name)
	}

	return tags.String()
}

type TaskList struct {
	Access bool   `json:"access"`
	Id     string `json:"id"`
	Name   string `json:"name"`
}

func (t TaskList) String() string {
	return fmt.Sprintf("%s (%s)", t.Name, t.Id)
}

type TaskFolder struct {
	Access bool   `json:"access"`
	Hidden bool   `json:"hidden"`
	Id     string `json:"id"`
	Name   string `json:"name"`
}

func (t TaskFolder) String() string {
	return fmt.Sprintf("%s (%s)", t.Name, t.Id)
}

type TaskSpace struct {
	Id string `json:"id"`
}

func (t Task) GetAssignees() string {
	assignees := []string{}
	for _, assignee := range t.Assignees {
		assignees = append(assignees, assignee.Username)
	}
	return strings.Join(assignees, ",")
}

type Assignee struct {
	Id             uint   `json:"id"`
	Username       string `json:"username"`
	Color          string `json:"color"`
	Initials       string `json:"initials"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
}

type Status struct {
	Status     string `json:"status"`
	Color      string `json:"color"`
	Orderindex int    `json:"orderindex"`
	Type       string `json:"type"`
}

type Creator struct {
	Id             int    `json:"id"`
	Username       string `json:"username"`
	Color          string `json:"color"`
	ProfilePicture string `json:"profilePicture"`
}

type RequestGetTasks struct {
	Tasks    []Task `json:"tasks"`
	LastPage bool   `json:"last_page"`
	Err      string `json:"err"`
}

type RequestGetTask struct {
	Task Task   `json:"task"`
	Err  string `json:"err"`
}

func (c *Client) GetTasksFromView(viewId string) ([]Task, error) {
	rawData, err := c.requestGet("/view/" + viewId + "/task")
	if err != nil {
		return nil, err
	}
	var objmap RequestGetTasks

	if err := json.Unmarshal(rawData, &objmap); err != nil {
		return nil, err
	}

	if objmap.Err != "" {
		return nil, fmt.Errorf(
			"Error occurs while getting tasks from view: %s. API response: %s",
			viewId, string(rawData))
	}

	return objmap.Tasks, nil
}

func (c *Client) GetTask(taskId string) (Task, error) {
	rawData, err := c.requestGet("/task/"+taskId, "include_markdown_description", "true")
	if err != nil {
		return Task{}, err
	}
	var objmap Task

	if err := json.Unmarshal(rawData, &objmap); err != nil {
		return Task{}, err
	}

	return objmap, nil
}
