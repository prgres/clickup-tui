package context

import (
	"github.com/prgrs/clickup/pkg/clickup"
	"github.com/prgrs/clickup/pkg/logger1"
	"github.com/prgrs/clickup/ui/theme"
)

type UserContext struct {
	Style   theme.Style
	Logger  logger1.Logger
	Clickup *clickup.Client
}

func NewUserContext(clickup *clickup.Client, logger logger1.Logger) UserContext {
	return UserContext{
		Style:   theme.NewStyle(*theme.DefaultTheme),
		Clickup: clickup,
		Logger:  logger,
	}
}

// type Style struct {
// 	TabBorder       lipgloss.Border
// 	TabStyle        lipgloss.Style
// 	Highlight       lipgloss.AdaptiveColor
// 	ActiveTabBorder lipgloss.Border
// 	ActiveTabStyle  lipgloss.Style
// 	TabGap          lipgloss.Style
// 	DocStyle        lipgloss.Style
// }

// func NewStyle() Style {
// 	highlight := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

// 	tabBorder := lipgloss.Border{
// 		Top:         "─",
// 		Bottom:      "─",
// 		Left:        "│",
// 		Right:       "│",
// 		TopLeft:     "╭",
// 		TopRight:    "╮",
// 		BottomLeft:  "┴",
// 		BottomRight: "┴",
// 	}

// 	tabStyle := lipgloss.NewStyle().
// 		Border(tabBorder, true).
// 		BorderForeground(highlight).
// 		Padding(0, 1)

// 	activeTabBorder := lipgloss.Border{
// 		Top:         "─",
// 		Bottom:      " ",
// 		Left:        "│",
// 		Right:       "│",
// 		TopLeft:     "╭",
// 		TopRight:    "╮",
// 		BottomLeft:  "┘",
// 		BottomRight: "└",
// 	}

// 	return Style{
// 		TabBorder:       tabBorder,
// 		TabStyle:        tabStyle,
// 		Highlight:       highlight,
// 		ActiveTabBorder: activeTabBorder,
// 		ActiveTabStyle:  tabStyle.Copy().Border(activeTabBorder, true),
// 		TabGap: tabStyle.Copy().
// 			BorderTop(false).
// 			BorderLeft(false).
// 			BorderRight(false),
// 		DocStyle: lipgloss.NewStyle().Padding(1, 2, 1, 2),
// 	}
// }
