package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
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
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	logger.Info("Starting up...")

	clickup := clickup.NewDefaultClientWithLogger(TOKEN, logger)
	ctx := context.NewUserContext(clickup, logger)

	mainModel := ui.InitialModel(&ctx)

	p := tea.NewProgram(mainModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
