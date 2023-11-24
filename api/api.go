package api

import (
	"encoding/json"
	"log/slog"

	"github.com/charmbracelet/log"

	"github.com/prgrs/clickup/pkg/cache"
	"github.com/prgrs/clickup/pkg/clickup"
)

const (
	// Cache namespace
	CacheNamespaceTeams   = "teams"
	CacheNamespaceSpaces  = "spaces"
	CacheNamespaceFolders = "folders"
	CacheNamespaceLists   = "lists"
	CacheNamespaceViews   = "views"
	CacheNamespaceTasks   = "tasks"
	CacheNamespaceTask    = "task"
)

type Api struct {
	Clickup *clickup.Client
	Cache   *cache.Cache
	logger  *log.Logger
}

func NewApi(logger *log.Logger, cache *cache.Cache, token string) Api {
	log := logger.WithPrefix("Api")
	log.Debug("Initializing ClickUp client...")

	clickup := clickup.NewDefaultClientWithLogger(
		token,
		slog.New(log.WithPrefix(log.GetPrefix()+"/ClickUp")),
	)

	return Api{
		Clickup: clickup,
		logger:  log,
		Cache:   cache,
	}
}

func (m *Api) GetSpaces(teamId string) ([]clickup.Space, error) {
	m.logger.Debug("Getting spaces for a team",
		"teamId", teamId)

	cacheNamespace := CacheNamespaceSpaces

	data, ok := m.Cache.Get(cacheNamespace, teamId)
	if ok {
		var spaces []clickup.Space
		if err := m.Cache.ParseData(data, &spaces); err != nil {
			return nil, err
		}

		return spaces, nil
	}

	client := m.Clickup

	m.logger.Debugf("Fetching spaces from API")
	spaces, err := client.GetSpaces(teamId)
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Found %d spaces for team: %s", len(spaces), teamId)

	m.Cache.Set(cacheNamespace, teamId, spaces)

	return spaces, nil
}

// Alias for GetTeams since they are the same thing
func (m *Api) GetWorkspaces() ([]clickup.Workspace, error) {
	return m.GetTeams()
}

func (m *Api) GetTeams() ([]clickup.Team, error) {
	m.logger.Debug("Getting Authorized Teams (Workspaces)")

	cacheNamespace := CacheNamespaceTeams
	data, ok := m.Cache.Get(cacheNamespace, "teams")
	if ok {
		var teams []clickup.Team
		if err := m.Cache.ParseData(data, &teams); err != nil {
			return nil, err
		}

		return teams, nil
	}

	client := m.Clickup

	m.logger.Debugf("Fetching teams from API")
	teams, err := client.GetTeams()
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Found %d teams", len(teams))

	m.Cache.Set(cacheNamespace, "teams", teams)

	return teams, nil
}

func (m *Api) GetFolders(spaceId string) ([]clickup.Folder, error) {
	m.logger.Debug("Getting folders for a space",
		"space", spaceId)

	cacheNamespace := CacheNamespaceFolders
	data, ok := m.Cache.Get(cacheNamespace, spaceId)
	if ok {
		var folders []clickup.Folder
		if err := m.Cache.ParseData(data, &folders); err != nil {
			return nil, err
		}

		return folders, nil
	}

	client := m.Clickup

	m.logger.Debugf("Fetching folders from API")
	folders, err := client.GetFolders(spaceId)
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Found %d folders for space: %s", len(folders), spaceId)

	m.Cache.Set(cacheNamespace, spaceId, folders)

	return folders, nil
}

func (m *Api) GetLists(folderId string) ([]clickup.List, error) {
	m.logger.Debug("Getting lists for a folder",
		"folderId", folderId)

	cacheNamespace := CacheNamespaceLists
	data, ok := m.Cache.Get(cacheNamespace, folderId)
	if ok {
		var lists []clickup.List
		if err := m.Cache.ParseData(data, &lists); err != nil {
			return nil, err
		}

		return lists, nil
	}

	client := m.Clickup

	m.logger.Debugf("Fetching lists from API")
	lists, err := client.GetListsFromFolder(folderId)
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Found %d lists for folder: %s", len(lists), folderId)

	for _, f := range lists {
		j, _ := json.MarshalIndent(f, "", "  ")
		m.logger.Debugf("%s", j)
	}

	m.Cache.Set(cacheNamespace, folderId, lists)

	return lists, nil
}

func (m *Api) GetTask(taskId string) (clickup.Task, error) {
	m.logger.Debug("Getting a task",
		"taskId", taskId)

	cacheNamespace := CacheNamespaceTask
	data, ok := m.Cache.Get(cacheNamespace, taskId)
	if ok {
		var task clickup.Task
		if err := m.Cache.ParseData(data, &task); err != nil {
			return clickup.Task{}, err
		}

		return task, nil
	}

	client := m.Clickup

	m.logger.Debug("Fetching task from API")
	task, err := client.GetTask(taskId)
	if err != nil {
		return clickup.Task{}, err
	}
	m.logger.Debugf("Found tasks %s", taskId)

	m.Cache.Set(cacheNamespace, taskId, task)

	return task, nil
}

