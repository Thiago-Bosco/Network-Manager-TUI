package network

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Connection represents a network connection
type Connection struct {
	Name      string
	Type      string
	IPAddress string
	SubnetMask string
	Gateway   string
	DNS       string
	Active    bool
}

// ConnectionsMsg is a message containing network connections
type ConnectionsMsg struct {
	Connections []Connection
}

// ConnectionCreatedMsg is a message indicating a connection was created
type ConnectionCreatedMsg struct {
	Connection Connection
}

// ConnectionStatusChangedMsg is a message indicating a connection status changed
type ConnectionStatusChangedMsg struct {
	Name   string
	Active bool
}

// GetConnections returns all configured network connections
func GetConnections() ([]Connection, error) {
	// Check if NetworkManager is available
	_, err := exec.LookPath("nmcli")
	if err == nil {
		return getConnectionsNM()
	}
	
	// Fallback to reading from our own connection files
	return getConnectionsFiles()
}

// CreateConnection creates a new network connection
func CreateConnection(conn Connection) error {
	// Check if NetworkManager is available
	_, err := exec.LookPath("nmcli")
	if err == nil {
		return createConnectionNM(conn)
	}
	
	// Fallback to creating our own connection files
	return createConnectionFile(conn)
}

// ActivateConnection activates a network connection
func ActivateConnection(name string) error {
	// Check if NetworkManager is available
	_, err := exec.LookPath("nmcli")
	if err == nil {
		return activateConnectionNM(name)
	}
	
	// Fallback to activating manually
	return activateConnectionManual(name)
}

// DeactivateConnection deactivates a network connection
func DeactivateConnection(name string) error {
	// Check if NetworkManager is available
	_, err := exec.LookPath("nmcli")
	if err == nil {
		return deactivateConnectionNM(name)
	}
	
	// Fallback to deactivating manually
	return deactivateConnectionManual(name)
}

// getConnectionsNM gets connections using NetworkManager
func getConnectionsNM() ([]Connection, error) {
	// Get all connections
	cmd := exec.Command("nmcli", "-t", "-f", "NAME,TYPE,DEVICE", "connection", "show")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get connections: %w", err)
	}
	
	connections := []Connection{}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		parts := strings.Split(line, ":")
		if len(parts) < 3 {
			continue
		}
		
		name := parts[0]
		connType := parts[1]
		device := parts[2]
		active := device != "" && device != "--"
		
		// Get IP details for active connections
		var ipAddress, subnetMask, gateway, dns string
		if active {
			ipDetails, err := getConnectionDetailsNM(name)
			if err == nil {
				ipAddress = ipDetails.ipAddress
				subnetMask = ipDetails.subnetMask
				gateway = ipDetails.gateway
				dns = ipDetails.dns
			}
		}
		
		connections = append(connections, Connection{
			Name:      name,
			Type:      connType,
			IPAddress: ipAddress,
			SubnetMask: subnetMask,
			Gateway:   gateway,
			DNS:       dns,
			Active:    active,
		})
	}
	
	return connections, nil
}

// connectionDetails holds IP details for a connection
type connectionDetails struct {
	ipAddress  string
	subnetMask string
	gateway    string
	dns        string
}

// getConnectionDetailsNM gets IP details for a connection using NetworkManager
func getConnectionDetailsNM(name string) (connectionDetails, error) {
	var details connectionDetails
	
	// Get connection details
	cmd := exec.Command("nmcli", "-t", "connection", "show", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return details, fmt.Errorf("failed to get connection details: %w", err)
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "IP4.ADDRESS[1]:") {
			// Format: IP4.ADDRESS[1]:192.168.1.100/24
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				addrParts := strings.Split(parts[1], "/")
				if len(addrParts) == 2 {
					details.ipAddress = addrParts[0]
					
					// Convert CIDR to subnet mask
					cidr := addrParts[1]
					subnetMask, err := cidrToMask(cidr)
					if err == nil {
						details.subnetMask = subnetMask
					}
				}
			}
		} else if strings.HasPrefix(line, "IP4.GATEWAY:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				details.gateway = parts[1]
			}
		} else if strings.HasPrefix(line, "IP4.DNS[") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				if details.dns == "" {
					details.dns = parts[1]
				} else {
					details.dns += ", " + parts[1]
				}
			}
		}
	}
	
	return details, nil
}

