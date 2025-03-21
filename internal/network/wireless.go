package network

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// WiFiNetwork represents a wireless network
type WiFiNetwork struct {
	SSID           string
	BSSID          string
	SignalStrength int  // 0 to 5, with 5 being the strongest
	Secured        bool
}

// WiFiScannedMsg is a message containing scanned WiFi networks
type WiFiScannedMsg struct {
	Networks []WiFiNetwork
}

// WiFiConnectedMsg is a message indicating a WiFi connection was established
type WiFiConnectedMsg struct {
	SSID string
}

// ScanWiFiNetworks scans for available WiFi networks
func ScanWiFiNetworks() ([]WiFiNetwork, error) {
	// First check if we have a wireless interface
	cmd := exec.Command("iwconfig")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run iwconfig: %w", err)
	}

	// Find wireless interfaces
	interfaceRe := regexp.MustCompile(`(\w+)\s+IEEE 802\.11`)
	matches := interfaceRe.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return nil, fmt.Errorf("no wireless interfaces found")
	}

	wirelessIface := matches[1]

	// Scan for networks
	scanCmd := exec.Command("iwlist", wirelessIface, "scan")
	scanOutput, err := scanCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to scan for networks: %w", err)
	}

	// Parse the output
	return parseWiFiScan(string(scanOutput)), nil
}

// ConnectToWiFi connects to a WiFi network
func ConnectToWiFi(ssid, password string) error {
	// First check if we have a wireless interface
	cmd := exec.Command("iwconfig")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run iwconfig: %w", err)
	}

	// Find wireless interfaces
	interfaceRe := regexp.MustCompile(`(\w+)\s+IEEE 802\.11`)
	matches := interfaceRe.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		return fmt.Errorf("no wireless interfaces found")
	}

	wirelessIface := matches[1]

	// Create a temporary wpa_supplicant configuration
	configContent := fmt.Sprintf(`
network={
    ssid="%s"
    %s
}
`, ssid, getAuthConfig(password))

	// Write config to temporary file
	configFile := "/tmp/wpa_supplicant.conf"
	writeCmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' > %s", configContent, configFile))
	if err := writeCmd.Run(); err != nil {
		return fmt.Errorf("failed to write wpa_supplicant config: %w", err)
	}

	// Stop any running wpa_supplicant instances
	stopCmd := exec.Command("killall", "wpa_supplicant")
	// Ignore errors - it might not be running
	stopCmd.Run()

	// Start wpa_supplicant
	wpaCmd := exec.Command("wpa_supplicant", "-B", "-i", wirelessIface, "-c", configFile)
	if err := wpaCmd.Run(); err != nil {
		return fmt.Errorf("failed to start wpa_supplicant: %w", err)
	}

	// Get IP via DHCP
	dhcpCmd := exec.Command("dhclient", "-v", wirelessIface)
	if err := dhcpCmd.Run(); err != nil {
		return fmt.Errorf("failed to get IP address: %w", err)
	}

	return nil
}

// parseWiFiScan parses the output of iwlist scan
func parseWiFiScan(scanOutput string) []WiFiNetwork {
	networks := []WiFiNetwork{}
	
	// Split output into cells
	cells := strings.Split(scanOutput, "Cell ")
	
	// Skip the first element, which is just header text
	cells = cells[1:]
	
	for _, cell := range cells {
		var network WiFiNetwork
		
		// Extract SSID
		ssidRe := regexp.MustCompile(`ESSID:"([^"]*)"`)
		ssidMatch := ssidRe.FindStringSubmatch(cell)
		if len(ssidMatch) < 2 {
			continue
		}
		network.SSID = ssidMatch[1]
		
		// Skip hidden networks
		if network.SSID == "" {
			continue
		}
		
		// Extract BSSID
		bssidRe := regexp.MustCompile(`Address: ([0-9A-F:]{17})`)
		bssidMatch := bssidRe.FindStringSubmatch(cell)
		if len(bssidMatch) < 2 {
			continue
		}
		network.BSSID = bssidMatch[1]
		
		// Extract signal strength
		signalRe := regexp.MustCompile(`Signal level=(-\d+) dBm`)
		signalMatch := signalRe.FindStringSubmatch(cell)
		if len(signalMatch) < 2 {
			// Alternative format
			signalRe = regexp.MustCompile(`Quality=(\d+)/70`)
			signalMatch = signalRe.FindStringSubmatch(cell)
			if len(signalMatch) < 2 {
				// Use default value
				network.SignalStrength = 3
			} else {
				quality, _ := strconv.Atoi(signalMatch[1])
				network.SignalStrength = mapSignalQuality(quality, 0, 70, 0, 5)
			}
		} else {
			// Convert dBm to our scale (0-5)
			signalDBm, _ := strconv.Atoi(signalMatch[1])
			network.SignalStrength = mapSignalStrength(signalDBm, -100, -50, 0, 5)
		}
		
		// Check encryption
		encryptionRe := regexp.MustCompile(`Encryption key:(.*)`)
		encryptionMatch := encryptionRe.FindStringSubmatch(cell)
		if len(encryptionMatch) < 2 || strings.TrimSpace(encryptionMatch[1]) == "off" {
			network.Secured = false
		} else {
			network.Secured = true
		}
		
		networks = append(networks, network)
	}
	
	return networks
}

// getAuthConfig returns the WPA configuration based on password
func getAuthConfig(password string) string {
	if password == "" {
		return "key_mgmt=NONE"
	}
	return fmt.Sprintf("psk=\"%s\"", password)
}

// mapSignalStrength maps a signal level (dBm) to a range
func mapSignalStrength(value, fromLow, fromHigh, toLow, toHigh int) int {
	// Ensure value is in range
	if value < fromLow {
		value = fromLow
	}
	if value > fromHigh {
		value = fromHigh
	}
	
	// Map to new range
	return (value-fromLow)*(toHigh-toLow)/(fromHigh-fromLow) + toLow
}

// mapSignalQuality maps a quality value to a range
func mapSignalQuality(value, fromLow, fromHigh, toLow, toHigh int) int {
	// Ensure value is in range
	if value < fromLow {
		value = fromLow
	}
	if value > fromHigh {
		value = fromHigh
	}
	
	// Map to new range
	return (value-fromLow)*(toHigh-toLow)/(fromHigh-fromLow) + toLow
}
