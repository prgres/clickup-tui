package config

type Config struct {
	Token string `fig:"token" validate:"required"`

	DefaultTeam   string `fig:"default_team" validate:"required"`
	DefaultSpace  string `fig:"default_space" validate:"required"`
	DefaultFolder string `fig:"default_folder" validate:"required"`
	DefaultList   string `fig:"default_list" validate:"required"`
}
