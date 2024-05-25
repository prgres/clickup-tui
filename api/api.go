package api

import (
	"errors"
	"log/slog"
	"slices"

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
	m.logger.Debug("Getting spaces for a team", "teamId", teamId)

	var data []clickup.Space
	cacheNamespace := CacheNamespaceSpaces

	ok, err := m.getFromCache(cacheNamespace, teamId, &data)
	if err != nil {
		return nil, err
	}
	if ok {
		return data, nil
	}

	m.logger.Debug("Fetching resources from API", "namespace", cacheNamespace)
	spaces, err := m.Clickup.GetSpacesFromTeam(teamId)
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

	var data []clickup.Team
	cacheNamespace := CacheNamespaceTeams

	ok, err := m.getFromCache(cacheNamespace, "teams", &data)
	if err != nil {
		return nil, err
	}
	if ok {
		return data, nil
	}

	m.logger.Debugf("Fetching teams from API")
	teams, err := m.Clickup.GetTeams()
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Found %d teams", len(teams))

	m.Cache.Set(cacheNamespace, "teams", teams)

	return teams, nil
}

func (m *Api) GetFolders(spaceId string) ([]clickup.Folder, error) {
	m.logger.Debug("Getting folders for a space", "space", spaceId)

	var data []clickup.Folder
	cacheNamespace := CacheNamespaceFolders
	ok, err := m.getFromCache(cacheNamespace, spaceId, &data)
	if err != nil {
		return nil, err
	}
	if ok {
		return data, nil
	}

	m.logger.Debugf("Fetching folders from API")
	folders, err := m.Clickup.GetFolders(spaceId)
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Found %d folders for space: %s", len(folders), spaceId)

	m.Cache.Set(cacheNamespace, spaceId, folders)

	return folders, nil
}

func (m *Api) GetLists(folderId string) ([]clickup.List, error) {
	m.logger.Debug("Getting lists for a folder", "folderId", folderId)

	var data []clickup.List
	cacheNamespace := CacheNamespaceLists

	ok, err := m.getFromCache(cacheNamespace, folderId, &data)
	if err != nil {
		return nil, err
	}
	if ok {
		return data, nil
	}

	m.logger.Debugf("Fetching lists from API")
	lists, err := m.Clickup.GetListsFromFolder(folderId)
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Found %d lists for folder: %s", len(lists), folderId)

	m.Cache.Set(cacheNamespace, folderId, lists)

	return lists, nil
}

func (m *Api) GetTask(taskId string) (clickup.Task, error) {
	m.logger.Debug("Getting a task",
		"taskId", taskId)

	var data clickup.Task
	cacheNamespace := CacheNamespaceTask
	ok, err := m.getFromCache(cacheNamespace, taskId, &data)
	if err != nil {
		return clickup.Task{}, err
	}
	if ok {
		return data, nil
	}

	m.logger.Debug("Fetching task from API")
	task, err := m.Clickup.GetTask(taskId)
	if err != nil {
		return clickup.Task{}, err
	}
	m.logger.Debugf("Found tasks %s", taskId)

	m.Cache.Set(cacheNamespace, taskId, task)

	return task, nil
}

func (m *Api) GetTasksFromList(listId string) ([]clickup.Task, error) {
	m.logger.Debug("Getting tasks for a list", "listId", listId)

	var data []clickup.Task
	cacheNamespace := CacheNamespaceTasks

	ok, err := m.getFromCache(cacheNamespace, listId, &data)
	if err != nil {
		return nil, err
	}
	if ok {
		return data, nil
	}

	m.logger.Debug("Fetching tasks from API")
	tasks, err := m.Clickup.GetTasksFromList(listId)
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

	var data []clickup.Task
	cacheNamespace := CacheNamespaceTasks

	ok, err := m.getFromCache(cacheNamespace, viewId, &data)
	if err != nil {
		return nil, err
	}
	if ok {
		return data, nil
	}

	m.logger.Debug("Fetching tasks from API")
	tasks, err := m.Clickup.GetTasksFromView(viewId)
	if err != nil {
		return nil, err
	}
	m.logger.Debugf("Found %d tasks in view %s", len(tasks), viewId)

	m.Cache.Set(cacheNamespace, viewId, tasks)

	return tasks, nil
}

