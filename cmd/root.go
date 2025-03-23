package cmd

import (
	"fmt"
	"os"

	"github.com/user/networkmanager-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	app := ui.NewApp()
	
	// Run the TUI program
	p := tea.NewProgram(
		app,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithoutMouseAllMotion(), // Disable mouse motion/scroll
	)
	
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		return err
	}

	return nil
}
