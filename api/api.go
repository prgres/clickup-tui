package api

import (
	"encoding/json"

	"github.com/prgrs/clickup/pkg/cache"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/pkg/logger1"
)

type Api struct {
	Clickup *clickup.Client
	Cache   *cache.Cache
	logger  logger1.Logger
}

func NewApi(clickup *clickup.Client, logger logger1.Logger, cache *cache.Cache) Api {
	return Api{
		Clickup: clickup,
		logger:  logger,
		Cache:   cache,
	}
}

func (m *Api) GetSpaces(team string) ([]clickup.Space, error) {
	m.logger.Infof("Getting spaces for team: %s", team)
	client := m.Clickup

	data, ok := m.Cache.Get("spaces", team)
	if ok {
		m.logger.Infof("Spaces found in cache")
		var spaces []clickup.Space
		if err := m.Cache.ParseData(data, &spaces); err != nil {
			return nil, err
		}

		return spaces, nil
	}
	m.logger.Infof("Spaces not found in cache")

	m.logger.Infof("Fetching spaces from API")
	spaces, err := client.GetSpaces(team)
	if err != nil {
		return nil, err
	}
	m.logger.Infof("Found %d spaces for team: %s", len(spaces), team)

	m.logger.Infof("Caching spaces")
	m.Cache.Set("spaces", "spaces", spaces)

	return spaces, nil
}

func (m *Api) GetWorkspaces() ([]clickup.Workspace, error) {
	teams, err := m.GetTeams()
	if err != nil {
		return nil, err
	}
	return teams, nil
}

func (m *Api) GetTeams() ([]clickup.Team, error) {
	m.logger.Info("Getting Authorized Teams (Workspaces)")
	client := m.Clickup

	data, ok := m.Cache.Get("teams", "teams")
	if ok {
		m.logger.Infof("Teams found in cache")
		var teams []clickup.Team
		if err := m.Cache.ParseData(data, &teams); err != nil {
			return nil, err
		}

		return teams, nil
	}
	m.logger.Infof("Teams not found in cache")

	m.logger.Infof("Fetching teams from API")
	teams, err := client.GetTeams()
	if err != nil {
		return nil, err
	}
	m.logger.Infof("Found %d teams ", len(teams))

	m.logger.Infof("Caching teams")
	m.Cache.Set("teams", "teams", teams)

	return teams, nil
}

func (m *Api) GetFolders(space string) ([]clickup.Folder, error) {
	m.logger.Infof("Getting folders for space: %s", space)
	client := m.Clickup

	data, ok := m.Cache.Get("folders", space)
	if ok {
		m.logger.Infof("Folders found in cache")
		var folders []clickup.Folder
		if err := m.Cache.ParseData(data, &folders); err != nil {
			return nil, err
		}

		return folders, nil
	}
	m.logger.Infof("Folders not found in cache")

	m.logger.Infof("Fetching folders from API")
	folders, err := client.GetFolders(space)
	if err != nil {
		return nil, err
	}
	m.logger.Infof("Found %d folders for space: %s", len(folders), space)

	m.logger.Infof("Caching folders")
	m.Cache.Set("folders", space, folders)

	return folders, nil
}

func (m *Api) GetLists(folderId string) ([]clickup.List, error) {
	m.logger.Infof("Getting lists for folder: %s", folderId)
	client := m.Clickup

	data, ok := m.Cache.Get("lists", folderId)
	if ok {
		m.logger.Infof("Lists found in cache")
		var lists []clickup.List
		if err := m.Cache.ParseData(data, &lists); err != nil {
			return nil, err
		}

		return lists, nil
	}
	m.logger.Infof("Lists not found in cache")

	m.logger.Infof("Fetching lists from API")
	lists, err := client.GetListsFromFolder(folderId)
	if err != nil {
		return nil, err
	}
	m.logger.Infof("Found %d lists for folder: %s", len(lists), folderId)

	for _, f := range lists {
		j, _ := json.MarshalIndent(f, "", "  ")
		m.logger.Infof("%s", j)
	}

	m.logger.Infof("Caching lists")
	m.Cache.Set("lists", folderId, lists)

	return lists, nil
}

func (m *Api) GetTask(id string) (clickup.Task, error) {
	m.logger.Infof("Getting task: %s", id)

	data, ok := m.Cache.Get("task", id)
	if ok {
		m.logger.Infof("Task found in cache")
		var task clickup.Task
		if err := m.Cache.ParseData(data, &task); err != nil {
			return clickup.Task{}, err
		}

		return task, nil
	}
	m.logger.Info("Task not found in cache")

	m.logger.Info("Fetching task from API")
	client := m.Clickup

	task, err := client.GetTask(id)
	if err != nil {
		return clickup.Task{}, err
	}
	m.logger.Infof("Found tasks %s", id)

	m.logger.Info("Caching tasks")
	m.Cache.Set("task", id, task)

	return task, nil
}

