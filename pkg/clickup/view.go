package clickup

import (
	"encoding/json"
	"fmt"
)

type View struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Type        string          `json:"type"`
	Parent      ViewParent      `json:"parent"`
	Grouping    ViewGrouping    `json:"grouping"`
	Divide      ViewDivide      `json:"divide"`
	Sorting     ViewSorting     `json:"sorting"`
	Filter      ViewFilter      `json:"filters"`
	Columns     ViewColumns     `json:"columns"`
	TeamSidebar ViewTeamSidebar `json:"team_sidebar"`
	Settings    ViewSettings    `json:"settings"`
}

type ViewParent struct {
	Id   string `json:"id"`
	Type int    `json:"type"`
}

type ViewGrouping struct {
	Field     string   `json:"field"`
	Dir       int      `json:"dir"`
	Collapsed []string `json:"collapsed"`
	Ignore    bool     `json:"ignore"`
}

type ViewDivide struct {
	Field     string   `json:"field"`
	Dir       int      `json:"dir"`
	Collapsed []string `json:"collapsed"`
}

type ViewSorting struct {
	Fields []interface{} `json:"fields"`
}

type ViewFilter struct {
	Op         string        `json:"op"`
	Fields     []interface{} `json:"fields"`
	Search     string        `json:"search"`
	ShowClosed bool          `json:"show_closed"`
}

type ViewColumns struct {
	Fields []interface{} `json:"fields"`
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
	return []View{r.List, r.Board, r.Box, r.Calendar}
}

func (c *Client) GetViewsFromSpace(spaceId string) ([]View, error) {
	rawData, err := c.requestGet("/space/" + spaceId + "/view")
	if err != nil {
		return nil, err
	}

	var objmap RequestGetViews
	if err := json.Unmarshal(rawData, &objmap); err != nil {
		return nil, err
	}
	if objmap.Err != "" {
		return nil, fmt.Errorf(objmap.Err)
	}
	return append(objmap.Views, objmap.RequiredViews.GetViews()...), nil
}
