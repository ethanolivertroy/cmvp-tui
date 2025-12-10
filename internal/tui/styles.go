package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ethanolivertroy/nist-cmvp-cli/internal/model"
)

var (
	// Color palette
	PrimaryColor   = lipgloss.Color("#7D56F4") // Purple
	SecondaryColor = lipgloss.Color("#04B575") // Green
	WarningColor   = lipgloss.Color("#FFCC00") // Yellow
	ErrorColor     = lipgloss.Color("#FF5F56") // Red
	SubtleColor    = lipgloss.Color("#626262") // Gray

	// Status colors
	ActiveColor     = lipgloss.Color("#04B575") // Green
	HistoricalColor = lipgloss.Color("#626262") // Gray
	InProcessColor  = lipgloss.Color("#FFCC00") // Yellow

	// App styles
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(PrimaryColor).
			Padding(0, 1).
			Bold(true)

	StatusBarStyle = lipgloss.NewStyle().
			Foreground(SubtleColor).
			Padding(0, 1)

	// Detail view styles
	DetailTitleStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor).
				Bold(true).
				MarginBottom(1)

	DetailLabelStyle = lipgloss.NewStyle().
				Foreground(SubtleColor).
				Width(18)

	DetailValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA"))

	DetailURLStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00AAFF")).
			Underline(true)

	// Status badge styles
	ActiveBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ActiveColor).
			Padding(0, 1)

	HistoricalBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(HistoricalColor).
			Padding(0, 1)

	InProcessBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(InProcessColor).
			Padding(0, 1)

	// Help style
	HelpStyle = lipgloss.NewStyle().
			Foreground(SubtleColor).
			MarginTop(1)

	// Caveat warning style (important security warnings)
	CaveatStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF6B6B")).
			Padding(0, 1).
			Bold(true)

	// Level badge styles (color coded by security level)
	Level1Badge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#04B575")).
			Padding(0, 1)

	Level2Badge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#000000")).
			Background(lipgloss.Color("#FFCC00")).
			Padding(0, 1)

	Level3Badge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF9500")).
			Padding(0, 1)

	Level4Badge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF5F56")).
			Padding(0, 1)

	// Algorithm tag style
	AlgorithmStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#5B5FC7")).
			Padding(0, 1).
			MarginRight(1)

	// Description style (for longer text)
	DescriptionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#CCCCCC")).
				Width(60)
)

// StatusBadge returns a styled status badge for the given status
func StatusBadge(status model.ModuleStatus) string {
	switch status {
	case model.StatusActive:
		return ActiveBadge.Render("ACTIVE")
	case model.StatusHistorical:
		return HistoricalBadge.Render("HISTORICAL")
	case model.StatusInProcess:
		return InProcessBadge.Render("IN PROCESS")
	default:
		return ""
	}
}

// LevelBadge returns a color-coded security level badge
func LevelBadge(level int) string {
	if level == 0 {
		return ""
	}
	switch level {
	case 1:
		return Level1Badge.Render("Level 1")
	case 2:
		return Level2Badge.Render("Level 2")
	case 3:
		return Level3Badge.Render("Level 3")
	case 4:
		return Level4Badge.Render("Level 4")
	default:
		return ""
	}
}