func (m *Api) GetTasksFromList(list string) ([]clickup.Task, error) {
	m.logger.Infof("Getting tasks for list: %s", list)

	data, ok := m.Cache.Get("tasks", list)
	if ok {
		m.logger.Infof("Tasks found in cache")
		var tasks []clickup.Task
		if err := m.Cache.ParseData(data, &tasks); err != nil {
			return nil, err
		}

		return tasks, nil
	}
	m.logger.Info("Tasks not found in cache")

	m.logger.Info("Fetching tasks from API")
	client := m.Clickup

	tasks, err := client.GetTasksFromList(list)
	if err != nil {
		return nil, err
	}
	m.logger.Infof("Found %d tasks in list %s", len(tasks), list)

	m.logger.Info("Caching tasks")
	m.Cache.Set("tasks", list, tasks)

	return tasks, nil
}
func (m *Api) GetTasksFromView(view string) ([]clickup.Task, error) {
	m.logger.Infof("Getting tasks for view: %s", view)

	data, ok := m.Cache.Get("tasks", view)
	if ok {
		m.logger.Infof("Tasks found in cache")
		var tasks []clickup.Task
		if err := m.Cache.ParseData(data, &tasks); err != nil {
			return nil, err
		}

		return tasks, nil
	}
	m.logger.Info("Tasks not found in cache")

	m.logger.Info("Fetching tasks from API")
	client := m.Clickup

	tasks, err := client.GetTasksFromView(view)
	if err != nil {
		return nil, err
	}
	m.logger.Infof("Found %d tasks in view %s", len(tasks), view)

	m.logger.Info("Caching tasks")
	m.Cache.Set("tasks", view, tasks)

	return tasks, nil
}

func (m *Api) GetViewsFromFolder(folder string) ([]clickup.View, error) {
	m.logger.Infof("Getting views for folder: %s", folder)

	data, ok := m.Cache.Get("views", folder)
	if ok {
		m.logger.Infof("Views found in cache")
		var views []clickup.View
		if err := m.Cache.ParseData(data, &views); err != nil {
			return nil, err
		}

		return views, nil
	}
	m.logger.Info("Views not found in cache")

	m.logger.Info("Fetching views from API")
	client := m.Clickup

	m.logger.Infof("Getting views from folder: %s", folder)
	views, err := client.GetViewsFromFolder(folder)
	if err != nil {
		return nil, err
	}
	m.logger.Infof("Found %d views in folder %s", len(views), folder)

	m.logger.Info("Caching views")
	m.Cache.Set("views", folder, views)

	return views, nil
}

func (m *Api) GetViewsFromList(listId string) ([]clickup.View, error) {
	m.logger.Infof("Getting views for list: %s", listId)

	data, ok := m.Cache.Get("views", listId)
	if ok {
		m.logger.Infof("Views found in cache")
		var views []clickup.View
		if err := m.Cache.ParseData(data, &views); err != nil {
			return nil, err
		}

		return views, nil
	}
	m.logger.Info("Views not found in cache")

	m.logger.Info("Fetching views from API")
	client := m.Clickup

	m.logger.Infof("Getting views from folder: %s", listId)
	views, err := client.GetViewsFromList(listId)
	if err != nil {
		return nil, err
	}
	m.logger.Infof("Found %d views in folder %s", len(views), listId)

	m.logger.Info("Caching views")
	m.Cache.Set("views", listId, views)

	return views, nil
}

func (m *Api) GetViewsFromSpace(space string) ([]clickup.View, error) {
	m.logger.Infof("Getting views for space: %s", space)

	data, ok := m.Cache.Get("views", space)
	if ok {
		m.logger.Infof("Views found in cache")
		var views []clickup.View
		if err := m.Cache.ParseData(data, &views); err != nil {
			return nil, err
		}

		return views, nil
	}
	m.logger.Info("Views not found in cache")

	m.logger.Info("Fetching views from API")
	client := m.Clickup

	m.logger.Infof("Getting views from space: %s", space)
	views, err := client.GetViewsFromSpace(space)
	if err != nil {
		return nil, err
	}
	m.logger.Infof("Found %d views in space %s", len(views), space)

	m.logger.Info("Caching views")
	m.Cache.Set("views", space, views)

	return views, nil
}
