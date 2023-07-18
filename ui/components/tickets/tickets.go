package tickets

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/context"

	"github.com/charmbracelet/bubbles/table"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

const (
	SPACE_SRE_LIST_COOL = "q5kna-61288"
	SPACE_SRE           = "48458830"
)

type Model struct {
	ctx               *context.UserContext
	table             table.Model
	columns           []table.Column
	tickets           map[string][]clickup.Task
	SelectedSpace     string
	PrevSelectedSpace string
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func InitialModel(ctx *context.UserContext) Model {
	columns := []table.Column{
		{Title: "Status", Width: 15},
		{Title: "Name", Width: 90},
		{Title: "Assignees", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	return Model{
		ctx:               ctx,
		table:             t,
		columns:           columns,
		tickets:           map[string][]clickup.Task{},
		SelectedSpace:     SPACE_SRE,
		PrevSelectedSpace: SPACE_SRE,
	}
}

func (m Model) syncItems() Model {
	m.ctx.Logger.Info(fmt.Sprintf("sync items %d", len(m.tickets[m.SelectedSpace])))
	items := make([]table.Row, len(m.tickets[m.SelectedSpace]))
	for i, ticket := range m.tickets[m.SelectedSpace] {
		items[i] = table.Row{
			ticket.Status.Status,
			ticket.Name,
			ticket.GetAssignees(),
		}
	}

	m.table.SetRows(items)
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.SelectedSpace != m.PrevSelectedSpace {
		m.ctx.Logger.Info("space changed")
		m.PrevSelectedSpace = m.SelectedSpace
		msg = m.Init()
		m = m.syncItems()
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case ticketsMsg:
		m.ctx.Logger.Info("ticketsMsg")
		m.tickets[m.SelectedSpace] = msg
		m = m.syncItems()
		// _, cmd := m.list.Update(msg)
		cmds = append(cmds, cmd)

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.table.SetWidth(msg.Width - h)
		m.table.SetHeight(msg.Height - v)
	}

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}

func (m Model) Init() tea.Msg {
	if m.tickets[m.SelectedSpace] != nil {
		return ticketsMsg(m.tickets[m.SelectedSpace])
	}

	client := m.ctx.Clickup
	m.ctx.Logger.Info("fetching tickets")
	m.ctx.Logger.Info("getting views from space " + m.SelectedSpace)
	views, err := client.GetViewsFromSpace(m.SelectedSpace)
	if err != nil {
		m.ctx.Logger.Info("error")
		m.ctx.Logger.Info(err)

		return tea.Quit
	}
	// m.ctx.Logger.Info("views ", views)
	m.ctx.Logger.Info("getting tasks from view " + views[0].Id + " " + views[0].Name)
	tasks, err := client.GetTasksFromView(views[0].Id)
	if err != nil {
		m.ctx.Logger.Info("error")
		m.ctx.Logger.Info(err)

		return tea.Quit
	}

	return ticketsMsg(tasks)
}

type ticketsMsg []clickup.Task
