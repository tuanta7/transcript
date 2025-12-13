package ui

import "github.com/charmbracelet/lipgloss"

var (
	accentBlue   = lipgloss.Color("#A7E9FF") // pastel cyan
	accentPurple = lipgloss.Color("#D8B7FF") // lavender
	accentGreen  = lipgloss.Color("#BFFFC5") // mint
	textPrimary  = lipgloss.Color("#F7F7FF") // very light off-white
	textMuted    = lipgloss.Color("#9AA2B2") // muted gray-blue

	logoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentBlue).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(textMuted).
			Italic(true).
			MarginBottom(1)

	menuBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentBlue).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(accentPurple).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(textMuted)

	cursorStyle = lipgloss.NewStyle().
			Foreground(accentGreen).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(textMuted).
			MarginTop(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(accentPurple).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(accentGreen).
			Bold(true).
			MarginBottom(1)

	recordingDotStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B")).
				Bold(true)

	transcriptBoxStyle = lipgloss.NewStyle().
				Foreground(textPrimary).
				BorderForeground(accentBlue).
				Border(lipgloss.RoundedBorder()).
				Padding(0, 1)

	transcriptHeaderStyle = lipgloss.NewStyle().
				Foreground(accentBlue).
				Bold(true).
				MarginBottom(1)

	transcriptTextStyle = lipgloss.NewStyle().
				Foreground(textPrimary)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B6B"))

	durationValueStyle = lipgloss.NewStyle().
				Foreground(accentGreen).
				Bold(true)
)
