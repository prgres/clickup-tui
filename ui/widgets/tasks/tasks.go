package tasks

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/ui/common"
	tabletasks "github.com/prgrs/clickup/ui/components/table-tasks"
	taskssidebar "github.com/prgrs/clickup/ui/components/tasks-sidebar"
	"github.com/prgrs/clickup/ui/context"
	"golang.design/x/clipboard"
)

const WidgetId = "tasks"

type Model struct {
	log       *log.Logger
	ctx       *context.UserContext
	WidgetId  common.WidgetId
	size      common.Size
	Focused   bool
	Hidden    bool
	ifBorders bool
	keyMap    KeyMap

	state common.ComponentId

	spinner     spinner.Model
	showSpinner bool

	componenetTasksTable   tabletasks.Model
	componenetTasksSidebar taskssidebar.Model

	copyMode bool // TODO make as a widget
}

func InitialModel(ctx *context.UserContext, logger *log.Logger) Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse

	log := logger.WithPrefix(logger.GetPrefix() + "/" + WidgetId)

	size := common.Size{
		Width:  0,
		Height: 0,
	}
	var (
		componenetTasksTable   = tabletasks.InitialModel(ctx, log)
		componenetTasksSidebar = taskssidebar.InitialModel(ctx, log).WithHidden(true)
	)

	return Model{
		WidgetId:  WidgetId,
		ctx:       ctx,
		size:      size,
		Focused:   false,
		Hidden:    false,
		keyMap:    DefaultKeyMap(),
		log:       log,
		ifBorders: true,

		componenetTasksTable:   componenetTasksTable,
		componenetTasksSidebar: componenetTasksSidebar,

		state: componenetTasksTable.ComponentId,

		spinner:     s,
		showSpinner: false,

		copyMode: false,
	}
}

type KeyMap struct {
	OpenTicketInWebBrowserBatch key.Binding
	OpenTicketInWebBrowser      key.Binding
	ToggleSidebar               key.Binding
	CopyMode                    key.Binding
	CopyTaskId                  key.Binding
	CopyTaskUrl                 key.Binding
	CopyTaskUrlMd               key.Binding
	LostFocus                   key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		OpenTicketInWebBrowserBatch: key.NewBinding(
			key.WithKeys("U"),
			key.WithHelp("U", "batch open in web browser"),
		),
		OpenTicketInWebBrowser: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "open in web browser"),
		),
		ToggleSidebar: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "toggle sidebar"),
		),
		CopyMode: key.NewBinding(
			key.WithKeys("y"),
			key.WithHelp("y", "toggle copy mode"),
		),
		CopyTaskId: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "copy task id to clipboard"),
		),
		CopyTaskUrl: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "copy task url to clipboard"),
		),
		CopyTaskUrlMd: key.NewBinding(
			key.WithKeys("m"),
			key.WithHelp("m", "copy task url as markdown to clipboard"),
		),
		LostFocus: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "lost pane focus"),
		),
	}
}

func (m Model) KeyMap() help.KeyMap {
	var km help.KeyMap

	if m.copyMode {
		return common.NewKeyMap(
			func() [][]key.Binding {
				return [][]key.Binding{
					{
						m.keyMap.CopyTaskId,
						m.keyMap.CopyTaskUrl,
						m.keyMap.CopyTaskUrlMd,
						m.keyMap.LostFocus,
					},
				}
			},
			func() []key.Binding {
				return []key.Binding{
					m.keyMap.CopyTaskId,
					m.keyMap.CopyTaskUrl,
					m.keyMap.CopyTaskUrlMd,
					m.keyMap.LostFocus,
				}
			},
		)
	}

	switch m.state {
	case m.componenetTasksSidebar.ComponentId:
		km = m.componenetTasksSidebar.KeyMap()
	case m.componenetTasksTable.ComponentId:
		km = m.componenetTasksTable.KeyMap()
	}

	return common.NewKeyMap(
		func() [][]key.Binding {
			return append(
				km.FullHelp(),
				[]key.Binding{
					m.keyMap.OpenTicketInWebBrowser,
					m.keyMap.ToggleSidebar,
				},
			)
		},
		func() []key.Binding {
			return append(
				km.ShortHelp(),
				m.keyMap.OpenTicketInWebBrowser,
				m.keyMap.CopyMode,
			)
		},
	)
}

