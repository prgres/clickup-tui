package context

import (
	"github.com/prgrs/clickup/api"
	"github.com/prgrs/clickup/internal/config"
	"github.com/prgrs/clickup/pkg/logger1"
	"github.com/prgrs/clickup/ui/theme"
)

type UserContext struct {
	Style      theme.Style
	Logger     logger1.Logger
	WindowSize WindowSize
	Api        *api.Api
	Config     *config.Config
}

type WindowSize struct {
	Width  int
	Height int
}

func NewUserContext(logger logger1.Logger, api *api.Api, config *config.Config) UserContext {
	return UserContext{
		Style:  theme.NewStyle(*theme.DefaultTheme),
		Logger: logger,
		WindowSize: WindowSize{
			Width:  0,
			Height: 0,
		},
		Api:    api,
		Config: config,
	}
}
