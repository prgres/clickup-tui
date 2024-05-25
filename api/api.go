package api

import (
	"errors"
	"fmt"
	"log/slog"
	"time"
	"reflect"
	"slices"

	"github.com/charmbracelet/log"

	"github.com/prgrs/clickup/pkg/cache"
	"github.com/prgrs/clickup/pkg/clickup"
)

const (
	// Cache namespace
	CacheNamespaceTeams     cache.Namespace = "teams"
	CacheNamespaceSpaces    cache.Namespace = "spaces"
	CacheNamespaceFolders   cache.Namespace = "folders"
	CacheNamespaceLists     cache.Namespace = "lists"
	CacheNamespaceViews     cache.Namespace = "views"
	CacheNamespaceTask      cache.Namespace = "tasks"
	CacheNamespaceTasksList cache.Namespace = "tasks-list"
	CacheNamespaceTasksView cache.Namespace = "tasks-view"
)

const (
	TTL                      = 10
	GarbageCollectorInterval = 5
)

type Api struct {
	Clickup *clickup.Client
	Cache   *cache.Cache
	logger  *log.Logger

	gcCloseChan chan struct{}
	interval    time.Duration
}

func NewApi(logger *log.Logger, cache *cache.Cache, token string) Api {
	log := logger.WithPrefix("Api")
	log.Debug("Initializing ClickUp client...")

	clickup := clickup.NewDefaultClientWithLogger(
		token,
		slog.New(log.WithPrefix(log.GetPrefix()+"/ClickUp")),
	)

	api := Api{
		Clickup:     clickup,
		logger:      log,
		Cache:       cache,
		gcCloseChan: make(chan struct{}),
		interval:    GarbageCollectorInterval * time.Second,
	}

	go api.garbageCollector()
	return api
}

func (m *Api) garbageCollector() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.logger.Debug("Garbage Collector: starting")
			// now := time.Now().Unix()
			if err := m.InvalidateCache(); err != nil {
				panic(err)
			}
		case <-m.gcCloseChan:
			return
		}
	}
}

