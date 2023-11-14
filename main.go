package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/api"
	"github.com/prgrs/clickup/pkg/cache"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui"
	"github.com/prgrs/clickup/ui/context"
)

const (
	TEAM_RAMP_NETWORK   = "24301226"
	SPACE_SRE           = "48458830"
	SPACE_SRE_LIST_COOL = "q5kna-61288"

	TOKEN = "pk_42381487_1IES0AC9MGLLQND6XQ2CWIPS4KJZIR34"
)

func main() {
	logger := log.Default()

	f, err := tea.LogToFileWith("debug.log", "debug", logger)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	logger.Info("Starting up...")

	logger.Info("Initializing cache...")
	cache := cache.NewCache(logger)
	defer cache.Dump()

	if err := cache.Load(); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Initializing clickup client...")
	clickup := clickup.NewDefaultClientWithLogger(TOKEN, logger)

	logger.Info("Initializing api...")
	api := api.NewApi(clickup, logger, cache)

	logger.Info("Initializing user context...")
	ctx := context.NewUserContext(logger, &api)

	logger.Info("Initializing main model...")
	mainModel := ui.InitialModel(&ctx)

	logger.Info("Initializing program...")
	p := tea.NewProgram(mainModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Fatal(err)
	}
}
