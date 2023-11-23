package clickup

import (
	"encoding/json"
	"fmt"
)

type Space struct {
	Id                string
	Name              string        `json:"name"`
	Private           bool          `json:"private"`
	Statuses          []SpaceStatus `json:"statuses"`
	MultipleAssignees bool          `json:"multiple_assignees"`
	Features          []interface{} `json:"-"` //	`json:"features"`
}

type SpaceStatus struct {
	Status     string
	Type       string
	OrderIndex int
	Color      string
}

type RequestGetSpaces struct {
	Spaces []Space `json:"spaces"`
	Err    string  `json:"err"`
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