func (m *Api) GetSpaces(teamId string) ([]clickup.Space, error) {
	m.logger.Debug("Getting spaces for a team", "teamId", teamId)

	var data []clickup.Space
	cacheNamespace := CacheNamespaceSpaces
	key := teamId
	fallback := func() (interface{}, error) { return m.Clickup.GetSpacesFromTeam(key) }

	if err := m.get(cacheNamespace, key, &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) syncSpaces(entry cache.Entry) error {
	m.logger.Debug("Sync spaces for a team", "teamId", entry.Key)

	client := m.Clickup

	m.logger.Debugf("Fetching spaces from API")
	data, err := client.GetSpaces(entry.Key)
	if err != nil {
		return err
	}
	m.logger.Debugf("Found %d spaces for team: %s", len(data), entry.Key)

	entry.UpdatedTimestamp = time.Now().Unix()
	entry.Value = data
	m.Cache.Update(entry)

	return nil
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

	if err := m.get(cacheNamespace, key, &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) syncTeams(entry cache.Entry) error {
	m.logger.Debug("Sync Authorized Teams (Workspaces)")

	client := m.Clickup

	m.logger.Debugf("Fetching teams from API")
	data, err := client.GetTeams()
	if err != nil {
		return err
	}
	m.logger.Debugf("Found %d teams", len(data))

	entry.UpdatedTimestamp = time.Now().Unix()
	entry.Value = data
	m.Cache.Update(entry)

	return nil
}

func (m *Api) GetFolders(spaceId string) ([]clickup.Folder, error) {
	m.logger.Debug("Getting folders for a space", "space", spaceId)

	var data []clickup.Folder
	cacheNamespace := CacheNamespaceFolders
	key := spaceId
	fallback := func() (interface{}, error) { return m.Clickup.GetFolders(key) }

	if err := m.get(cacheNamespace, key, &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) syncFolders(entry cache.Entry) error {
	spaceId := entry.Key
	m.logger.Debug("Sync folders for a space", "space", spaceId)

	client := m.Clickup

	m.logger.Debugf("Fetching folders from API")
	data, err := client.GetFolders(spaceId)
	if err != nil {
		return err
	}
	m.logger.Debugf("Found %d folders for space: %s", len(data), spaceId)

	entry.UpdatedTimestamp = time.Now().Unix()
	entry.Value = data
	m.Cache.Update(entry)

	return nil
}

func (m *Api) GetLists(folderId string) ([]clickup.List, error) {
	m.logger.Debug("Getting lists for a folder", "folderId", folderId)

	var data []clickup.List
	cacheNamespace := CacheNamespaceLists
	key := folderId
	fallback := func() (interface{}, error) { return m.Clickup.GetListsFromFolder(key) }

	if err := m.get(cacheNamespace, key, &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) syncLists(entry cache.Entry) error {
	folderId := entry.Key
	m.logger.Debug("Getting lists for a folder", "folderId", folderId)

	client := m.Clickup

	m.logger.Debugf("Fetching lists from API")
	data, err := client.GetListsFromFolder(folderId)
	if err != nil {
		return err
	}
	m.logger.Debugf("Found %d lists for folder: %s", len(data), folderId)

	entry.UpdatedTimestamp = time.Now().Unix()
	entry.Value = data
	m.Cache.Update(entry)

	return nil
}

func (m *Api) GetTask(taskId string) (clickup.Task, error) {
	m.logger.Debug("Getting a task", "taskId", taskId)

	var data clickup.Task
	cacheNamespace := CacheNamespaceTask
	key := taskId
	fallback := func() (interface{}, error) { return m.Clickup.GetTask(key) }

	if err := m.get(cacheNamespace, key, &data, fallback); err != nil {
		return clickup.Task{}, err
	}

	return data, nil
}

func (m *Api) syncTask(entry cache.Entry) error {
	taskId := entry.Key
	m.logger.Debug("Sync a task", "taskId", taskId)

	client := m.Clickup

	m.logger.Debug("Fetching task from API")
	data, err := client.GetTask(taskId)
	if err != nil {
		return err
	}
	m.logger.Debug("Found task", "task", taskId)

	entry.UpdatedTimestamp = time.Now().Unix()
	entry.Value = data
	m.Cache.Update(entry)

	return nil
}

func (m *Api) GetTasksFromList(listId string) ([]clickup.Task, error) {
	m.logger.Debug("Getting tasks for a list", "listId", listId)

	var data []clickup.Task
	cacheNamespace := CacheNamespaceTasks
	key := listId
	fallback := func() (interface{}, error) { return m.Clickup.GetTasksFromList(key) }

	if err := m.get(cacheNamespace, key, &data, fallback); err != nil {
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

	if err := m.get(cacheNamespace, key, &data, fallback); err != nil {
		return nil, err
	}

	return data, nil
}

func (m *Api) syncTasksFromView(entry cache.Entry) error {
	viewId := entry.Key
	m.logger.Debug("Sync tasks for a view", "viewId", viewId)

	client := m.Clickup

	m.logger.Debug("Fetching tasks from API")
	data, err := client.GetTasksFromView(viewId)
	if err != nil {
		return err
	}
	m.logger.Debugf("Found %d tasks in view %s", len(data), viewId)

	entry.UpdatedTimestamp = time.Now().Unix()
	entry.Value = data
	m.Cache.Update(entry)

	return nil
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

	if err := m.get(cacheNamespace, key, &data, fallback); err != nil {
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

	if err := m.get(cacheNamespace, key, &data, fallback); err != nil {
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

	if err := m.get(cacheNamespace, key, &data, fallback); err != nil {
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

	if err := m.get(cacheNamespace, key, &data, fallback); err != nil {
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

	if err := m.get(cacheNamespace, key, &data, fallback); err != nil {
		return clickup.List{}, err
	}

	return data, nil
}

func (m *Api) InvalidateCache() error {
	m.logger.Debug("Invalidating cache")

	entries := m.Cache.GetEntries()
	m.logger.Debug("Found cache entries", "count", len(entries))

	// if err := m.Cache.Invalidate(); err != nil {
	// 	m.logger.Error("Failed to invalidate cache", "error", err)
	// 	return err
	// }

	now := time.Now().Unix()
	for _, entry := range entries {
		if entry.UpdatedTimestamp+TTL > now {
			continue
		}

		var err error

		m.logger.Debug("Invalidating cache", "namespace", entry.Namespace)
		switch entry.Namespace {
		case CacheNamespaceTeams:
			err = m.syncTeams(entry)
		case CacheNamespaceSpaces:
			err = m.syncSpaces(entry)
		case CacheNamespaceFolders:
			err = m.syncFolders(entry)
		case CacheNamespaceLists:
			err = m.syncLists(entry)

		// TODO:
		// case CacheNamespaceViews:
		// 	m.logger.Debug("Invalidating views cache")
		// 	_, err := m.GetViewsFromSpace(entry.Key)
		// 	if err != nil {
		// 		m.logger.Error("Failed to invalidate views cache", "error", err)
		// 	}
		// case CacheNamespaceTasksList:
		// 	m.logger.Debug("Invalidating tasks cache")
		// 	_, err := m.GetTasksFromList(entry.Key)
		// 	if err != nil {
		// 		m.logger.Error("Failed to invalidate tasks cache", "error", err)
		// 	}
		case CacheNamespaceTasksView:
			err = m.syncTasksFromView(entry)
		case CacheNamespaceTask:
			err = m.syncTask(entry)
		}

		if err != nil {
			m.logger.Error("Failed to invalidate task cache", "error", err)
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

func (m *Api) getFromCache(namespace string, key string, v interface{}) (bool, error) {
	err := m.Cache.Get(namespace, key, v)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, cache.ErrKeyNotFoundInNamespace) {
		return false, nil
	}

	return false, err
}

func (m *Api) get(cacheNamespace string, id string, data interface{}, fallback func() (interface{}, error)) error {
	m.logger.Debug("Getting resources", "namespace", cacheNamespace, "id", id)

	ok, err := m.getFromCache(cacheNamespace, id, data)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	m.logger.Debug("Fetching resources from API", "namespace", cacheNamespace, "id", id)
	newData, err := fallback()
	if err != nil {
		return err
	}
	m.Cache.Set(cacheNamespace, id, newData)

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
