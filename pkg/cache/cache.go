package cache

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	TTL                      = 3 * 60  // 3m
	StaleInterval            = 30 * 60 // 30m
	GarbageCollectorInterval = 3       // 5s
)

type Entry struct {
	Value             interface{} `json:"value"`
	Key               string      `json:"key"`
	Namespace         Namespace   `json:"namespace"`
	Stale             bool        `json:"stale"`
	AccessedTimestamp int64       `json:"accessedtimestamp"`
	CreatedTimestamp  int64       `json:"createdtimestamp"`
	UpdatedTimestamp  int64       `json:"updatedTimestamp"`
}

type Key string

type Data map[string]Entry

type Namespace string

type Cache struct {
	logger    *slog.Logger
	data      map[Namespace]Data
	closeChan chan struct{}
	path      string
	interval  time.Duration
	mutex     sync.RWMutex
}

func NewCache(logger *slog.Logger, path string) *Cache {
	c := &Cache{
		path:      path,
		data:      map[Namespace]Data{},
		logger:    logger,
		interval:  GarbageCollectorInterval * time.Second,
		closeChan: make(chan struct{}),
	}

	go c.garbageCollector()

	return c
}

func (c *Cache) garbageCollector() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.logger.Debug("Garbage Collector: starting")
			now := time.Now().Unix()
			entries := c.GetEntries()
			for _, entry := range entries {
				entryId := string(entry.Namespace) + "/" + entry.Key
				if !entry.Stale {
					if now > entry.AccessedTimestamp+StaleInterval {
						c.logger.Debug("Garbage Collector: marking as stale", "entry", entryId)
						entry.Stale = true
						c.Update(entry)
					}

					continue
				}

				if now > entry.UpdatedTimestamp+TTL {
					c.logger.Debug("Garbage Collector: deleting stale", "entry", entryId)
					c.Delete(entry.Namespace, entry.Key)
				}
			}

		case <-c.closeChan:
			return
		}
	}
}

func (c *Cache) Delete(namespace Namespace, key string) {
	c.logger.Debug("Removing", "entry", string(namespace)+"/"+key)
	c.mutex.Lock()
	delete(c.data[namespace], key)
	c.mutex.Unlock()
}

func (c *Cache) Update(entry Entry) {
	c.logger.Debug("Updating", "entry", string(entry.Namespace)+"/"+entry.Key)
	c.mutex.Lock()
	c.data[entry.Namespace][entry.Key] = entry
	c.mutex.Unlock()
}

func (c *Cache) close() {
	c.closeChan <- struct{}{}
}

func (c *Cache) Close() error {
	c.close()

	return c.Dump()
}

func filterDir(files []os.DirEntry) []os.DirEntry {
	var dirs []os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file)
		}
	}
	return dirs
}

func (c *Cache) getNamespacesFromCacheFiles() ([]os.DirEntry, error) {
	namespaces, err := os.ReadDir(c.path)
	if err != nil {
		return nil, err
	}

	return filterDir(namespaces), nil
}

func (c *Cache) loadKey(namespace Namespace, key string) (Entry, error) {
	keyId := fmt.Sprintf("%s/%s", namespace, key)
	path := fmt.Sprintf("%s/%s/%s.json", c.path, namespace, key)

	c.logger.Debug("Loading key", "key", keyId)

	var entry Entry

	entry, err := c.loadFromFile(path)
	if err != nil {
		return entry, fmt.Errorf("loading key error id=%s err=%w", keyId, err)
	}

	return entry, nil
}

func (c *Cache) loadNamespace(namespace Namespace) (Data, error) {
	c.logger.Debug("Loading namespace", "namespace", namespace)

	keys, err := os.ReadDir(fmt.Sprintf("%s/%s", c.path, namespace))
	if err != nil {
		return nil, err
	}

	data := Data{}
	for _, key := range keys {
		if key.IsDir() {
			continue
		}

		keyName := strings.ReplaceAll(key.Name(), ".json", "")
		keyData, err := c.loadKey(namespace, keyName)
		if err != nil {
			return nil, err
		}

		data[keyName] = keyData
	}

	return data, nil
}

func (c *Cache) Load() error {
	c.logger.Debug("Loading cache", "path", c.path)
	namespaces, err := c.getNamespacesFromCacheFiles()
	if err != nil {
		return err
	}

	if len(namespaces) == 0 {
		c.logger.Debug("No namespaces found in cache")
		return nil
	}

	errgroup := new(errgroup.Group)
	for _, namespace := range namespaces {
		// nested function to prevent loop closure
		func(namespace os.DirEntry) {
			errgroup.Go(func() error {
				ns := Namespace(namespace.Name())
				data, err := c.loadNamespace(ns)
				if err != nil {
					return err
				}

				c.mutex.Lock()
				c.data[ns] = data
				c.mutex.Unlock()

				return nil
			})
		}(namespace)
	}

	return errgroup.Wait()
}

