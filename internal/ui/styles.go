package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Colors for the application
var (
	ColorPrimary   = lipgloss.Color("#5E81AC")
	ColorSecondary = lipgloss.Color("#81A1C1")
	ColorAccent    = lipgloss.Color("#88C0D0")
	ColorSuccess   = lipgloss.Color("#A3BE8C")
	ColorWarning   = lipgloss.Color("#EBCB8B")
	ColorError     = lipgloss.Color("#BF616A")
	ColorText      = lipgloss.Color("#ECEFF4")
	ColorSubtle    = lipgloss.Color("#D8DEE9")
	ColorBorder    = lipgloss.Color("#4C566A")
	ColorDim       = lipgloss.Color("#3B4252")
	ColorBg        = lipgloss.Color("#2E3440")
)

// Common styles
var (
	// Base styles
	BaseStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Foreground(ColorText).
		Background(ColorBg)

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
		Background(ColorDim).
		PaddingLeft(2).
		PaddingRight(2)

	// Box styles
	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
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
		BorderForeground(ColorBorder).
		Padding(0, 1)

	// Status styles
	StatusActiveStyle = lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true)

	StatusInactiveStyle = lipgloss.NewStyle().
		Foreground(ColorError).
		Bold(true)

	// Help styles
	HelpStyle = lipgloss.NewStyle().
		Foreground(ColorSubtle)
)
