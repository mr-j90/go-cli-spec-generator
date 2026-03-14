// Package tui implements the interactive terminal UI using Bubble Tea.
package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle        lipgloss.Style
	selectedStyle     lipgloss.Style
	normalStyle       lipgloss.Style
	descStyle         lipgloss.Style
	checkboxChecked   lipgloss.Style
	checkboxUnchecked lipgloss.Style
	cursorStyle       lipgloss.Style
	helpStyle         lipgloss.Style
)

// initStyles initializes all Lipgloss styles. Pass noColor=true to strip ANSI colors.
func initStyles(noColor bool) {
	if noColor {
		titleStyle        = lipgloss.NewStyle().Bold(true)
		selectedStyle     = lipgloss.NewStyle().Bold(true)
		normalStyle       = lipgloss.NewStyle()
		descStyle         = lipgloss.NewStyle()
		checkboxChecked   = lipgloss.NewStyle()
		checkboxUnchecked = lipgloss.NewStyle()
		cursorStyle       = lipgloss.NewStyle()
		helpStyle         = lipgloss.NewStyle()
	} else {
		titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("62"))
		selectedStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
		normalStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
		descStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		checkboxChecked   = lipgloss.NewStyle().Foreground(lipgloss.Color("78"))
		checkboxUnchecked = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		cursorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
		helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	}
}
