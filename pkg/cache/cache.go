package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

var ErrKeyNotFoundInNamespace = errors.New("key not found in namespace")

type Entry struct {
	Key       Key
	Namespace Namespace
	Value     interface{}
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

type Data map[Key]Entry

type Namespace string

type Key string

func (k Key) String() string {
	return string(k)
}

type Cache struct {
	logger *slog.Logger

	data  map[Namespace]Data
	path  string
	mutex sync.RWMutex
}

func NewCache(logger *slog.Logger, path string) *Cache {
	return &Cache{
		path:   path,
		data:   map[Namespace]Data{},
		logger: logger,
	}
}

func (c *Cache) Close() error {
	if err := c.clearCacheDir(); err != nil {
		return err
	}

	return c.Dump()
}
func (c *Cache) Load() error {
	c.logger.Debug("Loading cache from path...", "path", c.path)
	namespaces, err := c.getNamespacesFromCacheFiles()
	if err != nil {
		return err
	}

	errgroup := new(errgroup.Group)
	for _, namespace := range namespaces {
		// nested function to prevent loop closure
		func(namespace Namespace) {
			errgroup.Go(func() error {
				data, err := c.loadNamespace(namespace)
				if err != nil {
					return err
				}

				c.mutex.Lock()
				c.data[namespace] = data
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
		v = Data{}
	}

	c.mutex.Unlock()

	return v
}

func (c *Cache) Get(namespace Namespace, key Key, target interface{}) error {
	data := c.GetNamespace(namespace)

	// Check if the key exists in the cache
	value, ok := data[key]
	if !ok {
		// If not, try to load it from the file
		v, err := c.loadKey(namespace, key)
		if err != nil {
			c.logger.Debug("Key not found in cache",
				"namespace", namespace, "key", key)
			return ErrKeyNotFoundInNamespace
		}

		value = v
	}

	c.logger.Debug("Key found in cache", "namespace", namespace, "key", key)

	return c.parseData(value, target)
}

func (c *Cache) Set(namespace Namespace, key Key, value interface{}) {
	c.logger.Debug("Caching", "namespace", namespace, "key", key)
	data := c.GetNamespace(namespace)

	c.mutex.Lock()
	data[key] = Entry{
		Key:       key,
		Namespace: namespace,
		Value:     value,
	}
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
	entries := c.GetEntries()
	for _, entry := range entries {
		func(entry Entry) {
			errgroup.Go(func() error {
				return c.saveEntryToFile(entry)
			})
		}(entry)
	}

	return nil
}

func (c *Cache) GetEntries() []Entry {
	entries := []Entry{}
	c.logger.Debug("Getting all cache entries")

	for _, data := range c.data {
		for _, entry := range data {
			entries = append(entries, entry)
		}
	}

	return entries
}

func (c *Cache) Invalidate() error {
	c.logger.Debug("Invalidating all cache entries")

	// Clear the in-memory cache
	c.data = make(map[Namespace]Data)

	return c.clearCacheDir()
}

func (c *Cache) clearCacheDir() error {
	c.logger.Debug("Clearing cache dir")

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

func (c *Cache) getNamespacesFromCacheFiles() ([]Namespace, error) {
	rd, err := os.ReadDir(c.path)
	if err != nil {
		return nil, err
	}

	dirs := filterDir(rd)
	ns := make([]Namespace, len(dirs))

	for i := range dirs {
		ns[i] = Namespace(dirs[i].Name())
	}

	return ns, nil
}

func (c *Cache) loadKey(namespace Namespace, key Key) (Entry, error) {
	c.logger.Debug("Loading key", "key", fmt.Sprintf("%s/%s", namespace, key))

	data, err := c.loadFromFile(
		fmt.Sprintf("%s/%s/%s.json", c.path, namespace, key))
	if err != nil {
		return Entry{}, err
	}

	e := Entry{
		Key:       key,
		Namespace: namespace,
		Value:     data,
	}
	return e, nil
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

		keyName := Key(strings.ReplaceAll(key.Name(), ".json", ""))
		keyData, err := c.loadKey(namespace, keyName)
		if err != nil {
			return nil, err
		}

		data[keyName] = keyData
	}

	return data, nil
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

func (c *Cache) parseData(data interface{}, target interface{}) error {
	j, err := json.Marshal(data.(Entry).Value)
	if err != nil {
		return err
	}

	return json.Unmarshal(j, target)
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
