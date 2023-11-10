package viewtabs

import "github.com/charmbracelet/lipgloss"

const (
	SPACE_SRE_LIST_COOL = "q5kna-61288"
	SPACE_SRE           = "48458830"
	FOLDER_INITIATIVE   = "90050568353"
)

var (
	highlightColor   = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	activeTabStyle   = lipgloss.NewStyle().Background(highlightColor)
	inactiveTabStyle = lipgloss.NewStyle().Background(lipgloss.Color("0"))
)
