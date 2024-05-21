package cache

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Entry struct {
	// The key of the Entry
	Key string

	// The namespace of the Entryl
	Namespace string
}

type Data map[string]interface{}

type Cache struct {
	logger *slog.Logger

	data  map[string]Data
	path  string
	mutex sync.RWMutex
}

func NewCache(logger *slog.Logger, path string) *Cache {
	return &Cache{
		path:   path,
		data:   map[string]Data{},
		logger: logger,
	}
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

func (c *Cache) loadKey(namespace string, key string) (interface{}, error) {
	c.logger.Debug("Loading key", "key", namespace+"/"+key)

	rawData, err := c.loadFromFile(
		fmt.Sprintf("%s/%s/%s.json",
			c.path, namespace, key))
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func (c *Cache) loadNamespace(namespace string) (Data, error) {
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
	c.logger.Debug("Loading cache from path...", "path", c.path)
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
				data, err := c.loadNamespace(namespace.Name())
				if err != nil {
					return err
				}

				c.mutex.Lock()
				c.data[namespace.Name()] = data
				c.mutex.Unlock()

				return nil
			})
		}(namespace)
	}

	return errgroup.Wait()
}

func (c *Cache) GetNamespace(namespace string) Data {
	c.mutex.Lock()

	v, ok := c.data[namespace]
	if !ok {
		v = Data{}
	}

	c.mutex.Unlock()

	return v
}

// Get returns the value of the key in the namespace
// and a boolean indicating if the key exists in the cache
func (c *Cache) Get(namespace string, key string) (interface{}, bool) {
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

	return value, true
}

func (c *Cache) Set(namespace string, key string, value interface{}) {
	c.logger.Debug("Caching",
		"namespace", namespace, "key", key)
	data := c.GetNamespace(namespace)

	c.mutex.Lock()
	data[key] = value
	c.mutex.Unlock()

	path := fmt.Sprintf("%s/%s", c.path, namespace)
	filename := fmt.Sprintf("%s.json", key)

	if err := c.saveToFile(path, filename, value); err != nil {
		c.logger.Error(err.Error())
		panic(err)
	}
}

func (c *Cache) Dump() error {
	c.logger.Debug("Dumping cache")

	errgroup := new(errgroup.Group)
	for namespace, data := range c.data {
		for key, value := range data {
			func(namespace string, key string, value interface{}) {
				errgroup.Go(func() error {
					path := fmt.Sprintf("%s/%s", c.path, namespace)
					filename := fmt.Sprintf("%s.json", key)

					return c.saveToFile(path, filename, value)
				})
			}(namespace, key, value)
		}
	}

	return errgroup.Wait()
}

func (c *Cache) saveToFile(path string, filename string, value interface{}) error {
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

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) loadFromFile(filepath string) (interface{}, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data interface{}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
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

	j, _ := json.Marshal(c.data)
	c.logger.Debug(string(j))

	for ns, data := range c.data {
		for key := range data {
			entries = append(entries, Entry{
				Namespace: ns,
				Key:       key,
			})
		}
	}

	return entries
}

func (c *Cache) Invalidate() error {
	c.logger.Debug("Invalidating all cache entries")

	// Clear the in-memory cache
	c.data = make(map[string]Data)

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
	c.data = make(map[string]Data)

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
