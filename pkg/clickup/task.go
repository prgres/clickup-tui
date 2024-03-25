package clickup

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Task struct {
	Startdate           interface{}   `json:"startdate"`
	Duedate             interface{}   `json:"duedate"`
	Priority            interface{}   `json:"priority"`
	Parent              interface{}   `json:"parent"`
	Timeestimate        interface{}   `json:"timeestimate"`
	Timespent           interface{}   `json:"timespent"`
	DateCreated         string        `json:"date_created"`
	Orderindex          string        `json:"orderindex"`
	Id                  string        `json:"id"`
	DateUpdated         string        `json:"date_updated"`
	DateClosed          string        `json:"date_closed"`
	DateDone            string        `json:"date_done"`
	Url                 string        `json:"url"`
	Space               TaskSpace     `json:"space"`
	MarkdownDescription string        `json:"markdown_description"`
	Description         string        `json:"description"`
	TextContent         string        `json:"text_content"`
	Name                string        `json:"name"`
	CustomId            string        `json:"custom_id"`
	Status              Status        `json:"status"`
	Creator             Creator       `json:"creator"`
	List                TaskList      `json:"list"`
	Folder              TaskFolder    `json:"folder"`
	Tags                []TaskTag     `json:"tags"`
	Checklists          []interface{} `json:"checklists"`
	Assignees           []Assignee    `json:"assignees"`
}

type TaskTag struct {
	Name    string `json:"name"`
	Tag_bg  string `json:"tag_bg"`
	Tag_fg  string `json:"tag_fg"`
	Created int64  `json:"created"`
}

func (t Task) GetTags() string {
	tags := strings.Builder{}
	for _, tag := range t.Tags {
		tags.WriteString(tag.Name)
	}

	return tags.String()
}

type TaskList struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Access bool   `json:"access"`
}

func (t TaskList) String() string {
	return fmt.Sprintf("%s (%s)", t.Name, t.Id)
}

type TaskFolder struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Access bool   `json:"access"`
	Hidden bool   `json:"hidden"`
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
	Username       string `json:"username"`
	Color          string `json:"color"`
	Initials       string `json:"initials"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
	Id             uint   `json:"id"`
}

type Status struct {
	Status     string `json:"status"`
	Color      string `json:"color"`
	Type       string `json:"type"`
	Orderindex int    `json:"orderindex"`
}

type Creator struct {
	Username       string `json:"username"`
	Color          string `json:"color"`
	ProfilePicture string `json:"profilePicture"`
	Id             int    `json:"id"`
}

type RequestGetTasks struct {
	Err      string `json:"err"`
	Tasks    []Task `json:"tasks"`
	LastPage bool   `json:"last_page"`
}

type RequestGetTask struct {
	Err  string `json:"err"`
	Task Task   `json:"task"`
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
			"error occurs while getting tasks from view: %s. API response: %s",
			viewId, string(rawData))
	}

	return objmap.Tasks, nil
}

func (c *Client) GetTasksFromList(listId string) ([]Task, error) {
	rawData, err := c.requestGet("/list/" + listId + "/task")
	if err != nil {
		return nil, err
	}
	var objmap RequestGetTasks

	if err := json.Unmarshal(rawData, &objmap); err != nil {
		return nil, err
	}

	if objmap.Err != "" {
		return nil, fmt.Errorf(
			"error occurs while getting tasks from list: %s. API response: %s",
			listId, string(rawData))
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
