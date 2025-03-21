package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"

	"github.com/user/networkmanager-tui/cmd"
)

func main() {
	// Setup logger
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
	})

	// Check root privileges
	if os.Geteuid() != 0 {
		fmt.Println("This application requires root privileges to configure network interfaces.")
		fmt.Println("Please run with sudo or as root.")
		os.Exit(1)
	}

	// Execute the root command
	if err := cmd.Execute(); err != nil {
		logger.Error("application error", "err", err)
		os.Exit(1)
	}
}
