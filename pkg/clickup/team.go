package clickup

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
	Err   string `json:"err"`
}

func (r RequestGetTeams) Error() string {
	return r.Err
}

func (c *Client) GetTeams() ([]Team, error) {
	return c.getTeams("/team")
}

func (c *Client) getTeams(url string) ([]Team, error) {
	var objmap RequestGetTeams
	if err := c.get(url, &objmap); err != nil {
		return nil, err
	}
	return objmap.Teams, nil
}
