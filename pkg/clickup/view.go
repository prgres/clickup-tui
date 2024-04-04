package clickup

import (
	"encoding/json"
	"fmt"
)

type View struct {
	Name          string          `json:"name"`
	Type          string          `json:"type"`
	DateProtected string          `json:"date_protected"`
	Id            string          `json:"id"`
	ProtectedBy   string          `json:"protected_by"`
	ProtectedNote string          `json:"protected_note"`
	Visibility    string          `json:"visibility"`
	DateCreated   string          `json:"date_created"`
	Parent        ViewParent      `json:"parent"`
	Sorting       ViewSorting     `json:"sorting"`
	Columns       ViewColumns     `json:"columns"`
	Filter        ViewFilter      `json:"filters"`
	Divide        ViewDivide      `json:"divide"`
	TeamSidebar   ViewTeamSidebar `json:"team_sidebar"`
	Grouping      ViewGrouping    `json:"grouping"`
	Settings      ViewSettings    `json:"settings"`
	Creator       int             `json:"creator"`
	OrderIndex    int             `json:"order_index"`
	Protected     bool            `json:"protected"`
}

type ViewParent struct {
	Id   string `json:"id"`
	Type int    `json:"type"`
}

type ViewGrouping struct {
	Field     string   `json:"field"`
	Collapsed []string `json:"collapsed"`
	Dir       int      `json:"dir"`
	Ignore    bool     `json:"ignore"`
}

type ViewDivide struct {
	Field     string   `json:"field"`
	Collapsed []string `json:"collapsed"`
	Dir       int      `json:"dir"`
}

type ViewSorting struct {
	Fields []interface{} `json:"fields"`
}

type ViewFilter struct {
	Op         string        `json:"op"`
	Search     string        `json:"search"`
	Fields     []interface{} `json:"fields"`
	ShowClosed bool          `json:"show_closed"`
}

type ViewColumns struct {
	Fields []ColumnField `json:"fields"`
}

func (v ViewColumns) GetColumnsFields() []string {
	fields := []string{}
	for _, column := range v.Fields {
		fields = append(fields, column.Field)
	}

	return fields
}

type ColumnField struct {
	Field   string `json:"field"`
	Name    string `json:"name"`
	Display string `json:"display"`
	Idx     int    `json:"idx"`
	Width   int    `json:"width"`
	Hidden  bool   `json:"hidden"`
}

type ViewTeamSidebar struct {
	Assignees        []string `json:"assignees"`
	AssignedComments bool     `json:"assigned_comments"`
	UnassignedTasks  bool     `json:"unassigned_tasks"`
}

type ViewSettings struct {
	ShowTaskLocations      bool `json:"show_task_locations"`
	ShowSubtasks           int  `json:"show_subtasks"`
	ShowSubtaskParentNames bool `json:"show_subtask_parent_names"`
	ShowCompletedSubtasks  bool `json:"show_completed_subtasks"`
	ShowAssignees          bool `json:"show_assignees"`
	ShowImages             bool `json:"show_images"`
	CollapseEmptyColumns   bool `json:"collapse_empty_columns"`
	MeComments             bool `json:"me_comments"`
	MeSubtasks             bool `json:"me_subtasks"`
	MeChecklists           bool `json:"me_checklists"`
}

type RequestGetViews struct {
	Views         []View        `json:"views"`
	Err           string        `json:"err"`
	RequiredViews RequiredViews `json:"required_views"`
}

type RequiredViews struct {
	List     View `json:"list"`
	Board    View `json:"board"`
	Box      View `json:"box"`
	Calendar View `json:"calendar"`
}

func (r RequiredViews) GetViews() []View {
	views := []View{}
	if r.List.Id != "" {
		views = append(views, r.List)
	}

	// if r.Board.Id != "" {
	// 	views = append(views, r.Board)
	// }

	// if r.Box.Id != "" {
	// 	views = append(views, r.Box)
	// }

	// if r.Calendar.Id != "" {
	// 	views = append(views, r.Calendar)
	// }

	return views // []View{r.List, r.Board, r.Box, r.Calendar}
}

func filterListViews(views []View) []View {
	filteredViews := []View{}
	for _, view := range views {
		if view.Type == "list" {
			filteredViews = append(filteredViews, view)
		}
	}
	return filteredViews
}

