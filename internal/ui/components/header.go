package components

import (
	"github.com/charmbracelet/lipgloss"
	
	"github.com/user/networkmanager-tui/internal/ui"
)

// Header is the top component of the application
type Header struct {
	width int
}

// NewHeader creates a new header component
func NewHeader() Header {
	return Header{
		width: 80, // Default width
	}
}

// SetWidth sets the width of the header
func (h *Header) SetWidth(width int) {
	h.width = width
}

// Height returns the height of the header
func (h Header) Height() int {
	return 3 // Title + separator line + space
}

// View renders the header
func (h Header) View() string {
	title := ui.TitleStyle.Copy().Width(h.width).Render("Network Manager TUI")
	
	// Create a subtle separator line
	separator := lipgloss.NewStyle().
		Foreground(ui.ColorBorder).
		Width(h.width).
		Render("─" + lipgloss.RepeatString("─", h.width-2) + "─")
	
	return lipgloss.JoinVertical(lipgloss.Center, title, separator)
}
