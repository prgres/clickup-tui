package api

import (
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"time"

	"github.com/charmbracelet/log"
	"golang.org/x/sync/errgroup"

	"github.com/prgrs/clickup/pkg/cache"
	"github.com/prgrs/clickup/pkg/clickup"
)

const (
	CacheNamespaceTeams          cache.Namespace = "teams"
	CacheNamespaceSpaces         cache.Namespace = "spaces"
	CacheNamespaceFolders        cache.Namespace = "folders"
	CacheNamespaceLists          cache.Namespace = "lists"
	CacheNamespaceListsFolder    cache.Namespace = "lists-folder"
	CacheNamespaceViewsWorkspace cache.Namespace = "views-workspace"
	CacheNamespaceViewsSpace     cache.Namespace = "views-space"
	CacheNamespaceViewsFolder    cache.Namespace = "views-folder"
	CacheNamespaceViewsList      cache.Namespace = "views-list"
	CacheNamespaceTasks          cache.Namespace = "tasks"
	CacheNamespaceTasksList      cache.Namespace = "tasks-list"
	CacheNamespaceTasksView      cache.Namespace = "tasks-view"

	SyncInterval = 1
)

type Api struct {
	Clickup   *clickup.Client
	Cache     *cache.Cache
	logger    *log.Logger
	closeChan chan struct{}
	interval  time.Duration
}

func NewApi(logger *log.Logger, cache *cache.Cache, token string) *Api {
	log := logger.WithPrefix("Api")
	log.Debug("Initializing ClickUp client...")

	clickup := clickup.NewDefaultClientWithLogger(
		token,
		slog.New(log.WithPrefix(log.GetPrefix()+"/ClickUp")),
	)

	a := Api{
		Clickup:   clickup,
		logger:    log,
		Cache:     cache,
		interval:  SyncInterval * time.Second,
		closeChan: make(chan struct{}),
	}

	go a.sync()

	return &a
}

func (m *Api) sync() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := m.Sync(); err != nil {
				m.logger.Error(err)
			}

		case <-m.closeChan:
			return
		}
	}
}

func (m *Api) Close() {
	m.closeChan <- struct{}{}
}

func (m *Api) GetSpaces(teamId string) ([]clickup.Space, error) {
	return m.getSpaces(true, teamId)
}

func (m *Api) syncSpaces(teamId string) ([]clickup.Space, error) {
	return m.getSpaces(false, teamId)
}

func (m *Api) getSpaces(cached bool, teamId string) ([]clickup.Space, error) {
	m.logger.Debug("Getting spaces for a team", "teamId", teamId)

	var data []clickup.Space
	cacheNamespace := CacheNamespaceSpaces
	key := teamId
	fallback := func() (interface{}, error) { return m.Clickup.GetSpacesFromTeam(key) }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback, cached); err != nil {
		return nil, err
	}

	return data, nil
}

// Alias for GetTeams since they are the same thing
func (m *Api) GetWorkspaces() ([]clickup.Workspace, error) {
	return m.GetTeams()
}

func (m *Api) GetTeams() ([]clickup.Team, error) {
	return m.getTeams(true)
}

func (m *Api) syncTeams() ([]clickup.Team, error) {
	return m.getTeams(false)
}

func (m *Api) getTeams(cached bool) ([]clickup.Team, error) {
	m.logger.Debug("Getting Authorized Teams (Workspaces)")

	var data []clickup.Team
	cacheNamespace := CacheNamespaceTeams
	key := "teams"
	fallback := func() (interface{}, error) { return m.Clickup.GetTeams() }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback, cached); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetFolders(spaceId string) ([]clickup.Folder, error) {
	return m.getFolders(true, spaceId)
}

func (m *Api) syncFolders(spaceId string) ([]clickup.Folder, error) {
	return m.getFolders(false, spaceId)
}

func (m *Api) getFolders(cached bool, spaceId string) ([]clickup.Folder, error) {
	m.logger.Debug("Getting folders for a space", "space", spaceId)

	var data []clickup.Folder
	cacheNamespace := CacheNamespaceFolders
	key := spaceId
	fallback := func() (interface{}, error) { return m.Clickup.GetFolders(key) }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback, cached); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetListsFromFolder(folderId string) ([]clickup.List, error) {
	return m.getListsFromFolder(true, folderId)
}

func (m *Api) syncListsFromFolder(folderId string) ([]clickup.List, error) {
	return m.getListsFromFolder(false, folderId)
}

