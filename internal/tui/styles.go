package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("#6B50FF") // Atlas Purple
	ColorSecondary = lipgloss.Color("#FF5F87") // Atlas Pink
	ColorCyan      = lipgloss.Color("#00D7FF") // Atlas Cyan
	ColorText      = lipgloss.Color("#FFFDF5")
	ColorSubtext   = lipgloss.Color("#626262")
	ColorError     = lipgloss.Color("#FF5F87")
	ColorSuccess   = lipgloss.Color("#00D787")

	// Base Styles
	StyleBase = lipgloss.NewStyle().Foreground(ColorText)

	StyleSubtext = lipgloss.NewStyle().Foreground(ColorSubtext)

	// Auth Styles
	StyleAuthBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2).
			Align(lipgloss.Center)

	StyleAuthHeader = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			MarginBottom(1)

	// List Styles
	StyleListItem = lipgloss.NewStyle().
			PaddingLeft(2)

	StyleListItemSelected = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(ColorSecondary).
				SetString("❯ ")

	StyleListHeader = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Padding(0, 1).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(ColorSubtext)

	// Editor Styles
	StyleEditorLabel = lipgloss.NewStyle().
			Foreground(ColorSubtext).
			Width(10)

	StyleEditorActive = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorPrimary)

	StyleStatusBar = lipgloss.NewStyle().
			Foreground(ColorSubtext).
			Padding(0, 1)
)