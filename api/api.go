package api

import (
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

	if len(spaces) > 0 {
		m.Cache.Set(cacheNamespace, teamId, spaces)
	}

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

	if len(teams) == 0 {
		m.Cache.Set(cacheNamespace, "teams", teams)
	}

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

	if len(folders) > 0 {
		m.Cache.Set(cacheNamespace, spaceId, folders)
	}

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

	if len(lists) > 0 {
		m.Cache.Set(cacheNamespace, folderId, lists)
	}

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

	if len(tasks) > 0 {
		m.Cache.Set(cacheNamespace, listId, tasks)
	}

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

	if len(tasks) > 0 {
		m.Cache.Set(cacheNamespace, viewId, tasks)
	}

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

	if len(views) > 0 {
		m.Cache.Set(cacheNamespace, folderId, views)
	}

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

	if len(views) > 0 {
		m.Cache.Set(cacheNamespace, listId, views)
	}

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

	if len(views) > 0 {
		m.Cache.Set(cacheNamespace, spaceId, views)
	}

	return views, nil
}

func (m *Api) GetList(listId string) (clickup.List, error) {
	m.logger.Debug("Getting a list",
		"listId", listId)
	cacheNamespace := CacheNamespaceLists
	data, ok := m.Cache.Get(cacheNamespace, listId)
	if ok {
		var list clickup.List
		if err := m.Cache.ParseData(data, &list); err != nil {
			return clickup.List{}, err
		}
		return list, nil
	}
	client := m.Clickup
	m.logger.Debug("Fetching list from API")

	list, err := client.GetList(listId)
	if err != nil {
		return clickup.List{}, err
	}
	m.logger.Debugf("Found list %s", listId)

	m.Cache.Set(cacheNamespace, listId, list)

	return list, nil
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
			_, err := m.GetTeams()
			if err != nil {
				m.logger.Error("Failed to invalidate teams cache", "error", err)
			}
		case CacheNamespaceSpaces:
			m.logger.Debug("Invalidating spaces cache")
			_, err := m.GetSpaces(entry.Key)
			if err != nil {
				m.logger.Error("Failed to invalidate spaces cache", "error", err)
			}
		case CacheNamespaceFolders:
			m.logger.Debug("Invalidating folders cache")
			_, err := m.GetFolders(entry.Key)
			if err != nil {
				m.logger.Error("Failed to invalidate folders cache", "error", err)
			}
		case CacheNamespaceLists:
			m.logger.Debug("Invalidating lists cache")
			_, err := m.GetLists(entry.Key)
			if err != nil {
				m.logger.Error("Failed to invalidate lists cache", "error", err)
			}
		case CacheNamespaceViews:
			m.logger.Debug("Invalidating views cache")
			_, err := m.GetViewsFromSpace(entry.Key)
			if err != nil {
				m.logger.Error("Failed to invalidate views cache", "error", err)
			}
		case CacheNamespaceTasks:
			m.logger.Debug("Invalidating tasks cache")
			_, err := m.GetTasksFromList(entry.Key)
			if err != nil {
				m.logger.Error("Failed to invalidate tasks cache", "error", err)
			}
		case CacheNamespaceTask:
			m.logger.Debug("Invalidating task cache")
			_, err := m.GetTask(entry.Key)
			if err != nil {
				m.logger.Error("Failed to invalidate task cache", "error", err)
			}
		default:
			m.logger.Debug("Invalidating cache",
				"namespace", entry.Namespace, "key", entry.Key)
		}
	}
	return nil
}
