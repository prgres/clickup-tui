package clickup

import "encoding/json"

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
}

func (c *Client) GetFolders(spaceId string) ([]Folder, error) {
	rawData, err := c.requestGet("/space/" + spaceId + "/folder")
	if err != nil {
		return nil, err
	}

	var objmap RequestGetFolders
	if err := json.Unmarshal(rawData, &objmap); err != nil {
		return nil, err
	}

	return objmap.Folders, nil
}
