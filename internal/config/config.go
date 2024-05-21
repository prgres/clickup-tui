package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var (
	ErrMissingToken = fmt.Errorf("ClickUp API token is required")
	ErrFileNotFound = fmt.Errorf("file not found")
)

const (
	DefaultPathPrefix = ".config/clickup-tui"
	DefaultFilename   = "config.yaml"

	HowToGetToken = "Follow the steps: [ClickUp API docs: Generate your personal API token](https://clickup.com/api/developer-portal/authentication/#generate-your-personal-api-token) and please set it in the config file. See https://docs.clickup.com/en/articles/1367130-getting-started-with-the-clickup-api"
)

type Config struct {
	Token            string `yaml:"token"` // required
	DefaultWorkspace string `yaml:"default_workspace"`
	DefaultSpace     string `yaml:"default_space"`
	DefaultFolder    string `yaml:"default_folder"`
	DefaultList      string `yaml:"default_list"`
	Path             string `yaml:"-"`
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FindConfigFile(filename string, paths []string) (string, error) {
	for _, dir := range paths {
		path := filepath.Join(dir, filename)
		if fileExists(path) {
			return path, nil
		}
	}
	return "", ErrFileNotFound
}

func CreateEmptyCfgFile(filename string, path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	data, err := yaml.Marshal(Config{})
	if err != nil {
		return err
	}

	configPath := filepath.Join(path, filename)
	err = os.WriteFile(configPath, data, 0755)
	if err != nil {
		return fmt.Errorf("unable to write file: %w", err)
	}

	return nil
}

func (c *Config) Save() error {
	data, err := yaml.Marshal(*c)
	if err != nil {
		return err
	}

	configPath := filepath.Join(c.Path)
	err = os.WriteFile(configPath, data, 0755)
	if err != nil {
		return fmt.Errorf("unable to write file: %w", err)
	}

	return nil
}

func Init(path string) (*Config, error) {
	var cfg Config

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	byteValue, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(byteValue, &cfg); err != nil {
		return nil, err
	}
	cfg.Path = path

	if cfg.Token == "" {
		return &cfg, ErrMissingToken
	}

	return &cfg, nil
}
