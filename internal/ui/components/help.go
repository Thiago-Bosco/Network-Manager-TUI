package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	
	"github.com/user/networkmanager-tui/internal/ui"
)

// HelpItem represents a single help item
type HelpItem struct {
	Key         string
	Description string
}

// HelpBox displays help information in a styled box
type HelpBox struct {
	items []HelpItem
	width int
}

// NewHelpBox creates a new help box component
func NewHelpBox() HelpBox {
	return HelpBox{
		items: defaultHelpItems(),
		width: 80,
	}
}

// SetWidth sets the width of the help box
func (h *HelpBox) SetWidth(width int) {
	h.width = width
}

// SetItems replaces the help items
func (h *HelpBox) SetItems(items []HelpItem) {
	h.items = items
}

// View renders the help box
func (h HelpBox) View() string {
	var sb strings.Builder
	
	// Format each help item
	for i, item := range h.items {
		keyStyle := lipgloss.NewStyle().
			Foreground(ui.ColorPrimary).
			Bold(true).
			PaddingRight(1)
		
		descStyle := lipgloss.NewStyle().
			Foreground(ui.ColorText)
		
		keyText := keyStyle.Render(item.Key)
		descText := descStyle.Render(item.Description)
		
		helpText := lipgloss.JoinHorizontal(lipgloss.Left, keyText, descText)
		
		sb.WriteString(helpText)
		
		// Add separator between items, but not after the last one
		if i < len(h.items)-1 {
			sb.WriteString("  ")
		}
	}
	
	// Wrap the help text in a box
	return ui.BoxStyle.Copy().
		Width(h.width - 4).
		BorderForeground(ui.ColorBorder).
		Render(sb.String())
}

// defaultHelpItems returns the default help items
func defaultHelpItems() []HelpItem {
	return []HelpItem{
		{"↑/↓", "Navigate"},
		{"Enter", "Select"},
		{"Esc", "Back"},
		{"q", "Quit"},
		{"?", "Help"},
	}
}
