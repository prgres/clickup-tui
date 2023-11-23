package main

import (
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/kkyr/fig"
	"github.com/prgrs/clickup/api"
	"github.com/prgrs/clickup/internal/config"
	"github.com/prgrs/clickup/pkg/cache"
	"github.com/prgrs/clickup/ui"
	"github.com/prgrs/clickup/ui/context"
)

func main() {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		// ReportCaller:    true,
		ReportTimestamp: true,
	})

	f, err := tea.LogToFileWith("debug.log", logger.GetPrefix(), logger)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	logger.Info("Starting up...")

	logger.Info("Initializing config...")
	var cfg config.Config
	if err := fig.Load(&cfg,
		fig.File("config.yaml"),
		fig.Dirs(
			".",
			"/etc/myapp",
			"/home/user/myapp",
			"$HOME/.config/clickup-tui",
		),
	); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Initializing cache...")
	cache := cache.NewCache(slog.New(logger.WithPrefix("Cache")))
	defer func() {
		_ = cache.Dump()
	}()

	if err := cache.Load(); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Initializing api...")
	api := api.NewApi(logger, cache, cfg.Token)

	logger.Info("Initializing user context...")
	ctx := context.NewUserContext(logger, &api, &cfg)

	logger.Info("Initializing main model...")
	mainModel := ui.InitialModel(&ctx, logger)

	logger.Info("Initializing program...")
	p := tea.NewProgram(mainModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Fatal(err)
	}
}
