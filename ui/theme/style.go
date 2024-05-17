package theme

import "github.com/charmbracelet/lipgloss"

type Style struct {
	Borders lipgloss.Style
	// ListViewPort struct {
	// 	PagerStyle lipgloss.Style
	// }
	// Table struct {
	// 	CellStyle                lipgloss.Style
	// 	SelectedCellStyle        lipgloss.Style
	// 	TitleCellStyle           lipgloss.Style
	// 	SingleRuneTitleCellStyle lipgloss.Style
	// 	HeaderStyle              lipgloss.Style
	// 	RowStyle                 lipgloss.Style
	// }
	// Tabs struct {
	// 	Tab            lipgloss.Style
	// 	ActiveTab      lipgloss.Style
	// 	TabSeparator   lipgloss.Style
	// 	TabsRow        lipgloss.Style
	// 	ViewSwitcher   lipgloss.Style
	// 	ActiveView     lipgloss.Style
	// 	ViewsSeparator lipgloss.Style
	// 	InactiveView   lipgloss.Style
	// }
}

var (
	SearchHeight       = 3
	FooterHeight       = 1
	ExpandedHelpHeight = 11
	InputBoxHeight     = 8
	SingleRuneWidth    = 4
	MainContentPadding = 1
	TabsBorderHeight   = 1
	TabsContentHeight  = 2
	TabsHeight         = TabsBorderHeight + TabsContentHeight
	ViewSwitcherMargin = 1
	TableHeaderHeight  = 2
)

type CommonStyles struct {
	MainTextStyle lipgloss.Style
	FooterStyle   lipgloss.Style
	ErrorStyle    lipgloss.Style
	WaitingGlyph  string
	FailureGlyph  string
	SuccessGlyph  string
}

const (
	WaitingIcon = ""
	FailureIcon = "󰅙"
	SuccessIcon = ""
)

func BuildStyles(theme Theme) CommonStyles {
	var s CommonStyles

	// s.MainTextStyle = lipgloss.NewStyle().
	// 	Foreground(theme.PrimaryText).
	// 	Bold(true)
	// s.FooterStyle = lipgloss.NewStyle().
	// 	Background(theme.SelectedBackground).
	// 	Height(FooterHeight)
	//
	// s.ErrorStyle = s.FooterStyle.Copy().
	// 	Foreground(theme.WarningText).
	// 	MaxHeight(FooterHeight)
	//
	// s.WaitingGlyph = lipgloss.NewStyle().
	// 	Foreground(theme.FaintText).
	// 	Render(WaitingIcon)
	// s.FailureGlyph = lipgloss.NewStyle().
	// 	Foreground(theme.WarningText).
	// 	Render(FailureIcon)
	// s.SuccessGlyph = lipgloss.NewStyle().
	// 	Foreground(theme.SuccessText).
	// 	Render(SuccessIcon)
	//
	return s
}

var DefautlStyle = &Style{
	Borders: lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderBottom(true).
		BorderRight(true).
		BorderTop(true).
		BorderLeft(true),
}

func NewStyle(theme Theme) Style {
	var s Style

	// s.ListViewPort.PagerStyle = lipgloss.NewStyle().
	// 	Padding(0, 1).
	// 	Background(theme.SelectedBackground).
	// 	Foreground(theme.FaintText)
	//
	// s.Table.CellStyle = lipgloss.NewStyle().PaddingLeft(1).
	// 	PaddingRight(1).
	// 	MaxHeight(1)
	// s.Table.SelectedCellStyle = s.Table.CellStyle.Copy().
	// 	Background(theme.SelectedBackground)
	// s.Table.TitleCellStyle = s.Table.CellStyle.Copy().
	// 	Bold(true).
	// 	Foreground(theme.PrimaryText)
	// s.Table.SingleRuneTitleCellStyle = s.Table.TitleCellStyle.Copy().
	// 	Width(SingleRuneWidth)
	// s.Table.HeaderStyle = lipgloss.NewStyle().
	// 	BorderStyle(lipgloss.NormalBorder()).
	// 	BorderForeground(theme.FaintBorder).
	// 	BorderBottom(true)
	// s.Table.RowStyle = lipgloss.NewStyle().
	// 	BorderStyle(lipgloss.NormalBorder()).
	// 	BorderForeground(theme.FaintBorder)
	//
	// s.Tabs.Tab = lipgloss.NewStyle().
	// 	Faint(true).
	// 	Padding(0, 2)
	// s.Tabs.ActiveTab = s.Tabs.Tab.
	// 	Copy().
	// 	Faint(false).
	// 	Bold(true).
	// 	Background(theme.SelectedBackground).
	// 	Foreground(theme.PrimaryText)
	// s.Tabs.TabSeparator = lipgloss.NewStyle().
	// 	Foreground(theme.SecondaryBorder)
	// s.Tabs.TabsRow = lipgloss.NewStyle().
	// 	Height(TabsContentHeight).
	// 	PaddingTop(1).
	// 	PaddingBottom(0).
	// 	BorderBottom(true).
	// 	BorderStyle(lipgloss.ThickBorder()).
	// 	BorderBottomForeground(theme.PrimaryBorder)
	// s.Tabs.ViewSwitcher = lipgloss.NewStyle().
	// 	Background(theme.SecondaryText).
	// 	Foreground(theme.InvertedText).
	// 	Padding(0, 1).
	// 	Bold(true)
	//
	// s.Tabs.ActiveView = lipgloss.NewStyle().
	// 	Foreground(theme.PrimaryText).
	// 	Bold(true).
	// 	Background(theme.SelectedBackground)
	// s.Tabs.ViewsSeparator = lipgloss.NewStyle().
	// 	BorderForeground(theme.PrimaryBorder).
	// 	BorderStyle(lipgloss.NormalBorder()).
	// 	BorderRight(true)
	// s.Tabs.InactiveView = lipgloss.NewStyle().
	// 	Background(theme.FaintBorder).
	// 	Foreground(theme.SecondaryText)

	return s
}
