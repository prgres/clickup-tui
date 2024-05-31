package clickup

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Task struct {
	Startdate           interface{}   `json:"start_date"`
	Duedate             interface{}   `json:"due_date"`
	Priority            string        `json:"priority"`
	Parent              interface{}   `json:"parent"`
	Timeestimate        interface{}   `json:"time_estimate"`
	Timespent           interface{}   `json:"time_spent"`
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
	Points              int           `json:"points"`
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
	Tasks    []Task `json:"tasks"`
	LastPage bool   `json:"last_page"`
	Err      string `json:"err"`
}

func (r RequestGetTasks) Error() string {
	return r.Err
}

type RequestGetTask struct {
	Task Task   `json:"task"`
	Err  string `json:"err"`
}

type Assignees struct {
	Add []int `json:"add,omitempty"`
	Rem []int `json:"rem,omitempty"`
}

type Watchers struct {
	Add []int `json:"add,omitempty"`
	Rem []int `json:"rem,omitempty"`
}

type GroupAssignees struct {
	Add []int `json:"add,omitempty"`
	Rem []int `json:"rem,omitempty"`
}

type RequestPutTask struct {
	CustomItemId   int            `json:"custom_item_id,omitempty"`
	Id             string         `json:"id"`
	Name           string         `json:"name,omitempty"`
	Description    string         `json:"description,omitempty"`
	Status         string         `json:"status"`
	Priority       int32          `json:"priority,omitempty"`
	DueDate        int64          `json:"due_date,omitempty"`
	DueDateTime    bool           `json:"due_date_time,omitempty"`
	Parent         string         `json:"parent,omitempty"`
	TimeEstimate   int32          `json:"time_estimate,omitempty"`
	StartDate      int64          `json:"start_date,omitempty"`
	StartDateTime  bool           `json:"start_date_time,omitempty"`
	Points         int            `json:"points,omitempty"`
	Assignees      Assignees      `json:"assignees,omitempty"`
	GroupAssignees GroupAssignees `json:"group_assignees,omitempty"`
	Watchers       Watchers       `json:"watchers,omitempty"`
	Archived       bool           `json:"archived,omitempty"`
}

func (r RequestGetTask) Error() string {
	return r.Err
}

func (c *Client) GetTasksFromView(viewId string) ([]Task, error) {
	return c.getTasks("/view/" + viewId + "/task")
}

func (c *Client) GetTasksFromList(listId string) ([]Task, error) {
	return c.getTasks("/list/" + listId + "/task")
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

func (c *Client) getTasks(url string) ([]Task, error) {
	var objmap RequestGetTasks
	if err := c.get(url, &objmap); err != nil {
		return nil, err
	}
	return objmap.Tasks, nil
}

func (c *Client) UpdateTask(r RequestPutTask) (Task, error) {
	var objmap Task

	if err := c.update("/task/"+r.Id, r, &objmap); err != nil {
		return objmap, err
	}

	return objmap, nil
}
