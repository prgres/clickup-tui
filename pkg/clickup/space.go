package clickup

import (
	"encoding/json"
	"fmt"
)

type Space struct {
	Id                string        `json:"id"`
	Name              string        `json:"name"`
	Statuses          []SpaceStatus `json:"statuses"`
	Features          []interface{} `json:"-"`
	Private           bool          `json:"private"`
	MultipleAssignees bool          `json:"multiple_assignees"`
}

type SpaceStatus struct {
	Status     string
	Type       string
	Color      string
	OrderIndex int
}

type RequestGetSpaces struct {
	Err    string  `json:"err"`
	Spaces []Space `json:"spaces"`
}

func (c *Client) GetSpaces(teamId string) ([]Space, error) {
	rawData, err := c.requestGet("/team/" + teamId + "/space")
	if err != nil {
		return nil, err
	}

	var objmap RequestGetSpaces
	if err := json.Unmarshal(rawData, &objmap); err != nil {
		return nil, err
	}

	if objmap.Err != "" {
		return nil, fmt.Errorf(objmap.Err)
	}

	return objmap.Spaces, nil
}