func (m *Model) SetSpinner(f bool) {
	m.showSpinner = f
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.copyMode {
			switch {
			case key.Matches(msg, m.keyMap.CopyTaskId):
				task := m.componenetTasksTable.GetHighlightedTask()
				clipboard.Write(clipboard.FmtText, []byte(task.Id))
				m.copyMode = false

			case key.Matches(msg, m.keyMap.CopyTaskUrl):
				task := m.componenetTasksTable.GetHighlightedTask()
				clipboard.Write(clipboard.FmtText, []byte(task.Url))
				m.copyMode = false

			case key.Matches(msg, m.keyMap.CopyTaskUrlMd):
				task := m.componenetTasksTable.GetHighlightedTask()
				md := fmt.Sprintf("[[#%s] - %s](%s)", task.Id, task.Name, task.Url)
				clipboard.Write(clipboard.FmtText, []byte(md))
				m.copyMode = false

			case key.Matches(msg, m.keyMap.LostFocus):
				m.copyMode = false
			}

			return m, tea.Batch(cmds...)
		}

		switch {
		case key.Matches(msg, m.keyMap.OpenTicketInWebBrowserBatch):
			tasks := m.componenetTasksTable.GetSelectedTasks()
			for _, task := range tasks {
				m.log.Debug("Opening task in the web browser", "url", task.Url)
				if err := common.OpenUrlInWebBrowser(task.Url); err != nil {
					cmds = append(cmds, common.ErrCmd(err))
					return m, tea.Batch(cmds...)
				}
			}

		case key.Matches(msg, m.keyMap.OpenTicketInWebBrowser):
			task := m.componenetTasksTable.GetHighlightedTask()
			m.log.Debug("Opening task in the web browser", "url", task.Url)
			if err := common.OpenUrlInWebBrowser(task.Url); err != nil {
				cmds = append(cmds, common.ErrCmd(err))
				return m, tea.Batch(cmds...)
			}

		case key.Matches(msg, m.keyMap.ToggleSidebar):
			m.log.Debug("Toggle sidebar")
			m.componenetTasksSidebar.SetHidden(!m.componenetTasksSidebar.GetHidden())

		case key.Matches(msg, m.keyMap.CopyMode):
			m.log.Debug("Toggle copy mode")
			m.copyMode = true

		case key.Matches(msg, m.keyMap.LostFocus):
			switch m.state {
			case m.componenetTasksSidebar.ComponentId:
				m.state = m.componenetTasksTable.ComponentId
				m.componenetTasksSidebar.SetFocused(false)
				m.componenetTasksTable.SetFocused(true)

			case m.componenetTasksTable.ComponentId:
				m.componenetTasksSidebar.SetFocused(false)
				m.componenetTasksTable.SetFocused(false)
			}

			m.componenetTasksSidebar, cmd = m.componenetTasksSidebar.Update(msg)
			cmds = append(cmds, cmd)
			m.componenetTasksTable, cmd = m.componenetTasksTable.Update(msg)
			cmds = append(cmds, cmd)

			cmds = append(cmds, LostFocusCmd())
			return m, tea.Batch(cmds...)
		}

		switch m.state {
		case m.componenetTasksSidebar.ComponentId:
			m.componenetTasksSidebar, cmd = m.componenetTasksSidebar.Update(msg)
		case m.componenetTasksTable.ComponentId:
			m.componenetTasksTable, cmd = m.componenetTasksTable.Update(msg)
		}

		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	case tabletasks.TaskSelectedMsg:
		id := string(msg)
		m.log.Infof("Received: taskstable.TaskSelectedMsg: %s", id)

		m.state = m.componenetTasksSidebar.ComponentId

		m.componenetTasksSidebar.
			SetFocused(true).
			SetHidden(false)

		m.componenetTasksTable.SetFocused(false)

		if err := m.componenetTasksSidebar.SetTask(id); err != nil {
			cmds = append(cmds, common.ErrCmd(err))
		}

		cmds = append(cmds, cmd)
	}

	m.componenetTasksTable, cmd = m.componenetTasksTable.Update(msg)
	cmds = append(cmds, cmd)

	m.componenetTasksSidebar, cmd = m.componenetTasksSidebar.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) SetTasks(tasks []clickup.Task) error {
	m.showSpinner = false
	m.componenetTasksTable.SetTasks(tasks)

	if len(tasks) == 0 {
		m.componenetTasksSidebar.SetHidden(true)
		return nil
	}

	// TODO: check if it should yield at all or move it to cmd
	id := tasks[0].Id

	return m.componenetTasksSidebar.SetTask(id)
}