func (m *Api) GetViewsFromFolder(folderId string) ([]clickup.View, error) {
	m.logger.Debug("Getting views for folder", "folder", folderId)

	var data []clickup.View
	cacheNamespace := CacheNamespaceViews

	ok, err := m.getFromCache(cacheNamespace, folderId, &data)
	if err != nil {
		return nil, err
	}
	if ok {
		return data, nil
	}

	m.logger.Debug("Fetching views from API")
	views, err := m.Clickup.GetViewsFromFolder(folderId)
	if err != nil {
		return nil, err
	}
	views = filterViews(views, []clickup.ViewType{clickup.ViewTypeList})
	m.logger.Debugf("Found %d views in folder %s", len(views), folderId)

	m.Cache.Set(cacheNamespace, folderId, views)

	return views, nil
}

func (m *Api) GetViewsFromList(listId string) ([]clickup.View, error) {
	m.logger.Debug("Getting views for list", "listId", listId)

	var data []clickup.View
	cacheNamespace := CacheNamespaceViews

	ok, err := m.getFromCache(cacheNamespace, listId, &data)
	if err != nil {
		return nil, err
	}
	if ok {
		return data, nil
	}

	m.logger.Debug("Fetching views from API")
	views, err := m.Clickup.GetViewsFromList(listId)
	if err != nil {
		return nil, err
	}
	views = filterViews(views, []clickup.ViewType{clickup.ViewTypeList})
	m.logger.Debugf("Found %d views in folder %s", len(views), listId)

	m.Cache.Set(cacheNamespace, listId, views)

	return views, nil
}

func (m *Api) GetViewsFromSpace(spaceId string) ([]clickup.View, error) {
	m.logger.Info("Getting views for space", "spaceId", spaceId)

	var data []clickup.View
	cacheNamespace := CacheNamespaceViews

	ok, err := m.getFromCache(cacheNamespace, spaceId, &data)
	if err != nil {
		return nil, err
	}
	if ok {
		return data, nil
	}

	m.logger.Debug("Fetching views from API", "spaceId", spaceId)
	views, err := m.Clickup.GetViewsFromSpace(spaceId)
	if err != nil {
		return nil, err
	}
	views = filterViews(views, []clickup.ViewType{clickup.ViewTypeList})
	m.logger.Debugf("Found %d views in space %s", len(views), spaceId)

	m.Cache.Set(cacheNamespace, spaceId, views)

	return views, nil
}

func (m *Api) GetList(listId string) (clickup.List, error) {
	m.logger.Debug("Getting a list", "listId", listId)

	var data clickup.List
	cacheNamespace := CacheNamespaceLists

	ok, err := m.getFromCache(cacheNamespace, listId, &data)
	if err != nil {
		return clickup.List{}, err
	}
	if ok {
		return data, nil
	}

	m.logger.Debug("Fetching list from API")
	list, err := m.Clickup.GetList(listId)
	if err != nil {
		return clickup.List{}, err
	}
	m.logger.Debugf("Found list %s", listId)

	m.Cache.Set(cacheNamespace, listId, list)

	return list, nil
}

func (m *Api) GetViewsFromWorkspace(workspaceId string) ([]clickup.View, error) {
	m.logger.Debug("Getting views for workspace", "workspaceId", workspaceId)

		var data []clickup.View
	cacheNamespace := CacheNamespaceViews

	ok, err := m.getFromCache(cacheNamespace, workspaceId, &data)
	if err != nil {
		return nil, err
	}
	if ok {
		return data, nil
	}

	m.logger.Debug("Fetching views from API")
	views, err := m.Clickup.GetViewsFromWorkspace(workspaceId)
	if err != nil {
		return nil, err
	}
	views = filterViews(views, []clickup.ViewType{clickup.ViewTypeList})
	m.logger.Debugf("Found %d views in workspace %s", len(views), workspaceId)

	m.Cache.Set(cacheNamespace, workspaceId, views)

	return views, nil
}

func (m *Api) getFromCache(namespace string, key string, v interface{}) (bool, error) {
	err := m.Cache.Get(namespace, key, v)
	if err == nil {
		return true, nil
	}

	if !errors.Is(err, cache.ErrKeyNotFoundInNamespace) {
		return false, nil
	}

	return false, err
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

func filterViews(views []clickup.View, filters []clickup.ViewType) []clickup.View {
	result := []clickup.View{}
	for i := range views {
		if slices.Contains(filters, views[i].Type) {
			result = append(result, views[i])
		}
	}

	return result
}
