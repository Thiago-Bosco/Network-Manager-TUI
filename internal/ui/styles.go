
package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Colors for the application
var (
	ColorBackground = lipgloss.Color("#0D131E")  // Dark blue background
	ColorPrimary   = lipgloss.Color("#17A649")  // Green
	ColorSecondary = lipgloss.Color("#101827")  // Detail blue
	ColorText      = lipgloss.Color("#FFFFFF")  // White
)

// Common styles
var (
	// Base styles
	BaseStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(ColorText).
		Background(ColorBackground)

	// Title styles
	TitleStyle = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true).
		Padding(0, 1)

	// Header styles
	HeaderStyle = lipgloss.NewStyle().
		Foreground(ColorText).
		Background(ColorPrimary).
		Bold(true).
		PaddingLeft(2).
		PaddingRight(2)

	// Footer styles
	FooterStyle = lipgloss.NewStyle().
		Foreground(ColorText).
		Background(ColorSecondary).
		PaddingLeft(2).
		PaddingRight(2)

	// Box styles
	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSecondary).
		Padding(1, 2)

	// Active/Selected item styles
	SelectedItemStyle = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)

	// Regular item styles
	ItemStyle = lipgloss.NewStyle().
		Foreground(ColorText)

	// Button styles
	ButtonStyle = lipgloss.NewStyle().
		Foreground(ColorText).
		Background(ColorPrimary).
		Padding(0, 3).
		Margin(0, 1).
		Bold(true)

	// Input field styles
	InputFieldStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorSecondary).
		Padding(0, 1)

	// Status styles
	StatusActiveStyle = lipgloss.NewStyle().
		Foreground(ColorPrimary).
		Bold(true)

	StatusInactiveStyle = lipgloss.NewStyle().
		Foreground(ColorText).
		Bold(true)

	// Help styles
	HelpStyle = lipgloss.NewStyle().
		Foreground(ColorSecondary)
)
