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
