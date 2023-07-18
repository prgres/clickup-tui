package clickup

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Task struct {
	Id           string        `json:"id"`
	Name         string        `json:"name"`
	Status       Status        `json:"status"`
	Orderindex   string        `json:"orderindex"`
	DateCreated  string        `json:"date_created"`
	DateUpdated  string        `json:"date_updated"`
	DateClosed   interface{}   `json:"date_closed"`
	DateDone     interface{}   `json:"date_done"`
	Creator      Creator       `json:"creator"`
	Assignees    []Assignee    `json:"assignees"`
	Checklists   []interface{} `json:"checklists"`
	Tags         []interface{} `json:"tags"`
	Parent       interface{}   `json:"parent"`
	Priority     interface{}   `json:"priority"`
	Duedate      interface{}   `json:"duedate"`
	Startdate    interface{}   `json:"startdate"`
	Timeestimate interface{}   `json:"timeestimate"`
	Timespent    interface{}   `json:"timespent"`
	List         interface{}   `json:"list"`
	Folder       interface{}   `json:"folder"`
	Space        interface{}   `json:"space"`
	Url          string        `json:"url"`
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
	Id             uint   `json:"id"`
	Username       string `json:"username"`
	Color          string `json:"color"`
	ProfilePicture string `json:"profilePicture"`
}

type RequestGetTasks struct {
	Tasks    []Task `json:"tasks"`
	LastPage bool   `json:"last_page"`
	Err      string `json:"err"`
}

func (c *Client) GetTasksFromView(viewId string) ([]Task, error) {
	rawData, err := c.requestGet("/view/" + viewId + "/task")
	if err != nil {
		return nil, err
	}
	var objmap RequestGetTasks
	// c.logger.Info("rawData ", string(rawData))
	if err := json.Unmarshal(rawData, &objmap); err != nil {
		return nil, err
	}

	if objmap.Err != "" {
		return nil, fmt.Errorf(objmap.Err)
	}
	// c.logger.Info("tasks len ", len(objmap.Tasks))
	return objmap.Tasks, nil
}