// createConnectionNM creates a connection using NetworkManager
func createConnectionNM(conn Connection) error {
	args := []string{"connection", "add", "type", conn.Type, "con-name", conn.Name}
	
	// Add interface name if it's ethernet
	if conn.Type == "ethernet" {
		// Find first ethernet interface
		ifaces, err := GetInterfaces()
		if err != nil {
			return fmt.Errorf("failed to get interfaces: %w", err)
		}
		
		for _, iface := range ifaces {
			if iface.Type == "ethernet" {
				args = append(args, "ifname", iface.Name)
				break
			}
		}
	} else if conn.Type == "wifi" {
		// Find first wireless interface
		ifaces, err := GetInterfaces()
		if err != nil {
			return fmt.Errorf("failed to get interfaces: %w", err)
		}
		
		for _, iface := range ifaces {
			if iface.Type == "wireless" {
				args = append(args, "ifname", iface.Name)
				break
			}
		}
	}
	
	// Add IP configuration if provided
	if conn.IPAddress != "" && conn.SubnetMask != "" {
		// Convert subnet mask to CIDR
		cidr, err := maskToCIDR(conn.SubnetMask)
		if err != nil {
			return fmt.Errorf("invalid subnet mask: %w", err)
		}
		
		args = append(args, "ip4", conn.IPAddress+"/"+cidr)
		
		if conn.Gateway != "" {
			args = append(args, "gw4", conn.Gateway)
		}
	}
	
	// Add DNS if provided
	if conn.DNS != "" {
		args = append(args, "ipv4.dns", conn.DNS)
	}
	
	// Create the connection
	cmd := exec.Command("nmcli", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create connection: %s, %w", output, err)
	}
	
	return nil
}

// activateConnectionNM activates a connection using NetworkManager
func activateConnectionNM(name string) error {
	cmd := exec.Command("nmcli", "connection", "up", name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to activate connection: %s, %w", output, err)
	}
	
	return nil
}

// deactivateConnectionNM deactivates a connection using NetworkManager
func deactivateConnectionNM(name string) error {
	cmd := exec.Command("nmcli", "connection", "down", name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to deactivate connection: %s, %w", output, err)
	}
	
	return nil
}

// getConnectionsFiles gets connections from our own connection files
func getConnectionsFiles() ([]Connection, error) {
	// Create connection directory if it doesn't exist
	connectionDir := "/tmp/networkmanager-tui/connections"
	if err := os.MkdirAll(connectionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create connection directory: %w", err)
	}
	
	// List all files in the directory
	files, err := ioutil.ReadDir(connectionDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read connection directory: %w", err)
	}
	
	connections := []Connection{}
	
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			// Read the file
			filePath := filepath.Join(connectionDir, file.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				continue
			}
			
			// Parse the JSON
			var conn Connection
			if err := json.Unmarshal(data, &conn); err != nil {
				continue
			}
			
			// Check if the connection is active
			if conn.Name != "" {
				// For Ethernet connections, check if the interface has the IP
				if conn.Type == "ethernet" && conn.IPAddress != "" {
					ifaces, _ := GetInterfaces()
					for _, iface := range ifaces {
						if iface.IPAddress == conn.IPAddress {
							conn.Active = true
							break
						}
					}
				}
			}
			
			connections = append(connections, conn)
		}
	}
	
	return connections, nil
}

// createConnectionFile creates a connection file
func createConnectionFile(conn Connection) error {
	// Create connection directory if it doesn't exist
	connectionDir := "/tmp/networkmanager-tui/connections"
	if err := os.MkdirAll(connectionDir, 0755); err != nil {
		return fmt.Errorf("failed to create connection directory: %w", err)
	}
	
	// Convert connection to JSON
	data, err := json.Marshal(conn)
	if err != nil {
		return fmt.Errorf("failed to marshal connection: %w", err)
	}
	
	// Create file name from connection name
	fileName := strings.ReplaceAll(conn.Name, " ", "_") + ".json"
	filePath := filepath.Join(connectionDir, fileName)
	
	// Write to file
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write connection file: %w", err)
	}
	
	return nil
}

