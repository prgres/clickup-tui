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

	width = 96
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

	// client := NewClickUp(TOKEN)

	// teams, err := client.GetTeams()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Info(client.ToJson(teams))

	// spaces, err := client.GetSpaces(TEAM_RAMP_NETWORK)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Info(client.ToJson(spaces))

	// folders, err := client.GetFolders(SPACE_SRE)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Info(client.ToJson(folders))

	// viewsR, err := client.GetSpaceViews(SPACE_SRE)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// views := viewsR
	// // views := viewsR[0:2]
	// if err := os.WriteFile("views.json", client.ToJsonByte(views), 0644); err != nil {
	// 	log.Fatal(err)
	// }

	// for _, view := range views {
	// 	if strings.Contains(view.Name, "cool") {
	// log.Info(client.ToJson(view))
	// 	}
	// }

	// tasks, err := client.GetViewTasks(SPACE_SRE_LIST_COOL)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Info(client.ToJson(tasks))

	// if err := os.WriteFile("views-task.json", client.ToJsonByte(viewsTask), 0644); err != nil {
	// 	log.Fatal(err)
	// }

	// viewsByte, err := os.ReadFile("views.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var views []clickup.View
	// if err := json.Unmarshal(viewsByte, &views); err != nil {
	// 	log.Fatal(err)
	// }

	// viewsNames := make([]string, len(views))
	// for i, view := range views {
	// 	viewsNames[i] = view.Name
	// }

	// viewsTask := make([]string, len(views))
	// for i := range viewsTask {
	// 	doc := strings.Builder{}
	// 	// 	tasks, err := client.GetViewTasks(view.Id)
	// 	// 	if err != nil {
	// 	// 		log.Fatal(err)
	// 	// 	}

	// 	// 	// log.Info("------")
	// 	// 	// log.Info(client.ToJson(tasks))
	// 	for _, task := range []string{"a", "b", "c"} {
	// 		doc.WriteString(" *  " + task)
	// 		doc.WriteString("\n")
	// 	}
	// 	// 	// content := client.ToJson(tasks)
	// 	// 	// if len(content) > 10 {
	// 	// 	// 	content = content[0:10]
	// 	// 	// }
	// 	viewsTask[i] = doc.String()
	// 	// 	viewsTask[i] = client.ToJson(tasks)
	// }

	// p := tea.NewProgram(tabs.InitialModel(viewsNames, viewsTask), tea.WithAltScreen())
	// if _, err := p.Run(); err != nil {
	// 	log.Fatal(err)
	// }
}
