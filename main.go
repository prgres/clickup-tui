package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/kkyr/fig"
	"github.com/prgrs/clickup/api"
	"github.com/prgrs/clickup/internal/config"
	"github.com/prgrs/clickup/pkg/cache"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui"
	"github.com/prgrs/clickup/ui/context"
)

func main() {
	logger := log.Default()

	f, err := tea.LogToFileWith("debug.log", "debug", logger)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	logger.Info("Starting up...")

	logger.Info("Initializing config...")
	var cfg config.Config
	fig.Load(&cfg,
		fig.File("config.yaml"),
		fig.Dirs(
			".",
			"/etc/myapp",
			"/home/user/myapp",
			"$HOME/.config/clickup-tui",
		),
	)

	logger.Info("Initializing cache...")
	cache := cache.NewCache(logger)
	defer cache.Dump()

	if err := cache.Load(); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Initializing clickup client...")
	clickup := clickup.NewDefaultClientWithLogger(cfg.Token, logger)

	logger.Info("Initializing api...")
	api := api.NewApi(clickup, logger, cache)

	logger.Info("Initializing user context...")
	ctx := context.NewUserContext(logger, &api, &cfg)

	logger.Info("Initializing main model...")
	mainModel := ui.InitialModel(&ctx)

	logger.Info("Initializing program...")
	p := tea.NewProgram(mainModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Fatal(err)
	}
}
