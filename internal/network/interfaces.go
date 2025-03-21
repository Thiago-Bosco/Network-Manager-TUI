package network

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strings"
)

// Interface represents a network interface
type Interface struct {
	Name       string
	Type       string
	IPAddress  string
	SubnetMask string
	Gateway    string
	DNS        string
	MACAddress string
	IsActive   bool
}

// InterfacesMsg is a message containing network interfaces
type InterfacesMsg struct {
	Interfaces []Interface
}

// ErrorMsg is a message containing an error
type ErrorMsg struct {
	Err error
}

// ConfigAppliedMsg is a message indicating a configuration was applied
type ConfigAppliedMsg struct {
	Interface Interface
}

// GetInterfaces returns all network interfaces on the system
func GetInterfaces() ([]Interface, error) {
	// Get physical interfaces
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	// Run ip addr to get additional info
	ipAddrCmd := exec.Command("ip", "addr")
	ipAddrOutput, err := ipAddrCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run 'ip addr': %w", err)
	}

	// Run ip route to get gateway info
	ipRouteCmd := exec.Command("ip", "route")
	ipRouteOutput, err := ipRouteCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run 'ip route': %w", err)
	}

	// Get DNS info
	dnsInfo, err := getDNSServers()
	if err != nil {
		// Non-fatal, just use empty string
		dnsInfo = ""
	}

	interfaces := []Interface{}

	for _, iface := range netInterfaces {
		// Skip loopback interfaces
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Check if interface is up
		isActive := iface.Flags&net.FlagUp != 0

		// Extract information from ip addr output
		ipInfo := extractIPInfo(string(ipAddrOutput), iface.Name)
		
		// Extract gateway from ip route output
		gateway := extractGateway(string(ipRouteOutput), iface.Name)

		// Determine interface type
		ifaceType := determineInterfaceType(iface.Name)

		// Create interface object
		netInterface := Interface{
			Name:       iface.Name,
			Type:       ifaceType,
			IPAddress:  ipInfo.ipAddress,
			SubnetMask: ipInfo.subnetMask,
			Gateway:    gateway,
			DNS:        dnsInfo,
			MACAddress: iface.HardwareAddr.String(),
			IsActive:   isActive,
		}

		interfaces = append(interfaces, netInterface)
	}

	return interfaces, nil
}

// Apply a network interface configuration
func ApplyInterfaceConfig(iface Interface) error {
	// Validate IP address
	if iface.IPAddress != "" {
		if net.ParseIP(iface.IPAddress) == nil {
			return fmt.Errorf("invalid IP address: %s", iface.IPAddress)
		}
	}

	// Validate subnet mask
	if iface.SubnetMask != "" {
		if net.ParseIP(iface.SubnetMask) == nil {
			return fmt.Errorf("invalid subnet mask: %s", iface.SubnetMask)
		}
	}

	// Set IP address and subnet
	if iface.IPAddress != "" && iface.SubnetMask != "" {
		// Convert subnet mask to CIDR notation
		cidr, err := maskToCIDR(iface.SubnetMask)
		if err != nil {
			return fmt.Errorf("failed to convert subnet mask to CIDR: %w", err)
		}

		// First flush existing IP addresses
		flushCmd := exec.Command("ip", "addr", "flush", "dev", iface.Name)
		if err := flushCmd.Run(); err != nil {
			return fmt.Errorf("failed to flush IP addresses: %w", err)
		}

		// Set new IP address
		addrCmd := exec.Command("ip", "addr", "add", iface.IPAddress+"/"+cidr, "dev", iface.Name)
		if err := addrCmd.Run(); err != nil {
			return fmt.Errorf("failed to set IP address: %w", err)
		}
	}

	// Set default gateway if provided
	if iface.Gateway != "" {
		// First delete existing default route
		delRouteCmd := exec.Command("ip", "route", "del", "default", "dev", iface.Name)
		// Ignore errors - the route might not exist
		delRouteCmd.Run()

		// Add new default route
		addRouteCmd := exec.Command("ip", "route", "add", "default", "via", iface.Gateway, "dev", iface.Name)
		if err := addRouteCmd.Run(); err != nil {
			return fmt.Errorf("failed to set default gateway: %w", err)
		}
	}

	// Set DNS servers if provided
	if iface.DNS != "" {
		if err := setDNSServers(iface.DNS); err != nil {
			return fmt.Errorf("failed to set DNS servers: %w", err)
		}
	}

	// Make sure the interface is up
	upCmd := exec.Command("ip", "link", "set", "dev", iface.Name, "up")
	if err := upCmd.Run(); err != nil {
		return fmt.Errorf("failed to bring interface up: %w", err)
	}

	return nil
}

