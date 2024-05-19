package theme

import "github.com/charmbracelet/lipgloss"

type Style struct {
	Borders lipgloss.Style
}

var DefautlStyle = &Style{
	Borders: lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderBottom(true).
		BorderRight(true).
		BorderTop(true).
		BorderLeft(true),
}