// activateConnectionManual activates a connection manually
func activateConnectionManual(name string) error {
	// Get connection details
	conn, err := getConnectionByName(name)
	if err != nil {
		return err
	}
	
	// For Ethernet connections, we need to set the IP on an interface
	if conn.Type == "ethernet" {
		// Find first ethernet interface
		ifaces, err := GetInterfaces()
		if err != nil {
			return fmt.Errorf("failed to get interfaces: %w", err)
		}
		
		for _, iface := range ifaces {
			if iface.Type == "ethernet" {
				// Create a temporary interface config
				tmpIface := Interface{
					Name:       iface.Name,
					IPAddress:  conn.IPAddress,
					SubnetMask: conn.SubnetMask,
					Gateway:    conn.Gateway,
					DNS:        conn.DNS,
				}
				
				// Apply the config
				if err := ApplyInterfaceConfig(tmpIface); err != nil {
					return fmt.Errorf("failed to apply interface config: %w", err)
				}
				
				break
			}
		}
	} else if conn.Type == "wifi" {
		// For WiFi, we need to scan and connect
		// (This is simplified - a real implementation would need to handle passwords)
		networks, err := ScanWiFiNetworks()
		if err != nil {
			return fmt.Errorf("failed to scan WiFi networks: %w", err)
		}
		
		for _, network := range networks {
			if network.SSID == name {
				if err := ConnectToWiFi(name, ""); // Assuming open network for simplicity
				if err != nil {
					return fmt.Errorf("failed to connect to WiFi: %w", err)
				}
				break
			}
		}
	}
	
	// Update the connection file to mark it as active
	conn.Active = true
	return createConnectionFile(conn)
}

// deactivateConnectionManual deactivates a connection manually
func deactivateConnectionManual(name string) error {
	// Get connection details
	conn, err := getConnectionByName(name)
	if err != nil {
		return err
	}
	
	// For Ethernet connections, we need to bring down the interface
	if conn.Type == "ethernet" {
		// Find first ethernet interface
		ifaces, err := GetInterfaces()
		if err != nil {
			return fmt.Errorf("failed to get interfaces: %w", err)
		}
		
		for _, iface := range ifaces {
			if iface.Type == "ethernet" && iface.IPAddress == conn.IPAddress {
				// Bring down the interface
				cmd := exec.Command("ip", "link", "set", "dev", iface.Name, "down")
				if err := cmd.Run(); err != nil {
					return fmt.Errorf("failed to bring down interface: %w", err)
				}
				
				break
			}
		}
	} else if conn.Type == "wifi" {
		// Find wireless interface
		cmd := exec.Command("iwconfig")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to run iwconfig: %w", err)
		}
		
		// Find wireless interfaces
		interfaceRe := regexp.MustCompile(`(\w+)\s+IEEE 802\.11`)
		matches := interfaceRe.FindStringSubmatch(string(output))
		if len(matches) >= 2 {
			wirelessIface := matches[1]
			
			// Bring down the interface
			cmd := exec.Command("ip", "link", "set", "dev", wirelessIface, "down")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to bring down interface: %w", err)
			}
		}
	}
	
	// Update the connection file to mark it as inactive
	conn.Active = false
	return createConnectionFile(conn)
}

// getConnectionByName gets a connection by name
func getConnectionByName(name string) (Connection, error) {
	// Get all connections
	connections, err := getConnectionsFiles()
	if err != nil {
		return Connection{}, fmt.Errorf("failed to get connections: %w", err)
	}
	
	// Find the connection with the given name
	for _, conn := range connections {
		if conn.Name == name {
			return conn, nil
		}
	}
	
	return Connection{}, fmt.Errorf("connection not found: %s", name)
}