func (m *Api) GetTasksFromList(listId string) ([]clickup.Task, error) {
	m.logger.Debug("Getting tasks for a list",
		"listId", listId)

	cacheNamespace := CacheNamespaceTasks
	data, ok := m.Cache.Get(cacheNamespace, listId)
	if ok {
		var tasks []clickup.Task
		if err := m.Cache.ParseData(data, &tasks); err != nil {
			return nil, err
		}

		return tasks, nil
	}

	client := m.Clickup

	m.logger.Debug("Fetching tasks from API")
	tasks, err := client.GetTasksFromList(listId)
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Found %d tasks in list %s", len(tasks), listId)

	m.Cache.Set(cacheNamespace, listId, tasks)

	return tasks, nil
}
func (m *Api) GetTasksFromView(viewId string) ([]clickup.Task, error) {
	m.logger.Debug("Getting tasks for a view",
		"viewId", viewId)

	cacheNamespace := CacheNamespaceTasks
	data, ok := m.Cache.Get(cacheNamespace, viewId)
	if ok {
		var tasks []clickup.Task
		if err := m.Cache.ParseData(data, &tasks); err != nil {
			return nil, err
		}

		return tasks, nil
	}

	client := m.Clickup

	m.logger.Debug("Fetching tasks from API")
	tasks, err := client.GetTasksFromView(viewId)
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Found %d tasks in view %s", len(tasks), viewId)

	m.Cache.Set(cacheNamespace, viewId, tasks)

	return tasks, nil
}

func (m *Api) GetViewsFromFolder(folderId string) ([]clickup.View, error) {
	m.logger.Debug("Getting views for folder",
		"folder", folderId)

	cacheNamespace := CacheNamespaceViews
	data, ok := m.Cache.Get(cacheNamespace, folderId)
	if ok {
		var views []clickup.View
		if err := m.Cache.ParseData(data, &views); err != nil {
			return nil, err
		}

		return views, nil
	}

	client := m.Clickup

	m.logger.Debug("Fetching views from API")
	views, err := client.GetViewsFromFolder(folderId)
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Found %d views in folder %s", len(views), folderId)

	m.Cache.Set(cacheNamespace, folderId, views)

	return views, nil
}

func (m *Api) GetViewsFromList(listId string) ([]clickup.View, error) {
	m.logger.Debug("Getting views for list",
		"listId", listId)

	cacheNamespace := CacheNamespaceViews
	data, ok := m.Cache.Get(cacheNamespace, listId)
	if ok {
		var views []clickup.View
		if err := m.Cache.ParseData(data, &views); err != nil {
			return nil, err
		}

		return views, nil
	}

	client := m.Clickup

	m.logger.Debug("Fetching views from API")
	views, err := client.GetViewsFromList(listId)
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Found %d views in folder %s", len(views), listId)

	m.Cache.Set(cacheNamespace, listId, views)

	return views, nil
}

func (m *Api) GetViewsFromSpace(spaceId string) ([]clickup.View, error) {
	m.logger.Info("Getting views for space",
		"spaceId", spaceId)

	cacheNamespace := CacheNamespaceViews
	data, ok := m.Cache.Get(cacheNamespace, spaceId)
	if ok {
		var views []clickup.View
		if err := m.Cache.ParseData(data, &views); err != nil {
			return nil, err
		}

		return views, nil
	}

	client := m.Clickup

	m.logger.Debug("Fetching views from API",
		"spaceId", spaceId)
	views, err := client.GetViewsFromSpace(spaceId)
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Found %d views in space %s", len(views), spaceId)

	m.Cache.Set(cacheNamespace, spaceId, views)

	return views, nil
}

//nolint:unused
func (m *Api) getFromCache(namespace string, key string, v interface{}) (bool, error) {
	data, ok := m.Cache.Get(namespace, key)
	if !ok {
		return false, nil
	}

	if err := m.Cache.ParseData(data, &v); err != nil {
		return false, err
	}

	return true, nil
}

func (m *Api) InvalidateCache() error {
	m.logger.Debug("Invalidating cache")

	entries := m.Cache.GetEntries()
	m.logger.Debug("Found cache entries", "count", len(entries))

	if err := m.Cache.Invalidate(); err != nil {
		m.logger.Error("Failed to invalidate cache", "error", err)
		return err
	}

	for _, entry := range entries {
		switch entry.Namespace {
		case CacheNamespaceTeams:
			m.logger.Debug("Invalidating teams cache")
			m.GetTeams()
		case CacheNamespaceSpaces:
			m.logger.Debug("Invalidating spaces cache")
			m.GetSpaces(entry.Key)
		case CacheNamespaceFolders:
			m.logger.Debug("Invalidating folders cache")
			m.GetFolders(entry.Key)
		case CacheNamespaceLists:
			m.logger.Debug("Invalidating lists cache")
			m.GetLists(entry.Key)
		case CacheNamespaceViews:
			m.logger.Debug("Invalidating views cache")
			m.GetViewsFromSpace(entry.Key)
		case CacheNamespaceTasks:
			m.logger.Debug("Invalidating tasks cache")
			m.GetTasksFromList(entry.Key)
		case CacheNamespaceTask:
			m.logger.Debug("Invalidating task cache")
			m.GetTask(entry.Key)
		default:
			m.logger.Debug("Invalidating cache",
				"namespace", entry.Namespace, "key", entry.Key)
		}
	}
	return nil
}
