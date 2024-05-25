package clickup

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
	Spaces []Space `json:"spaces"`
	Err    string  `json:"err"`
}

func (r RequestGetSpaces) Error() string {
	return r.Err
}

func (c *Client) GetSpacesFromTeam(teamId string) ([]Space, error) {
	return c.getSpaces("/team/" + teamId + "/space")
}

func (c *Client) getSpaces(url string) ([]Space, error) {
	var objmap RequestGetSpaces
	if err := c.get(url, &objmap); err != nil {
		return nil, err
	}
	return objmap.Spaces, nil
}
