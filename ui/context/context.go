package context

import (
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/api"
	"github.com/prgrs/clickup/internal/config"
	"github.com/prgrs/clickup/ui/theme"
)

type UserContext struct {
	Api        *api.Api
	Config     *config.Config
	Style      theme.Style
	WindowSize WindowSize
}

type WindowSize struct {
	Width      int
	Height     int
	MetaHeight int
}

func NewUserContext(logger *log.Logger, api *api.Api, config *config.Config) UserContext {
	return UserContext{
		Style: theme.NewStyle(*theme.DefaultTheme),
		WindowSize: WindowSize{
			Width:      0,
			Height:     0,
			MetaHeight: 0,
		},
		Api:    api,
		Config: config,
	}
}
