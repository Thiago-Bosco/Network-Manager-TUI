package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/user/networkmanager-tui/internal/network"
	"github.com/user/networkmanager-tui/internal/ui"
	"github.com/user/networkmanager-tui/internal/ui/components"
)

// EditConnection is the page for editing a network connection
type EditConnection struct {
	interfaces   components.InterfacesList
	selectedIface *network.Interface
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

// NewEditConnection creates a new edit connection page
func NewEditConnection() EditConnection {
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
	dns.Width = 20
	
	return EditConnection{
		interfaces:   components.NewInterfacesList(),
		selectedIface: nil,
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

// Init initializes the edit connection page
func (e EditConnection) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		func() tea.Msg {
			ifaces, err := network.GetInterfaces()
			if err != nil {
				return network.ErrorMsg{Err: err}
			}
			return network.InterfacesMsg{Interfaces: ifaces}
		},
	)
}

// SetSize sets the size of the edit connection page
func (e *EditConnection) SetSize(width, height int) {
	e.width = width
	e.height = height
	e.interfaces.SetSize(width, height-10) // Adjust for form inputs
}

// updateInputs updates the textinputs with data from the selected interface
func (e *EditConnection) updateInputs() {
	if e.selectedIface == nil {
		e.ipInput.SetValue("")
		e.subnetInput.SetValue("")
		e.gatewayInput.SetValue("")
		e.dnsInput.SetValue("")
		return
	}
	
	e.ipInput.SetValue(e.selectedIface.IPAddress)
	e.subnetInput.SetValue(e.selectedIface.SubnetMask)
	e.gatewayInput.SetValue(e.selectedIface.Gateway)
	e.dnsInput.SetValue(e.selectedIface.DNS)
}

// focusInput sets focus on the current input
func (e *EditConnection) focusInput() {
	inputs := []textinput.Model{e.ipInput, e.subnetInput, e.gatewayInput, e.dnsInput}
	
	for i := range inputs {
		if i == e.focusIndex {
			inputs[i].Focus()
			continue
		}
		inputs[i].Blur()
	}
	
	e.ipInput = inputs[0]
	e.subnetInput = inputs[1]
	e.gatewayInput = inputs[2]
	e.dnsInput = inputs[3]
}

// Update handles user input and updates the page state
func (e EditConnection) Update(msg tea.Msg) (EditConnection, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case network.InterfacesMsg:
		e.interfaces.SetInterfaces(msg.Interfaces)
		
	case network.ErrorMsg:
		e.error = msg.Err.Error()
	
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab":
			if e.selectedIface == nil {
				// No interface selected, tab does nothing
				break
			}
			
			if msg.String() == "tab" {
				e.focusIndex = (e.focusIndex + 1) % 4
			} else {
				e.focusIndex = (e.focusIndex - 1 + 4) % 4
			}
			
			e.focusInput()
			
		case "enter":
			if e.selectedIface == nil {
				// No interface selected yet, so select one
				e.selectedIface = e.interfaces.SelectedInterface()
				if e.selectedIface != nil {
					e.updateInputs()
					e.focusInput()
					e.saved = false
				}
			} else {
				// Save the changes
				e.selectedIface.IPAddress = e.ipInput.Value()
				e.selectedIface.SubnetMask = e.subnetInput.Value()
				e.selectedIface.Gateway = e.gatewayInput.Value()
				e.selectedIface.DNS = e.dnsInput.Value()
				
				// Apply the changes
				cmds = append(cmds, func() tea.Msg {
					err := network.ApplyInterfaceConfig(*e.selectedIface)
					if err != nil {
						return network.ErrorMsg{Err: err}
					}
					return network.ConfigAppliedMsg{Interface: *e.selectedIface}
				})
				
				e.saved = true
			}
		}
	}
	
	// Update interface list
	newInterfaces, cmd := e.interfaces.Update(msg)
	e.interfaces = newInterfaces
	cmds = append(cmds, cmd)
	
	// Update text inputs
	if e.selectedIface != nil {
		var cmd tea.Cmd
		e.ipInput, cmd = e.ipInput.Update(msg)
		cmds = append(cmds, cmd)
		
		e.subnetInput, cmd = e.subnetInput.Update(msg)
		cmds = append(cmds, cmd)
		
		e.gatewayInput, cmd = e.gatewayInput.Update(msg)
		cmds = append(cmds, cmd)
		
		e.dnsInput, cmd = e.dnsInput.Update(msg)
		cmds = append(cmds, cmd)
	}
	
	return e, tea.Batch(cmds...)
}

// View renders the edit connection page
func (e EditConnection) View() string {
	var content string
	
	title := ui.TitleStyle.Copy().Render("Edit Connection")
	
	if e.error != "" {
		errorBox := ui.BoxStyle.Copy().
			BorderForeground(ui.ColorError).
			Render("Error: " + e.error)
			
		content = lipgloss.JoinVertical(lipgloss.Left, title, errorBox, e.interfaces.View())
	} else if e.selectedIface == nil {
		instructions := lipgloss.NewStyle().
			Foreground(ui.ColorSubtle).
			Render("Select a network interface to edit")
			
		content = lipgloss.JoinVertical(lipgloss.Left, title, instructions, e.interfaces.View())
	} else {
		// Form layout
		ifaceInfo := fmt.Sprintf("Editing interface: %s (%s)", 
			ui.SelectedItemStyle.Render(e.selectedIface.Name), 
			e.selectedIface.MACAddress)
			
		ipLabel := "IP Address:"
		subnetLabel := "Subnet Mask:"
		gatewayLabel := "Gateway:"
		dnsLabel := "DNS Servers:"
		
		ipField := lipgloss.JoinHorizontal(lipgloss.Left, 
			lipgloss.NewStyle().Width(15).Render(ipLabel), 
			e.ipInput.View())
			
		subnetField := lipgloss.JoinHorizontal(lipgloss.Left, 
			lipgloss.NewStyle().Width(15).Render(subnetLabel), 
			e.subnetInput.View())
			
		gatewayField := lipgloss.JoinHorizontal(lipgloss.Left, 
			lipgloss.NewStyle().Width(15).Render(gatewayLabel), 
			e.gatewayInput.View())
			
		dnsField := lipgloss.JoinHorizontal(lipgloss.Left, 
			lipgloss.NewStyle().Width(15).Render(dnsLabel), 
			e.dnsInput.View())
			
		form := lipgloss.JoinVertical(lipgloss.Left, ipField, subnetField, gatewayField, dnsField)
		
		helpText := "Tab: Next field • Shift+Tab: Previous field • Enter: Save changes"
		help := ui.HelpStyle.Render(helpText)
		
		saveMsg := ""
		if e.saved {
			saveMsg = ui.StatusActiveStyle.Render("Changes saved successfully!")
		}
		
		content = lipgloss.JoinVertical(lipgloss.Left, 
			title, 
			lipgloss.NewStyle().PaddingTop(1).Render(ifaceInfo),
			lipgloss.NewStyle().PaddingTop(1).Render(form),
			lipgloss.NewStyle().PaddingTop(1).Render(saveMsg),
			lipgloss.NewStyle().PaddingTop(1).Render(help))
	}
	
	return content
}
