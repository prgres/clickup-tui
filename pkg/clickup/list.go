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
	Id         string `json:"id"`
	Name       string `json:"name"`
	OrderIndex int    `json:"orderindex"`
	Content    string `json:"content"`
	Status     string `json:"status"`
	// Priority         string     `json:"priority"`
	Assignee         string     `json:"assignee"`
	TaskCount        int        `json:"task_count"`
	DueDate          string     `json:"due_date"`
	StartDate        string     `json:"start_date"`
	Folder           ListFolder `json:"folder"`
	Space            ListSpace  `json:"space"`
	Archived         bool       `json:"archived"`
	OverrideStatuses bool       `json:"override_statuses"`
	PermissionLevel  string     `json:"permission_level"`
}

type RequestGetLists struct {
	Lists []List `json:"lists"`
	Err   string `json:"err"`
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
			"Error occurs while getting lists from folders: %s. API response: %s",
			folderId, string(rawData))
	}

	return objmap.Lists, nil
}