func (c *Cache) GetNamespace(namespace Namespace) Data {
	c.mutex.Lock()

	v, ok := c.data[namespace]
	if !ok {
		c.logger.Debug("Namespace not found. Creating a new one", "namespace", namespace)
		v = make(Data)
	}

	c.mutex.Unlock()

	return v
}

// Get returns the value of the key in the namespace
// and a boolean indicating if the key exists in the cache
func (c *Cache) Get(namespace Namespace, key string) (interface{}, bool) {
	data := c.GetNamespace(namespace)

	// Check if the key exists in the cache
	value, ok := data[key]
	if !ok {
		// If not, try to load it from the file
		v, err := c.loadKey(namespace, key)
		if err != nil {
			c.logger.Debug("Key not found in cache",
				"namespace", namespace, "key", key)
			return nil, false
		}
		value = v
	}

	c.logger.Debug("Key found in cache",
		"namespace", namespace, "key", key)

	value.AccessedTimestamp = time.Now().Unix()
	c.data[namespace][key] = value

	return value.Value, true
}

func (c *Cache) Set(namespace Namespace, key string, value interface{}) {
	c.logger.Debug("Caching",
		"namespace", namespace, "key", key)
	data := c.GetNamespace(namespace)

	c.mutex.Lock()
	ts := time.Now().Unix()
	data[key] = Entry{
		Key:               key,
		Namespace:         namespace,
		Value:             value,
		AccessedTimestamp: ts,
		CreatedTimestamp:  ts,
		UpdatedTimestamp:  ts,
		Stale:             false,
	}

	c.data[namespace] = data
	c.mutex.Unlock()
}

func (c *Cache) Dump() error {
	c.logger.Debug("Dumping cache")

	if err := c.Invalidate(); err != nil {
		return err
	}

	errgroup := new(errgroup.Group)

	for namespace, data := range c.data {
		for key := range data {
			c.mutex.Lock()
			entry := c.data[namespace][key]
			errgroup.Go(func() error {
				return c.saveEntryToFile(entry)
			})
			c.mutex.Unlock()
		}
	}

	return errgroup.Wait()
}

func (c *Cache) saveEntryToFile(entry Entry) error {
	namespace := entry.Namespace
	key := entry.Key
	path := fmt.Sprintf("%s/%s", c.path, namespace)
	filename := fmt.Sprintf("%s.json", key)

	c.logger.Debug("Writing entry", "namespaces", namespace, "key", key)

	if err := c.saveToFile(path, filename, entry); err != nil {
		return err
	}

	return nil
}

func (c *Cache) saveToFile(path string, filename string, entry Entry) error {
	if err := os.MkdirAll(path, 0777); err != nil {
		c.logger.Error(err.Error())
		panic(err)
	}

	filepath := fmt.Sprintf("%s/%s", path, filename)
	c.logger.Debug("Saving cache to file", "file", filepath)

	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) loadFromFile(filepath string) (Entry, error) {
	var data Entry

	f, err := os.Open(filepath)
	if err != nil {
		return data, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return data, err
	}

	return data, nil
}

func (c *Cache) ParseData(data interface{}, target interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(j, target); err != nil {
		return err
	}

	return nil
}

func (c *Cache) GetEntries() []Entry {
	entries := []Entry{}
	c.logger.Debug("Getting all cache entries",
		"entries", len(c.data))

	for _, data := range c.data {
		for _, entry := range data {
			entries = append(entries, entry)
		}
	}

	return entries
}

func (c *Cache) Invalidate() error {
	c.logger.Debug("Invalidating all cache entries")

	contents, err := filepath.Glob(c.path + "/*")
	if err != nil {
		return err
	}

	for _, item := range contents {
		if strings.Contains(item, ".gitkeep") {
			continue
		}

		c.logger.Debug("Removing:", "path", item)

		if err = os.RemoveAll(item); err != nil {
			return err
		}
	}

	return nil
}

func (c *Cache) Invalidate_Deprecated() error {
	c.logger.Debug("Invalidating all cache entries")

	// Clear the in-memory cache
	c.data = make(map[Namespace]Data)

	// Remove subdirectories and nested files within the cache directory
	subdirs, err := os.ReadDir(c.path)
	if err != nil {
		return err
	}

	errgroup := new(errgroup.Group)
	for _, subdir := range subdirs {
		if subdir.IsDir() {
			func(subdir os.DirEntry) {
				errgroup.Go(func() error {
					subdirPath := filepath.Join(c.path, subdir.Name())

					// Remove files within the subdirectory
					files, err := os.ReadDir(subdirPath)
					if err != nil {
						return err
					}

					for _, file := range files {
						func(file os.DirEntry) {
							errgroup.Go(func() error {
								filePath := filepath.Join(subdirPath, file.Name())
								return os.Remove(filePath)
							})
						}(file)
					}

					// Remove the subdirectory
					if err := os.Remove(subdirPath); err != nil {
						return err
					}

					return nil
				})
			}(subdir)
		}
	}

	return errgroup.Wait()
}
