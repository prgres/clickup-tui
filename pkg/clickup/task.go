package clickup

import "encoding/json"

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
	Assignees    []string      `json:"assignees"`
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
	return objmap.Tasks, nil
}
