package ui

import "github.com/charmbracelet/lipgloss"

var (
	accentBlue   = lipgloss.Color("#00D9FF")
	accentPurple = lipgloss.Color("#B57EDC")
	accentGreen  = lipgloss.Color("#50FA7B")
	textPrimary  = lipgloss.Color("#F8F8F2")
	textMuted    = lipgloss.Color("#6272A4")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentBlue).
			MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(accentPurple).
			Bold(true).
			PaddingLeft(2)

	normalStyle = lipgloss.NewStyle().
			Foreground(textMuted).
			PaddingLeft(2)

	helpStyle = lipgloss.NewStyle().
			Foreground(textMuted).
			MarginTop(1)

	statusStyle = lipgloss.NewStyle().
			Foreground(accentGreen).
			MarginBottom(1)

	transcriptBoxStyle = lipgloss.NewStyle().
				Foreground(textPrimary).
				BorderForeground(accentBlue).
				Padding(1)
)
