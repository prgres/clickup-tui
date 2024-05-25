package api

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"slices"

	"github.com/charmbracelet/log"

	"github.com/prgrs/clickup/pkg/cache"
	"github.com/prgrs/clickup/pkg/clickup"
)

const (
	CacheNamespaceTeams   cache.Namespace = "teams"
	CacheNamespaceSpaces  cache.Namespace = "spaces"
	CacheNamespaceFolders cache.Namespace = "folders"
	CacheNamespaceLists   cache.Namespace = "lists"
	CacheNamespaceViews   cache.Namespace = "views"
	CacheNamespaceTasks   cache.Namespace = "tasks"
	CacheNamespaceTask    cache.Namespace = "task"
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
	key := teamId
	fallback := func() (interface{}, error) { return m.Clickup.GetSpacesFromTeam(key) }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

// Alias for GetTeams since they are the same thing
func (m *Api) GetWorkspaces() ([]clickup.Workspace, error) {
	return m.GetTeams()
}

func (m *Api) GetTeams() ([]clickup.Team, error) {
	m.logger.Debug("Getting Authorized Teams (Workspaces)")

	var data []clickup.Team
	cacheNamespace := CacheNamespaceTeams
	key := "teams"
	fallback := func() (interface{}, error) { return m.Clickup.GetTeams() }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetFolders(spaceId string) ([]clickup.Folder, error) {
	m.logger.Debug("Getting folders for a space", "space", spaceId)

	var data []clickup.Folder
	cacheNamespace := CacheNamespaceFolders
	key := spaceId
	fallback := func() (interface{}, error) { return m.Clickup.GetFolders(key) }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetLists(folderId string) ([]clickup.List, error) {
	m.logger.Debug("Getting lists for a folder", "folderId", folderId)

	var data []clickup.List
	cacheNamespace := CacheNamespaceLists
	key := folderId
	fallback := func() (interface{}, error) { return m.Clickup.GetListsFromFolder(key) }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetTask(taskId string) (clickup.Task, error) {
	m.logger.Debug("Getting a task", "taskId", taskId)

	var data clickup.Task
	cacheNamespace := CacheNamespaceTask
	key := taskId
	fallback := func() (interface{}, error) { return m.Clickup.GetTask(key) }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback); err != nil {
		return clickup.Task{}, err
	}

	return data, nil
}

func (m *Api) GetTasksFromList(listId string) ([]clickup.Task, error) {
	m.logger.Debug("Getting tasks for a list", "listId", listId)

	var data []clickup.Task
	cacheNamespace := CacheNamespaceTasks
	key := listId
	fallback := func() (interface{}, error) { return m.Clickup.GetTasksFromList(key) }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetTasksFromView(viewId string) ([]clickup.Task, error) {
	m.logger.Debug("Getting tasks for a view", "viewId", viewId)

	var data []clickup.Task
	cacheNamespace := CacheNamespaceTasks
	key := viewId
	fallback := func() (interface{}, error) { return m.Clickup.GetTasksFromView(key) }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetViewsFromFolder(folderId string) ([]clickup.View, error) {
	m.logger.Debug("Getting views for folder", "folder", folderId)

	var data []clickup.View
	cacheNamespace := CacheNamespaceViews
	key := folderId
	fallback := func() (interface{}, error) {
		v, err := m.Clickup.GetViewsFromFolder(key)
		if err != nil {
			return nil, err
		}

		return filterViews(v, []clickup.ViewType{clickup.ViewTypeList}), nil
	}

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetViewsFromList(listId string) ([]clickup.View, error) {
	m.logger.Debug("Getting views for list", "listId", listId)

	var data []clickup.View
	cacheNamespace := CacheNamespaceViews
	key := listId
	fallback := func() (interface{}, error) {
		v, err := m.Clickup.GetViewsFromList(key)
		if err != nil {
			return nil, err
		}

		return filterViews(v, []clickup.ViewType{clickup.ViewTypeList}), nil
	}

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetViewsFromSpace(spaceId string) ([]clickup.View, error) {
	m.logger.Info("Getting views for space", "spaceId", spaceId)

	var data []clickup.View
	cacheNamespace := CacheNamespaceViews
	key := spaceId
	fallback := func() (interface{}, error) {
		v, err := m.Clickup.GetViewsFromSpace(key)
		if err != nil {
			return nil, err
		}

		return filterViews(v, []clickup.ViewType{clickup.ViewTypeList}), nil
	}

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetViewsFromWorkspace(workspaceId string) ([]clickup.View, error) {
	m.logger.Debug("Getting views for workspace", "workspaceId", workspaceId)

	var data []clickup.View
	cacheNamespace := CacheNamespaceViews
	key := workspaceId
	fallback := func() (interface{}, error) {
		v, err := m.Clickup.GetViewsFromWorkspace(key)
		if err != nil {
			return nil, err
		}

		return filterViews(v, []clickup.ViewType{clickup.ViewTypeList}), nil
	}

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetList(listId string) (clickup.List, error) {
	m.logger.Debug("Getting a list", "listId", listId)

	var data clickup.List
	cacheNamespace := CacheNamespaceLists
	key := listId
	fallback := func() (interface{}, error) { return m.Clickup.GetList(key) }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback); err != nil {
		return clickup.List{}, err
	}

	return data, nil
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
			_, err := m.GetSpaces(entry.Key.String())
			if err != nil {
				m.logger.Error("Failed to invalidate spaces cache", "error", err)
			}
		case CacheNamespaceFolders:
			m.logger.Debug("Invalidating folders cache")
			_, err := m.GetFolders(entry.Key.String())
			if err != nil {
				m.logger.Error("Failed to invalidate folders cache", "error", err)
			}
		case CacheNamespaceLists:
			m.logger.Debug("Invalidating lists cache")
			_, err := m.GetLists(entry.Key.String())
			if err != nil {
				m.logger.Error("Failed to invalidate lists cache", "error", err)
			}
		case CacheNamespaceViews:
			m.logger.Debug("Invalidating views cache")
			_, err := m.GetViewsFromSpace(entry.Key.String())
			if err != nil {
				m.logger.Error("Failed to invalidate views cache", "error", err)
			}
		case CacheNamespaceTasks:
			m.logger.Debug("Invalidating tasks cache")
			_, err := m.GetTasksFromList(entry.Key.String())
			if err != nil {
				m.logger.Error("Failed to invalidate tasks cache", "error", err)
			}
		case CacheNamespaceTask:
			m.logger.Debug("Invalidating task cache")
			_, err := m.GetTask(entry.Key.String())
			if err != nil {
				m.logger.Error("Failed to invalidate task cache", "error", err)
			}
		default:
			m.logger.Debug("Invalidating cache",
				"namespace", entry.Namespace, "key", entry.Key.String())
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

func (m *Api) getFromCache(namespace cache.Namespace, key cache.Key, v interface{}) (bool, error) {
	err := m.Cache.Get(namespace, key, v)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, cache.ErrKeyNotFoundInNamespace) {
		return false, nil
	}

	return false, err
}

func (m *Api) get(cacheNamespace cache.Namespace, cacheKey cache.Key, data interface{}, fallback func() (interface{}, error)) error {
	m.logger.Debug("Getting resources", "namespace", cacheNamespace, "id", cacheKey)

	ok, err := m.getFromCache(cacheNamespace, cacheKey, data)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	m.logger.Debug("Fetching resources from API", "namespace", cacheNamespace, "id", cacheKey)
	newData, err := fallback()
	if err != nil {
		return err
	}
	m.Cache.Set(cacheNamespace, cacheKey, newData)

	// Use reflection to set the value of the data
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return fmt.Errorf("data must be a non-nil pointer")
	}

	// Ensure newData is assignable to data
	val = val.Elem()
	newVal := reflect.ValueOf(newData)
	if !newVal.Type().AssignableTo(val.Type()) {
		return fmt.Errorf("cannot assign type %T to type %T", newData, data)
	}

	val.Set(newVal)

	return nil
}
