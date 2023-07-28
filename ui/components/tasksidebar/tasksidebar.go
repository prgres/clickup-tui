package tasksidebar

import (
	"encoding/json"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	"github.com/prgrs/clickup/ui/context"
)

type TaskSelectedMsg string

func TaskSelectedCmd(task string) tea.Cmd {
	return func() tea.Msg {
		return TaskSelectedMsg(task)
	}
}

type Model struct {
	ctx      *context.UserContext
	viewport viewport.Model
}

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		ctx:      ctx,
		viewport: viewport.New(30, 30),
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case InitMsg:
		m.ctx.Logger.Info("TaskSidebar receive InitMsg")
		m.viewport.SetContent(content)

	case common.WindowSizeMsg:
		m.ctx.Logger.Info("TaskSidebar receive tea.WindowSizeMsg")

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 2

	case TaskSelectedMsg:
		m.ctx.Logger.Info("TaskSidebar receive TaskSelectedMsg")
		id := strings.Split(string(msg), " ")[0]
		task, err := m.getTask(id)
		if err != nil {
			m.ctx.Logger.Fatal(err)
			return m, nil
		}

		taskJson, err := json.MarshalIndent(task, "", "  ")
		if err != nil {
			m.ctx.Logger.Fatal(err)
			return m, nil
		}

		m.viewport.SetContent(string(taskJson))
		m.viewport.GotoTop()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.viewport.View()
}

func (m Model) Init() tea.Cmd {
	m.ctx.Logger.Info("Initializing component: TaskSidebar")
	// m.viewport.SetContent("Hello, world!")
	// m.viewport.YPosition = 10

	return InitCmd()
}

func (m Model) getTask(id string) (clickup.Task, error) {
	m.ctx.Logger.Infof("Getting task: %s", id)

	// data, ok := m.ctx.Cache.Get("task", id)
	// if ok {
	// 	m.ctx.Logger.Infof("Task found in cache")
	// 	var task clickup.Task
	// 	if err := m.ctx.Cache.ParseData(data, &task); err != nil {
	// 		return clickup.Task{}, err
	// 	}

	// 	return task, nil
	// }
	m.ctx.Logger.Info("Task not found in cache")

	m.ctx.Logger.Info("Fetching task from API")
	client := m.ctx.Clickup

	task, err := client.GetTask(id)
	if err != nil {
		return clickup.Task{}, err
	}
	m.ctx.Logger.Infof("Found tasks %s", id)

	m.ctx.Logger.Info("Caching tasks")
	m.ctx.Cache.Set("task", id, task)

	return task, nil
}

var content = `Glow
====

A casual introduction. 你好世界!

## Let’s talk about artichokes

The _artichoke_ is mentioned as a garden plant in the 8th century BC by Homer
**and** Hesiod. The naturally occurring variant of the artichoke, the cardoon,
which is native to the Mediterranean area, also has records of use as a food
among the ancient Greeks and Romans. Pliny the Elder mentioned growing of
_carduus_ in Carthage and Cordoba.

> He holds him with a skinny hand,
> ‘There was a ship,’ quoth he.
> ‘Hold off! unhand me, grey-beard loon!’
> An artichoke, dropt he.

--Samuel Taylor Coleridge, [The Rime of the Ancient Mariner][rime]

[rime]: https://poetryfoundation.org/poems/43997/

## Other foods worth mentioning

1. Carrots
1. Celery
1. Tacos
    * Soft
    * Hard
1. Cucumber

## Things to eat today

* [x] Carrots
* [x] Ramen
* [ ] Currywurst

### Power levels of the aforementioned foods

| Name       | Power | Comment          |
| ---        | ---   | ---              |
| Carrots    | 9001  | It’s over 9000?! |
| Ramen      | 9002  | Also over 9000?! |
| Currywurst | 10000 | What?!           |

## Currying Artichokes

_Alcachofa_, if you were wondering, is artichoke in Spanish.`
