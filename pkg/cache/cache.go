package cache

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Data map[string]interface{}

type Cache struct {
	path   string
	logger *slog.Logger
	data   map[string]Data
}

func NewCache(logger *slog.Logger) *Cache {
	path := "cache"
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
	c.logger.Info("Loading key", "key", namespace+"/"+key)

	rawData, err := c.loadFromFile(
		fmt.Sprintf("%s/%s/%s.json",
			c.path, namespace, key))
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func (c *Cache) loadNamespace(namespace string) (Data, error) {
	c.logger.Info("Loading namespace", "namespace", namespace)

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
	c.logger.Info("Loading cache from path...", "path", c.path)

	namespaces, err := c.getNamespacesFromCacheFiles()
	if err != nil {
		return err
	}

	if len(namespaces) == 0 {
		c.logger.Info("No namespaces found in cache")
		return nil
	}

	for _, namespace := range namespaces {
		data, err := c.loadNamespace(namespace.Name())
		if err != nil {
			return err
		}

		c.data[namespace.Name()] = data
	}

	return nil
}

func (c *Cache) GetNamespace(namespace string) Data {
	v, ok := c.data[namespace]
	if !ok {
		v = Data{}
	}
	return v
}

func (c *Cache) Get(namespace string, key string) (interface{}, bool) {
	data := c.GetNamespace(namespace)
	value, ok := data[key]
	if !ok {
		v, err := c.loadKey(namespace, key)
		if err != nil {
			return nil, false
		}
		value = v
	}

	return value, true
}

func (c *Cache) Set(namespace string, key string, value interface{}) {
	data := c.GetNamespace(namespace)
	data[key] = value

	path := fmt.Sprintf("%s/%s", c.path, namespace)
	filename := fmt.Sprintf("%s.json", key)

	if err := c.saveToFile(path, filename, value); err != nil {
		c.logger.Error(err.Error())
		panic(err)
	}
}

func (c *Cache) Dump() error {
	c.logger.Info("Dumping cache")
	for namespace, data := range c.data {
		for key, value := range data {
			path := fmt.Sprintf("%s/%s", c.path, namespace)
			filename := fmt.Sprintf("%s.json", key)
			if err := c.saveToFile(path, filename, value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Cache) saveToFile(path string, filename string, value interface{}) error {
	if err := os.MkdirAll(path, 0777); err != nil {
		c.logger.Error(err.Error())
		panic(err)
	}

	filepath := fmt.Sprintf("%s/%s", path, filename)
	c.logger.Info("Saving cache to file", "file", filepath)

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
