package viewtabs

import "github.com/charmbracelet/lipgloss"

var (
	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	activeTabStyle   = lipgloss.NewStyle().Background(highlightColor)
	inactiveTabStyle = lipgloss.NewStyle().Background(lipgloss.Color("0"))
)
