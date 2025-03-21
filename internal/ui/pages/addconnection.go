package pages

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/user/networkmanager-tui/internal/network"
	"github.com/user/networkmanager-tui/internal/ui"
)

// ConnectionType represents a type of network connection
type ConnectionType int

const (
	Ethernet ConnectionType = iota
	WiFi
	VPN
)

// AddConnection is the page for adding a new network connection
type AddConnection struct {
	nameInput    textinput.Model
	typeIndex    int
	ipInput      textinput.Model
	subnetInput  textinput.Model
	gatewayInput textinput.Model
	dnsInput     textinput.Model
	focusIndex   int
	width        int
	height       int
	saved        bool
	error        string
}

// NewAddConnection creates a new add connection page
func NewAddConnection() AddConnection {
	name := textinput.New()
	name.Placeholder = "New Connection"
	name.CharLimit = 50
	name.Width = 30
	name.Focus()
	
	ip := textinput.New()
	ip.Placeholder = "192.168.1.10"
	ip.CharLimit = 15
	ip.Width = 20
	
	subnet := textinput.New()
	subnet.Placeholder = "255.255.255.0"
	subnet.CharLimit = 15
	subnet.Width = 20
	
	gateway := textinput.New()
	gateway.Placeholder = "192.168.1.1"
	gateway.CharLimit = 15
	gateway.Width = 20
	
	dns := textinput.New()
	dns.Placeholder = "8.8.8.8, 1.1.1.1"
	dns.CharLimit = 50
	dns.Width = 30
	
	return AddConnection{
		nameInput:    name,
		typeIndex:    0,
		ipInput:      ip,
		subnetInput:  subnet,
		gatewayInput: gateway,
		dnsInput:     dns,
		focusIndex:   0,
		width:        80,
		height:       20,
		saved:        false,
		error:        "",
	}
}

// Init initializes the add connection page
func (a AddConnection) Init() tea.Cmd {
	return textinput.Blink
}

// SetSize sets the size of the add connection page
func (a *AddConnection) SetSize(width, height int) {
	a.width = width
	a.height = height
}

// focusInput sets focus on the current input
func (a *AddConnection) focusInput() {
	inputs := []textinput.Model{a.nameInput, a.ipInput, a.subnetInput, a.gatewayInput, a.dnsInput}
	
	for i := range inputs {
		if i == a.focusIndex {
			inputs[i].Focus()
			continue
		}
		inputs[i].Blur()
	}
	
	a.nameInput = inputs[0]
	a.ipInput = inputs[1]
	a.subnetInput = inputs[2]
	a.gatewayInput = inputs[3]
	a.dnsInput = inputs[4]
}

// Update handles user input and updates the page state
func (a AddConnection) Update(msg tea.Msg) (AddConnection, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case network.ErrorMsg:
		a.error = msg.Err.Error()
	
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
			a.focusIndex = (a.focusIndex + 1) % 5
			a.focusInput()
			
		case key.Matches(msg, key.NewBinding(key.WithKeys("shift+tab"))):
			a.focusIndex = (a.focusIndex - 1 + 5) % 5
			a.focusInput()
			
		case key.Matches(msg, key.NewBinding(key.WithKeys("left"))):
			if a.focusIndex == -1 { // Type selection mode
				a.typeIndex = (a.typeIndex - 1 + 3) % 3
			}
			
		case key.Matches(msg, key.NewBinding(key.WithKeys("right"))):
			if a.focusIndex == -1 { // Type selection mode
				a.typeIndex = (a.typeIndex + 1) % 3
			}
			
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if strings.TrimSpace(a.nameInput.Value()) == "" {
				a.error = "Connection name cannot be empty"
				break
			}
			
			// Create a new connection
			conn := network.Connection{
				Name:      a.nameInput.Value(),
				Type:      getConnectionTypeString(ConnectionType(a.typeIndex)),
				IPAddress: a.ipInput.Value(),
				SubnetMask: a.subnetInput.Value(),
				Gateway:   a.gatewayInput.Value(),
				DNS:       a.dnsInput.Value(),
			}
			
			cmds = append(cmds, func() tea.Msg {
				err := network.CreateConnection(conn)
				if err != nil {
					return network.ErrorMsg{Err: err}
				}
				return network.ConnectionCreatedMsg{Connection: conn}
			})
			
			a.saved = true
		}
	}
	
	// Update text inputs
	var cmd tea.Cmd
	a.nameInput, cmd = a.nameInput.Update(msg)
	cmds = append(cmds, cmd)
	
	a.ipInput, cmd = a.ipInput.Update(msg)
	cmds = append(cmds, cmd)
	
	a.subnetInput, cmd = a.subnetInput.Update(msg)
	cmds = append(cmds, cmd)
	
	a.gatewayInput, cmd = a.gatewayInput.Update(msg)
	cmds = append(cmds, cmd)
	
	a.dnsInput, cmd = a.dnsInput.Update(msg)
	cmds = append(cmds, cmd)
	
	return a, tea.Batch(cmds...)
}