func (m *Api) getListsFromFolder(cached bool, folderId string) ([]clickup.List, error) {
	m.logger.Debug("Getting lists for a folder", "folderId", folderId)

	var data []clickup.List
	cacheNamespace := CacheNamespaceListsFolder
	key := folderId
	fallback := func() (interface{}, error) { return m.Clickup.GetListsFromFolder(key) }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback, cached); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetTask(taskId string) (clickup.Task, error) {
	return m.getTask(true, taskId)
}

func (m *Api) SyncTask(taskId string) (clickup.Task, error) {
	return m.getTask(false, taskId)
}

func (m *Api) getTask(cached bool, taskId string) (clickup.Task, error) {
	m.logger.Debug("Getting a task", "taskId", taskId)

	var data clickup.Task
	cacheNamespace := CacheNamespaceTasks
	key := taskId
	fallback := func() (interface{}, error) { return m.Clickup.GetTask(key) }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback, cached); err != nil {
		return clickup.Task{}, err
	}

	return data, nil
}

func (m *Api) GetTasksFromList(listId string) ([]clickup.Task, error) {
	return m.getTasksFromList(true, listId)
}

func (m *Api) syncTasksFromList(listId string) ([]clickup.Task, error) {
	return m.getTasksFromList(false, listId)
}

func (m *Api) getTasksFromList(cached bool, listId string) ([]clickup.Task, error) {
	m.logger.Debug("Getting tasks for a list", "listId", listId)

	var data []clickup.Task
	cacheNamespace := CacheNamespaceTasksList
	key := listId
	fallback := func() (interface{}, error) { return m.Clickup.GetTasksFromList(key) }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback, cached); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetTasksFromView(viewId string) ([]clickup.Task, error) {
	return m.getTasksFromView(true, viewId)
}

func (m *Api) syncTasksFromView(viewId string) ([]clickup.Task, error) {
	return m.getTasksFromView(false, viewId)
}

func (m *Api) getTasksFromView(cached bool, viewId string) ([]clickup.Task, error) {
	m.logger.Debug("Getting tasks for a view", "viewId", viewId)

	var data []clickup.Task
	cacheNamespace := CacheNamespaceTasksView
	key := viewId
	fallback := func() (interface{}, error) { return m.Clickup.GetTasksFromView(key) }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback, cached); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetViewsFromFolder(folderId string) ([]clickup.View, error) {
	return m.getViewsFromFolder(true, folderId)
}

func (m *Api) syncViewsFromFolder(folderId string) ([]clickup.View, error) {
	return m.getViewsFromFolder(false, folderId)
}

func (m *Api) getViewsFromFolder(cached bool, folderId string) ([]clickup.View, error) {
	m.logger.Debug("Getting views for folder", "folder", folderId)

	var data []clickup.View
	cacheNamespace := CacheNamespaceViewsFolder
	key := folderId
	fallback := func() (interface{}, error) {
		v, err := m.Clickup.GetViewsFromFolder(key)
		if err != nil {
			return nil, err
		}

		return filterViews(v, []clickup.ViewType{clickup.ViewTypeList}), nil
	}

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback, cached); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetViewsFromList(listId string) ([]clickup.View, error) {
	return m.getViewsFromList(true, listId)
}

func (m *Api) syncViewsFromList(listId string) ([]clickup.View, error) {
	return m.getViewsFromList(false, listId)
}

func (m *Api) getViewsFromList(cached bool, listId string) ([]clickup.View, error) {
	m.logger.Debug("Getting views for list", "listId", listId)

	var data []clickup.View
	cacheNamespace := CacheNamespaceViewsList
	key := listId
	fallback := func() (interface{}, error) {
		v, err := m.Clickup.GetViewsFromList(key)
		if err != nil {
			return nil, err
		}

		return filterViews(v, []clickup.ViewType{clickup.ViewTypeList}), nil
	}

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback, cached); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetViewsFromSpace(spaceId string) ([]clickup.View, error) {
	return m.getViewsFromSpace(true, spaceId)
}

func (m *Api) syncViewsFromSpace(spaceId string) ([]clickup.View, error) {
	return m.getViewsFromSpace(false, spaceId)
}

func (m *Api) getViewsFromSpace(cached bool, spaceId string) ([]clickup.View, error) {
	m.logger.Info("Getting views for space", "spaceId", spaceId)

	var data []clickup.View
	cacheNamespace := CacheNamespaceViewsSpace
	key := spaceId
	fallback := func() (interface{}, error) {
		v, err := m.Clickup.GetViewsFromSpace(key)
		if err != nil {
			return nil, err
		}

		return filterViews(v, []clickup.ViewType{clickup.ViewTypeList}), nil
	}

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback, cached); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetViewsFromWorkspace(workspaceId string) ([]clickup.View, error) {
	return m.getViewsFromWorkspace(true, workspaceId)
}

