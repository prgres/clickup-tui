package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	BordersColorActive   lipgloss.Color
	BordersColorInactive lipgloss.Color
}

var DefaultTheme = &Theme{
	BordersColorActive:   lipgloss.Color("#8909FF"),
	BordersColorInactive: lipgloss.Color("#FFF"),
}
