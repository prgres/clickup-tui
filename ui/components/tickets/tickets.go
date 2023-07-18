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
)

type Model struct {
	ctx     *context.UserContext
	table   table.Model
	columns []table.Column
	tickets []clickup.Task
}

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

func InitialModel(ctx *context.UserContext) Model {
	columns := []table.Column{
		{Title: "Status", Width: 15},
		{Title: "Name", Width: 90},
		// {Title: "Assignees", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
	)

	return Model{
		ctx:     ctx,
		table:   t,
		columns: columns,
	}
}

func (m Model) syncItems() Model {
	m.ctx.Logger.Info(fmt.Sprintf("sync items %d", len(m.tickets)))
	items := make([]table.Row, len(m.tickets))
	for i := range m.tickets {
		items[i] = table.Row{
			m.tickets[i].Status.Status,
			m.tickets[i].Name,
			// strings.Join(m.tickets[i].Assignees, ","),
		}
	}

	m.table.SetRows(items)
	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case ticketsMsg:
		m.ctx.Logger.Info("ticketsMsg")
		m.tickets = msg
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
	client := m.ctx.Clickup
	m.ctx.Logger.Info("fetching tickets")

	tasks, err := client.GetTasksFromView(SPACE_SRE_LIST_COOL)
	if err != nil {
		return tea.Quit
	}

	return ticketsMsg(tasks)
}

type ticketsMsg []clickup.Task