func (c *Client) GetViewsFromWorkspace(workspaceId string) ([]View, error) {
	errMsg := "Error occurs while getting views from workspace: %s. Error: %s. Raw data: %s"

	errApiMsg := errMsg + " API response: %s"

	rawData, err := c.requestGet("/team/" + workspaceId + "/view")
	if err != nil {
		return nil, fmt.Errorf(errMsg, workspaceId, err, "none")
	}

	var objmap RequestGetViews
	if err := json.Unmarshal(rawData, &objmap); err != nil {
		return nil, fmt.Errorf(
			errApiMsg, workspaceId, err, string(rawData))
	}

	if objmap.Err != "" {
		return nil, fmt.Errorf(
			errApiMsg, workspaceId, objmap.Err, string(rawData))
	}

	allViews := append(objmap.Views, objmap.RequiredViews.GetViews()...)
	for _, v := range allViews {
		if v.Id == "" || v.Name == "" {
			return nil, fmt.Errorf(
				"View id or name is empty, API response: %s", string(rawData))
		}
	}

	if len(allViews) == 0 {
		c.logger.Error("No views found in workspace",
			"workspace", workspaceId)
		return []View{}, nil
	}

	return allViews, nil
}

func (c *Client) GetViewsFromSpace(spaceId string) ([]View, error) {
	errMsg := "Error occurs while getting views from space: %s. Error: %s. Raw data: %s"

	errApiMsg := errMsg + " API response: %s"

	rawData, err := c.requestGet("/space/" + spaceId + "/view")
	if err != nil {
		return nil, fmt.Errorf(errMsg, spaceId, err, "none")
	}

	var objmap RequestGetViews
	if err := json.Unmarshal(rawData, &objmap); err != nil {
		return nil, fmt.Errorf(
			errApiMsg, spaceId, err, string(rawData))
	}

	if objmap.Err != "" {
		return nil, fmt.Errorf(
			errMsg, spaceId, "API response contains error.", string(rawData))
	}

	allViews := append(objmap.Views, objmap.RequiredViews.GetViews()...)
	// TODO: find out why it was needed
	for _, v := range allViews {
		if v.Id == "" || v.Name == "" {
			return nil, fmt.Errorf(
				"View id or name is empty, API response: %s", string(rawData))
		}
	}

	if len(allViews) == 0 {
		c.logger.Error("No views found in space",
			"space", spaceId)
		return []View{}, nil
	}

	// filteredViews := filterListViews(allViews)

	return objmap.RequiredViews.GetViews(), nil
	// return append(filteredViews, objmap.RequiredViews.GetViews()...), nil
}

func (c *Client) GetViewsFromFolder(folderId string) ([]View, error) {
	errMsg := "Error occurs while getting views from folder: %s. Error: %s"
	errApiMsg := errMsg + " API response: %s"

	rawData, err := c.requestGet("/folder/" + folderId + "/view")
	if err != nil {
		return nil, fmt.Errorf(errMsg,
			folderId, err)
	}

	var objmap RequestGetViews
	if err := json.Unmarshal(rawData, &objmap); err != nil {
		return nil, fmt.Errorf(errApiMsg,
			folderId, err, string(rawData))
	}

	if objmap.Err != "" {
		return nil, fmt.Errorf(errApiMsg,
			folderId, objmap.Err, string(rawData))
	}

	allViews := append(objmap.Views, objmap.RequiredViews.GetViews()...)
	for _, v := range allViews {
		if v.Id == "" || v.Name == "" {
			return nil, fmt.Errorf(errApiMsg,
				folderId, "View id or name is empty", string(rawData))
		}
	}

	if len(allViews) == 0 {
		return []View{}, fmt.Errorf(errMsg,
			folderId, "No views found in folder")
	}

	filteredViews := filterListViews(allViews)

	return append(filteredViews, objmap.RequiredViews.GetViews()...), nil
}

func (c *Client) GetViewsFromList(listId string) ([]View, error) {
	errMsg := "Error occurs while getting views from list: %s. Error: %s"
	errApiMsg := errMsg + " API response: %s"

	rawData, err := c.requestGet("/list/" + listId + "/view")
	if err != nil {
		return nil, fmt.Errorf(errMsg,
			listId, err)
	}

	var objmap RequestGetViews
	if err := json.Unmarshal(rawData, &objmap); err != nil {
		return nil, fmt.Errorf(errApiMsg,
			listId, err, string(rawData))
	}

	if objmap.Err != "" {
		return nil, fmt.Errorf(errApiMsg,
			listId, objmap.Err, string(rawData))
	}

	allViews := append(objmap.Views, objmap.RequiredViews.GetViews()...)
	for _, v := range allViews {
		if v.Id == "" || v.Name == "" {
			return nil, fmt.Errorf(errApiMsg,
				listId, "View id or name is empty", string(rawData))
		}
	}

	if len(allViews) == 0 {
		return []View{}, nil
	}

	// filteredViews := filterListViews(allViews)

	return allViews, nil
	// return append(filteredViews, objmap.RequiredViews.GetViews()...), nil //TODO: find out why it was needed
}