func (m Model) View() string {
	bColor := m.ctx.Theme.BordersColorInactive
	if m.Focused {
		bColor = m.ctx.Theme.BordersColorActive
	}
	if m.copyMode {
		bColor = m.ctx.Theme.BordersColorCopyMode
	}

	style := lipgloss.NewStyle().
		Width(m.size.Width).
		MaxWidth(m.size.Width).
		Height(m.size.Height).
		MaxHeight(m.size.Height)

	styleBorders := m.ctx.Style.Borders.Copy().
		BorderForeground(bColor)

	borderMargin := 0
	if m.ifBorders {
		borderMargin = 2
	}

	size := common.Size{
		Width:  m.size.Width - borderMargin,
		Height: m.size.Height - borderMargin,
	}

	if m.showSpinner {
		return style.Copy().
			Inherit(styleBorders).
			Width(m.size.Width - borderMargin).
			MaxWidth(m.size.Width + borderMargin).
			Height(m.size.Height - borderMargin).
			MaxHeight(m.size.Height + borderMargin).
			Render(
				lipgloss.Place(
					size.Width, size.Height,
					lipgloss.Center,
					lipgloss.Center,
					fmt.Sprintf("%s Loading lists...", m.spinner.View()),
				),
			)
	}

	if m.componenetTasksTable.TotalRows() == 0 {
		return style.Copy().
			Inherit(styleBorders).
			Width(m.size.Width - borderMargin).
			MaxWidth(m.size.Width + borderMargin).
			Height(m.size.Height - borderMargin).
			MaxHeight(m.size.Height + borderMargin).
			Render(
				lipgloss.Place(
					size.Width, size.Height,
					lipgloss.Center,
					lipgloss.Center,
					"No tasks found",
				),
			)
	}

	tasksTableBorders := m.ctx.Theme.BordersColorInactive
	if m.componenetTasksTable.GetFocused() {
		tasksTableBorders = m.ctx.Theme.BordersColorActive
	}
	if m.copyMode {
		tasksTableBorders = m.ctx.Theme.BordersColorCopyMode
	}

	tmpStyle := style.Copy().
		Inherit(styleBorders)

	var (
		contentTasksTable   string
		contentTasksSidebar string
	)

	if m.componenetTasksSidebar.Hidden {
		m.componenetTasksTable.SetSize(size)

		tmpStyle = tmpStyle.Copy().
			BorderForeground(tasksTableBorders).
			Width(size.Width).
			MaxWidth(size.Width + borderMargin).
			Height(size.Height).
			MaxHeight(m.size.Height + borderMargin)

		m.componenetTasksTable.SetSize(size)
		contentTasksTable = tmpStyle.Render(m.componenetTasksTable.View())
		contentTasksSidebar = ""
	} else {
		// TODO: WTF?!
		size.Width /= 2

		size.Height += borderMargin
		m.componenetTasksSidebar.SetSize(size)
		size.Height -= borderMargin

		size.Width -= borderMargin // size.Width -= 2 * borderMargin

		tmpStyle = tmpStyle.Copy().
			Width(size.Width).
			BorderForeground(tasksTableBorders).
			MaxWidth(size.Width + borderMargin).
			Height(size.Height).
			MaxHeight(m.size.Height + borderMargin)

		m.componenetTasksTable.SetSize(size)

		contentTasksTable = tmpStyle.Render(m.componenetTasksTable.View())
		contentTasksSidebar = m.componenetTasksSidebar.View()
	}

	return style.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			contentTasksTable,
			contentTasksSidebar,
		))
}

func (m Model) SetFocused(f bool) Model {
	m.Focused = f

	switch m.state {
	case m.componenetTasksSidebar.ComponentId:
		m.componenetTasksSidebar.SetFocused(f)
	case m.componenetTasksTable.ComponentId:
		m.componenetTasksTable.SetFocused(f)
	}

	return m
}

func (m *Model) SetSize(s common.Size) {
	m.size = s
}

func (m *Model) Init() error {
	m.log.Info("Initializing...")
	return nil
}
