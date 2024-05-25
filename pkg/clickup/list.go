package clickup

import (
	"encoding/json"
)

type ListFolder struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Hidden bool   `json:"hidden"`
	Access bool   `json:"access"`
}

type ListSpace struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Access bool   `json:"access"`
}

type List struct {
	StartDate        string     `json:"start_date"`
	Name             string     `json:"name"`
	PermissionLevel  string     `json:"permission_level"`
	Content          string     `json:"content"`
	Status           string     `json:"status"`
	Assignee         string     `json:"assignee"`
	Id               string     `json:"id"`
	DueDate          string     `json:"due_date"`
	Folder           ListFolder `json:"folder"`
	Space            ListSpace  `json:"space"`
	TaskCount        int        `json:"task_count"`
	OrderIndex       int        `json:"orderindex"`
	Archived         bool       `json:"archived"`
	OverrideStatuses bool       `json:"override_statuses"`
}

type RequestGetLists struct {
	Lists []List `json:"lists"`
	Err   string `json:"err"`
}

func (r RequestGetLists) Error() string {
	return r.Err
}

func (c *Client) GetListsFromFolder(folderId string) ([]List, error) {
	return c.getLists("/folder/" + folderId + "/list")
}

func (c *Client) getLists(url string) ([]List, error) {
	var objmap RequestGetLists
	if err := c.get(url, &objmap); err != nil {
		return nil, err
	}
	return objmap.Lists, nil
}

type RequestGetList struct {
	List List   `json:"list"`
	Err  string `json:"err"`
}

func (r RequestGetList) Error() string {
	return r.Err
}

func (c *Client) GetList(listId string) (List, error) {
	rawData, err := c.requestGet("/list/" + listId)
	if err != nil {
		return List{}, err
	}

	var objmap List

	if err := json.Unmarshal(rawData, &objmap); err != nil {
		return List{}, err
	}

	return objmap, nil
}
