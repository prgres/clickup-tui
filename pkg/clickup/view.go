package clickup

type ViewType string

const (
	ViewTypeCalendar     ViewType = "calendar"
	ViewTypeGantt        ViewType = "gantt"
	ViewTypeTable        ViewType = "table"
	ViewTypeList         ViewType = "list"
	ViewTypeDoc          ViewType = "doc"
	ViewTypeBoard        ViewType = "board"
	ViewTypeTimeline     ViewType = "Timeline"
	ViewTypeWorkload     ViewType = "workload"
	ViewTypeActivity     ViewType = "activity"
	ViewTypeMap          ViewType = "map"
	ViewTypeConversation ViewType = "conversation"
)

type View struct {
	Name          string          `json:"name"`
	Type          ViewType        `json:"type"`
	DateProtected string          `json:"date_protected"`
	Id            string          `json:"id"`
	ProtectedBy   ViewProtectedBy `json:"protected_by"`
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

type ViewProtectedBy struct {
	Id             int    `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Color          string `json:"color"`
	Initials       string `json:"initials"`
	ProfilePicture string `json:"profilePicture"`
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
	RequiredViews RequiredViews `json:"required_views"`
	Err           string        `json:"err"`
}

func (r RequestGetViews) Error() string {
	return r.Err
}

type RequiredViews struct {
	List     View `json:"list"`
	Board    View `json:"board"`
	Box      View `json:"box"`
	Calendar View `json:"calendar"`
}

func (r RequiredViews) GetNonEmptyViews() []View {
	views := []View{}
	if r.List.Id != "" {
		views = append(views, r.List)
	}

	if r.Board.Id != "" {
		views = append(views, r.Board)
	}

	if r.Box.Id != "" {
		views = append(views, r.Box)
	}

	if r.Calendar.Id != "" {
		views = append(views, r.Calendar)
	}

	return views // []View{r.List, r.Board, r.Box, r.Calendar}
}

func (c *Client) GetViewsFromWorkspace(workspaceId string) ([]View, error) {
	return c.getViews("/team/" + workspaceId + "/view")
}

func (c *Client) GetViewsFromSpace(spaceId string) ([]View, error) {
	return c.getViews("/space/" + spaceId + "/view")
}

func (c *Client) GetViewsFromFolder(folderId string) ([]View, error) {
	return c.getViews("/folder/" + folderId + "/view")
}

func (c *Client) GetViewsFromList(listId string) ([]View, error) {
	return c.getViews("/list/" + listId + "/view")
}

func (c *Client) getViews(url string) ([]View, error) {
	var objmap RequestGetViews
	if err := c.get(url, &objmap); err != nil {
		return nil, err
	}

	return append(objmap.Views,
		objmap.RequiredViews.GetNonEmptyViews()...,
	), nil
}