func (m *Api) syncViewsFromWorkspace(workspaceId string) ([]clickup.View, error) {
	return m.getViewsFromWorkspace(false, workspaceId)
}

func (m *Api) getViewsFromWorkspace(cached bool, workspaceId string) ([]clickup.View, error) {
	m.logger.Debug("Getting views for workspace", "workspaceId", workspaceId)

	var data []clickup.View
	cacheNamespace := CacheNamespaceViewsWorkspace
	key := workspaceId
	fallback := func() (interface{}, error) {
		v, err := m.Clickup.GetViewsFromWorkspace(key)
		if err != nil {
			return nil, err
		}

		return filterViews(v, []clickup.ViewType{clickup.ViewTypeList}), nil
	}

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback, cached); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) GetList(listId string) (clickup.List, error) {
	return m.getList(true, listId)
}

func (m *Api) syncList(listId string) (clickup.List, error) {
	return m.getList(false, listId)
}

func (m *Api) getList(cached bool, listId string) (clickup.List, error) {
	m.logger.Debug("Getting a list", "listId", listId)

	var data clickup.List
	cacheNamespace := CacheNamespaceLists
	key := listId
	fallback := func() (interface{}, error) { return m.Clickup.GetList(key) }

	if err := m.get(cacheNamespace, cache.Key(key), &data, fallback, cached); err != nil {
		return clickup.List{}, err
	}

	return data, nil
}

func (m *Api) Sync() error {
	m.logger.Debug("Sync API")

	entries := m.Cache.GetEntries()
	errgroup := new(errgroup.Group)

	for _, entry := range entries {
		func(entry cache.Entry) {
			errgroup.Go(func() error {
				if !entry.Stale {
					return nil
				}

				m.logger.Debug("Invalidating cache", "entry", entry.Id())

				var err error
				key := entry.Key.String()

				switch entry.Namespace {
				case CacheNamespaceTeams:
					_, err = m.syncTeams()
				case CacheNamespaceSpaces:
					_, err = m.syncSpaces(key)
				case CacheNamespaceFolders:
					_, err = m.syncFolders(key)
				case CacheNamespaceListsFolder:
					_, err = m.syncListsFromFolder(key)
				case CacheNamespaceLists:
					_, err = m.syncList(key)
				case CacheNamespaceViewsWorkspace:
					_, err = m.syncViewsFromWorkspace(key)
				case CacheNamespaceViewsSpace:
					_, err = m.syncViewsFromSpace(key)
				case CacheNamespaceViewsFolder:
					_, err = m.syncViewsFromFolder(key)
				case CacheNamespaceViewsList:
					_, err = m.syncViewsFromList(key)
				case CacheNamespaceTasksList:
					_, err = m.syncTasksFromList(key)
				case CacheNamespaceTasksView:
					_, err = m.syncTasksFromView(key)
				case CacheNamespaceTasks:
					_, err = m.SyncTask(key)
				default:
					m.logger.Warn("Removing cache entry due to invalid namespace", "entry", entry.Id(), "namespace", entry.Namespace)
				}

				if err != nil {
					return fmt.Errorf("failed to invalidate task cache %e", err)
				}
				return nil
			})
		}(entry)
	}

	if err := errgroup.Wait(); err != nil {
		m.logger.Error(err)
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

func (m *Api) get(cacheNamespace cache.Namespace, cacheKey cache.Key, data interface{}, fallback func() (interface{}, error), cache bool) error {
	m.logger.Debug("Getting resources", "namespace", cacheNamespace, "id", cacheKey)

	if cache {
		ok, err := m.getFromCache(cacheNamespace, cacheKey, data)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
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

func (m *Api) update(cacheNamespace cache.Namespace, cacheKey cache.Key, data interface{}, updateFn func() (interface{}, error)) error {
	m.logger.Debug("Updating resources", "namespace", cacheNamespace, "id", cacheKey)

	m.logger.Debug("Fetching resources from API", "namespace", cacheNamespace, "id", cacheKey)

	newData, err := updateFn()
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

func (m *Api) UpdateTask(task clickup.Task) (clickup.Task, error) {
	r := clickup.RequestPutTask{
		Id:          task.Id,
		Name:        task.Name,
		Description: task.Description,
		Status:      task.Status.Status,
		Points:      task.Points,
	}

	t, err := m.Clickup.UpdateTask(r)
	if err != nil {
		return clickup.Task{}, err
	}

	return m.SyncTask(t.Id)
}
