package clickup

import "encoding/json"

type Workspace = Team

type Team struct {
	Id      string        `json:"id"`
	Name    string        `json:"name"`
	Color   string        `json:"color"`
	Avatar  string        `json:"avatar"`
	Members []interface{} `json:"members"`
}

type RequestGetTeams struct {
	Teams []Team `json:"teams"`
}

func (c *Client) GetTeams() ([]Team, error) {
	rawData, err := c.requestGet("/team")
	if err != nil {
		return nil, err
	}

	var objmap RequestGetTeams
	if err := json.Unmarshal(rawData, &objmap); err != nil {
		return nil, err
	}

	return objmap.Teams, nil
}
