package clickup

import "encoding/json"

type Team struct {
	Id      string
	Name    string
	Color   string
	Avatar  string
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
