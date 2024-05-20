package theme

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	ColorWhite = lipgloss.Color("#FFFFFF")
	ColorBlack = lipgloss.Color("#000000")
)

type Theme struct {
	BordersColorActive   lipgloss.Color
	BordersColorInactive lipgloss.Color
	BordersColorCopyMode lipgloss.Color
}

var DefaultTheme = &Theme{
	BordersColorActive:   lipgloss.Color("#8909FF"),
	BordersColorInactive: lipgloss.Color("#FFF"),
	BordersColorCopyMode: lipgloss.Color("#e6cc00"),
}