// IPInfo represents extracted IP information
type ipInfo struct {
	ipAddress  string
	subnetMask string
}

// extractIPInfo extracts IP information from ip addr output
func extractIPInfo(ipAddrOutput, interfaceName string) ipInfo {
	var info ipInfo

	// Find the section for the specific interface
	pattern := fmt.Sprintf(`(?m)^\d+: %s:.*?(?:\n\d+:|$)`, regexp.QuoteMeta(interfaceName))
	re := regexp.MustCompile(pattern)
	ifaceSection := re.FindString(ipAddrOutput)

	// Extract IPv4 address and CIDR
	ipv4Pattern := `inet\s+(\d+\.\d+\.\d+\.\d+)/(\d+)`
	ipv4Re := regexp.MustCompile(ipv4Pattern)
	ipv4Match := ipv4Re.FindStringSubmatch(ifaceSection)

	if len(ipv4Match) >= 3 {
		info.ipAddress = ipv4Match[1]
		
		// Convert CIDR to subnet mask
		cidr := ipv4Match[2]
		subnetMask, err := cidrToMask(cidr)
		if err == nil {
			info.subnetMask = subnetMask
		}
	}

	return info
}

// extractGateway extracts the default gateway for an interface
func extractGateway(ipRouteOutput, interfaceName string) string {
	// Look for default route
	pattern := fmt.Sprintf(`default via (\d+\.\d+\.\d+\.\d+) dev %s`, regexp.QuoteMeta(interfaceName))
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(ipRouteOutput)

	if len(match) >= 2 {
		return match[1]
	}

	return ""
}

// determineInterfaceType guesses the interface type based on its name
func determineInterfaceType(name string) string {
	if strings.HasPrefix(name, "wl") {
		return "wireless"
	} else if strings.HasPrefix(name, "en") || strings.HasPrefix(name, "eth") {
		return "ethernet"
	} else if strings.HasPrefix(name, "tun") || strings.HasPrefix(name, "tap") {
		return "vpn"
	} else if strings.HasPrefix(name, "docker") || strings.HasPrefix(name, "br") {
		return "virtual"
	}
	return "unknown"
}

// getDNSServers reads DNS server information from resolv.conf
func getDNSServers() (string, error) {
	// Read /etc/resolv.conf
	cmd := exec.Command("cat", "/etc/resolv.conf")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to read resolv.conf: %w", err)
	}

	// Extract nameserver lines
	pattern := `nameserver\s+(\d+\.\d+\.\d+\.\d+)`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(string(output), -1)

	servers := []string{}
	for _, match := range matches {
		if len(match) >= 2 {
			servers = append(servers, match[1])
		}
	}

	return strings.Join(servers, ", "), nil
}

// setDNSServers writes DNS server information to resolv.conf
func setDNSServers(dnsServers string) error {
	// Parse DNS servers
	servers := strings.Split(dnsServers, ",")
	
	// Build new resolv.conf content
	var content strings.Builder
	content.WriteString("# Generated by Network Manager TUI\n")
	
	for _, server := range servers {
		server = strings.TrimSpace(server)
		if net.ParseIP(server) != nil {
			content.WriteString(fmt.Sprintf("nameserver %s\n", server))
		}
	}
	
	// Add some common options
	content.WriteString("options timeout:2 attempts:3\n")
	
	// Write to resolv.conf
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' > /etc/resolv.conf", content.String()))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to write to resolv.conf: %w", err)
	}
	
	return nil
}

// cidrToMask converts CIDR notation to subnet mask
func cidrToMask(cidr string) (string, error) {
	cidrInt, err := parseInt(cidr)
	if err != nil {
		return "", err
	}
	
	// Create a netmask
	mask := net.CIDRMask(cidrInt, 32)
	
	// Convert to string format
	return fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3]), nil
}

// maskToCIDR converts a subnet mask to CIDR notation
func maskToCIDR(mask string) (string, error) {
	// Parse the mask
	ip := net.ParseIP(mask)
	if ip == nil {
		return "", fmt.Errorf("invalid subnet mask")
	}
	
	// Count the bits
	bits := 0
	for _, b := range ip.To4() {
		bits += countBits(b)
	}
	
	return fmt.Sprintf("%d", bits), nil
}

// countBits counts the number of bits set in a byte
func countBits(b byte) int {
	count := 0
	for i := 0; i < 8; i++ {
		if (b & (1 << uint(i))) != 0 {
			count++
		}
	}
	return count
}

// parseInt converts a string to an integer
func parseInt(s string) (int, error) {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