// View renders the add connection page
func (a AddConnection) View() string {
	title := ui.TitleStyle.Copy().Render("Add Connection")
	
	// Connection type selection
	typeTitle := "Connection Type:"
	ethernetBtn := getTypeButton("Ethernet", ConnectionType(a.typeIndex) == Ethernet)
	wifiBtn := getTypeButton("WiFi", ConnectionType(a.typeIndex) == WiFi)
	vpnBtn := getTypeButton("VPN", ConnectionType(a.typeIndex) == VPN)
	
	typeSelection := lipgloss.JoinHorizontal(lipgloss.Left, 
		lipgloss.NewStyle().Width(18).Render(typeTitle),
		ethernetBtn, wifiBtn, vpnBtn)
	
	// Form fields
	nameLabel := "Connection Name:"
	ipLabel := "IP Address:"
	subnetLabel := "Subnet Mask:"
	gatewayLabel := "Gateway:"
	dnsLabel := "DNS Servers:"
	
	nameField := lipgloss.JoinHorizontal(lipgloss.Left, 
		lipgloss.NewStyle().Width(18).Render(nameLabel), 
		a.nameInput.View())
		
	ipField := lipgloss.JoinHorizontal(lipgloss.Left, 
		lipgloss.NewStyle().Width(18).Render(ipLabel), 
		a.ipInput.View())
		
	subnetField := lipgloss.JoinHorizontal(lipgloss.Left, 
		lipgloss.NewStyle().Width(18).Render(subnetLabel), 
		a.subnetInput.View())
		
	gatewayField := lipgloss.JoinHorizontal(lipgloss.Left, 
		lipgloss.NewStyle().Width(18).Render(gatewayLabel), 
		a.gatewayInput.View())
		
	dnsField := lipgloss.JoinHorizontal(lipgloss.Left, 
		lipgloss.NewStyle().Width(18).Render(dnsLabel), 
		a.dnsInput.View())
	
	form := lipgloss.JoinVertical(lipgloss.Left, 
		nameField, 
		lipgloss.NewStyle().PaddingTop(1).Render(typeSelection),
		lipgloss.NewStyle().PaddingTop(1).Render(ipField),
		subnetField,
		gatewayField,
		dnsField)
	
	// Help text
	helpText := "Tab: Next field • Shift+Tab: Previous field • Enter: Save connection"
	help := ui.HelpStyle.Render(helpText)
	
	// Message area
	var message string
	if a.error != "" {
		message = ui.StatusInactiveStyle.Render("Error: " + a.error)
	} else if a.saved {
		message = ui.StatusActiveStyle.Render("Connection created successfully!")
	}
	
	// Combine all elements
	content := lipgloss.JoinVertical(lipgloss.Left, 
		title,
		lipgloss.NewStyle().PaddingTop(1).Render(form),
		lipgloss.NewStyle().PaddingTop(1).Render(message),
		lipgloss.NewStyle().PaddingTop(1).Render(help))
	
	return lipgloss.NewStyle().
		Width(a.width).
		Padding(1, 2).
		Render(content)
}

// getTypeButton returns a styled button for the connection type
func getTypeButton(label string, selected bool) string {
	if selected {
		return ui.ButtonStyle.Copy().
			Background(ui.ColorPrimary).
			Render(label)
	}
	return ui.ButtonStyle.Copy().
		Background(ui.ColorDim).
		Render(label)
}

// getConnectionTypeString converts a ConnectionType to its string representation
func getConnectionTypeString(t ConnectionType) string {
	switch t {
	case Ethernet:
		return "ethernet"
	case WiFi:
		return "wifi"
	case VPN:
		return "vpn"
	default:
		return "ethernet"
	}
}
