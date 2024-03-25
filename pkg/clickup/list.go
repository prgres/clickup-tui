package clickup

import (
	"encoding/json"
	"fmt"
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
	Err   string `json:"err"`
	Lists []List `json:"lists"`
}

func (c *Client) GetListsFromFolder(folderId string) ([]List, error) {
	rawData, err := c.requestGet("/folder/" + folderId + "/list")
	if err != nil {
		return nil, err
	}
	var objmap RequestGetLists

	if err := json.Unmarshal(rawData, &objmap); err != nil {
		return nil, err
	}

	if objmap.Err != "" {
		return nil, fmt.Errorf(
			"error occurs while getting lists from folders: %s. API response: %s",
			folderId, string(rawData))
	}

	return objmap.Lists, nil
}
