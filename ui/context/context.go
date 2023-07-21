package context

import (
	"github.com/prgrs/clickup/pkg/cache"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/pkg/logger1"
	"github.com/prgrs/clickup/ui/theme"
)

type UserContext struct {
	Style   theme.Style
	Logger  logger1.Logger
	Clickup *clickup.Client

	WindowSize WindowSize

	Cache *cache.Cache
}

type WindowSize struct {
	Width  int
	Height int
}

func NewUserContext(clickup *clickup.Client, logger logger1.Logger, cache *cache.Cache) UserContext {
	return UserContext{
		Style:   theme.NewStyle(*theme.DefaultTheme),
		Logger:  logger,
		Clickup: clickup,

		WindowSize: WindowSize{
			Width:  0,
			Height: 0,
		},

		Cache: cache,
	}
}
