package clickup

type Folder struct {
	Id               string        `json:"id"`
	Name             string        `json:"name"`
	TaskCount        string        `json:"task_count"`
	Space            Space         `json:"space"`
	Lists            []FolderSpace `json:"lists"`
	OrderIndex       int           `json:"orderindex"`
	OverrideStatuses bool          `json:"override_statuses"`
	Hidden           bool          `json:"hidden"`
}

type FolderSpace struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Access bool   `json:"access"`
}

type RequestGetFolders struct {
	Folders []Folder `json:"folders"`
	Err     string   `json:"err"`
}

func (r RequestGetFolders) Error() string {
	return r.Err
}

func (c *Client) GetFolders(spaceId string) ([]Folder, error) {
	return c.getFolders("/space/" + spaceId + "/folder")
}

func (c *Client) getFolders(url string) ([]Folder, error) {
	var objmap RequestGetFolders
	if err := c.get(url, &objmap); err != nil {
		return nil, err
	}

	return objmap.Folders, nil
}
