package tabs

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/prgrs/clickup/ui/context"
)

var (
	width = 96
)

type Model struct {
	ctx         *context.UserContext
	Tabs        []string
	TabsContent []string
	activeTab   int
}

func InitialModel(ctx *context.UserContext) Model {
	return Model{
		ctx:         ctx,
		Tabs:        []string{},
		TabsContent: []string{},
		activeTab:   0,
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l", "n", "tab":
			m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
			return m, nil
		case "left", "h", "p", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, nil
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m Model) View() string {
	return ""
	// doc := strings.Builder{}

	// var renderedTabs []string

	// for i, t := range m.Tabs {
	// 	style := m.ctx.Style.TabStyle.Copy()
	// 	// var style lipgloss.Style
	// 	_, _, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
	// 	// isFirst, isLast, isActive := i == 0, i == len(m.Tabs)-1, i == m.activeTab
	// 	if isActive {
	// 		style = m.ctx.Style.ActiveTabStyle.Copy()
	// 	}
	// 	renderedTabs = append(renderedTabs, style.Render(t))
	// }
	// row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// gap := m.ctx.Style.TabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
	// row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
	// doc.WriteString(row)
	// doc.WriteString("\n")
	// // doc.WriteString(windowStyle.Width((lipgloss.Width(row) - windowStyle.GetHorizontalFrameSize())).Render(m.TabsContent[m.activeTab]))
	// doc.WriteString(m.TabsContent[m.activeTab])
	// // doc.Max

	// physicalWidth, physicalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	// doc.WriteString(fmt.Sprintf("%d %d", physicalWidth, physicalHeight))

	// m.ctx.Style.DocStyle = m.ctx.Style.DocStyle.MaxHeight(physicalHeight - 1)
	// m.ctx.Style.DocStyle = m.ctx.Style.DocStyle.MaxWidth(physicalWidth - 1)
	// return m.ctx.Style.DocStyle.Render(doc.String())

}

func (m Model) Init() tea.Cmd {
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
