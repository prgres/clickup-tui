package config

type Config struct {
	Token string `fig:"token"` // required

	DefaultWorkspace string `fig:"default_workspace"`
	DefaultSpace     string `fig:"default_space"`
	DefaultFolder    string `fig:"default_folder"`
	DefaultList      string `fig:"default_list"`
}
