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

func nextView(views []clickup.View, SelectedView string) string {
	for i, view := range views {
		if view.Id == SelectedView {
			if i+1 < len(views) {
				return views[i+1].Id
			}
			return views[0].Id
		}
	}
	return views[0].Id
}

func prevView(views []clickup.View, SelectedView string) string {
	for i, view := range views {
		if view.Id == SelectedView {
			if i-1 >= 0 {
				return views[i-1].Id
			}
			return views[len(views)-1].Id
		}
	}
	return views[0].Id
}
