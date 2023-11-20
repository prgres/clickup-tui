package taskstabs

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/prgrs/clickup/pkg/clickup"
)

var (
	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	activeTabStyle   = lipgloss.NewStyle().Background(highlightColor)
	inactiveTabStyle = lipgloss.NewStyle().Background(lipgloss.Color("0"))
)

func removeView(views []clickup.View, s int) []clickup.View {
	return append(views[:s], views[s+1:]...)
}

func viewsToIdList(views []clickup.View) []string {
	ids := []string{}
	for _, view := range views {
		ids = append(ids, view.Id)
	}
	return ids
}

func nextTab(tabs []Tab, SelectedTab int) int {
	if SelectedTab+1 < len(tabs) {
		return SelectedTab + 1
	}
	return 0
}

func prevTab(tabs []Tab, SelectedTab int) int {
	if SelectedTab-1 >= 0 {
		return SelectedTab - 1
	}
	return len(tabs) - 1
}
