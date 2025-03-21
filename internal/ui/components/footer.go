package components

import (
	"runtime"
	"time"

	"github.com/charmbracelet/lipgloss"
	
	"github.com/user/networkmanager-tui/internal/ui"
)

// Footer is the bottom component of the application
type Footer struct {
	width int
}

// NewFooter creates a new footer component
func NewFooter() Footer {
	return Footer{
		width: 80, // Default width
	}
}

// SetWidth sets the width of the footer
func (f *Footer) SetWidth(width int) {
	f.width = width
}

// Height returns the height of the footer
func (f Footer) Height() int {
	return 1
}

// View renders the footer
func (f Footer) View() string {
	now := time.Now().Format("15:04:05")
	leftText := "Network Manager TUI"
	rightText := "Go " + runtime.Version() + " | " + now
	
	leftSide := ui.FooterStyle.Copy().
		Width((f.width / 2) - 2).
		Align(lipgloss.Left).
		Render(leftText)
	
	rightSide := ui.FooterStyle.Copy().
		Width((f.width / 2) - 2).
		Align(lipgloss.Right).
		Render(rightText)
	
	return lipgloss.JoinHorizontal(lipgloss.Center, leftSide, rightSide)
}
