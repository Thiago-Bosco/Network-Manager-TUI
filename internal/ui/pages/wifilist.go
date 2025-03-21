package pages

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/user/networkmanager-tui/internal/network"
	"github.com/user/networkmanager-tui/internal/ui"
)

// WiFiNetwork represents a Wi-Fi network
type WiFiNetwork struct {
	SSID       string
	SignalStrength int
	Secured    bool
}

// Title returns the SSID of the Wi-Fi network
func (w WiFiNetwork) Title() string { 
	if w.Secured {
		return fmt.Sprintf("%s 🔒", w.SSID)
	}
	return w.SSID 
}

// Description returns the signal strength of the Wi-Fi network
func (w WiFiNetwork) Description() string { 
	bars := ""
	for i := 0; i < 5; i++ {
		if i < w.SignalStrength {
			bars += "█"
		} else {
			bars += "░"
		}
	}
	return fmt.Sprintf("Signal: %s", bars)
}

// FilterValue returns the SSID of the Wi-Fi network for filtering
func (w WiFiNetwork) FilterValue() string { return w.SSID }

// WiFiList is the page for listing and connecting to Wi-Fi networks
type WiFiList struct {
	list           list.Model
	selectedSSID   string
	passwordInput  textinput.Model
	connecting     bool
	connected      bool
	showPassword   bool
	width          int
	height         int
	error          string
	networks       []network.WiFiNetwork
}

// NewWiFiList creates a new Wi-Fi list page
func NewWiFiList() WiFiList {
	passwordInput := textinput.New()
	passwordInput.Placeholder = "Enter password"
	passwordInput.Width = 30
	passwordInput.EchoMode = textinput.EchoPassword
	
	return WiFiList{
		list:          list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		selectedSSID:  "",
		passwordInput: passwordInput,
		connecting:    false,
		connected:     false,
		showPassword:  false,
		width:         80,
		height:        20,
		error:         "",
		networks:      []network.WiFiNetwork{},
	}
}

// Init initializes the Wi-Fi list page
func (w WiFiList) Init() tea.Cmd {
	w.list.Title = "Available WiFi Networks"
	w.list.SetShowStatusBar(false)
	w.list.SetFilteringEnabled(false)
	w.list.Styles.Title = ui.TitleStyle
	w.list.Styles.HelpStyle = ui.HelpStyle
	
	return func() tea.Msg {
		networks, err := network.ScanWiFiNetworks()
		if err != nil {
			return network.ErrorMsg{Err: err}
		}
		return network.WiFiScannedMsg{Networks: networks}
	}
}

// SetSize sets the size of the Wi-Fi list page
func (w *WiFiList) SetSize(width, height int) {
	w.width = width
	w.height = height
	w.list.SetSize(width, height-5) // Adjust for password input area
}

// Update handles user input and updates the page state
func (w WiFiList) Update(msg tea.Msg) (WiFiList, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case network.WiFiScannedMsg:
		w.networks = msg.Networks
		items := make([]list.Item, len(msg.Networks))
		for i, network := range msg.Networks {
			items[i] = WiFiNetwork{
				SSID:          network.SSID,
				SignalStrength: network.SignalStrength,
				Secured:       network.Secured,
			}
		}
		
		// Sort networks by signal strength (highest first)
		sort.Slice(items, func(i, j int) bool {
			return items[i].(WiFiNetwork).SignalStrength > items[j].(WiFiNetwork).SignalStrength
		})
		
		w.list.SetItems(items)
		
	case network.ErrorMsg:
		w.error = msg.Err.Error()
		w.connecting = false
		
	case network.WiFiConnectedMsg:
		w.connected = true
		w.connecting = false
		w.selectedSSID = msg.SSID
		
	case tea.KeyMsg:
		// Handle keyboard shortcuts
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if w.selectedSSID == "" {
				// No network selected yet, select the current one
				if w.list.SelectedItem() != nil {
					network := w.list.SelectedItem().(WiFiNetwork)
					w.selectedSSID = network.SSID
					w.showPassword = network.Secured
					if w.showPassword {
						w.passwordInput.Focus()
					} else {
						// Connect to open network
						w.connecting = true
						cmds = append(cmds, func() tea.Msg {
							err := network.ConnectToWiFi(w.selectedSSID, "")
							if err != nil {
								return network.ErrorMsg{Err: err}
							}
							return network.WiFiConnectedMsg{SSID: w.selectedSSID}
						})
					}
				}
			} else if w.showPassword {
				// Connect with password
				password := w.passwordInput.Value()
				w.connecting = true
				cmds = append(cmds, func() tea.Msg {
					err := network.ConnectToWiFi(w.selectedSSID, password)
					if err != nil {
						return network.ErrorMsg{Err: err}
					}
					return network.WiFiConnectedMsg{SSID: w.selectedSSID}
				})
			}
			
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			if w.selectedSSID != "" {
				// Go back to list view
				w.selectedSSID = ""
				w.showPassword = false
				w.connecting = false
				w.passwordInput.Reset()
			}
			
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			// Refresh network list
			w.connecting = false
			w.connected = false
			w.selectedSSID = ""
			w.error = ""
			cmds = append(cmds, func() tea.Msg {
				networks, err := network.ScanWiFiNetworks()
				if err != nil {
					return network.ErrorMsg{Err: err}
				}
				return network.WiFiScannedMsg{Networks: networks}
			})
		}
	}
	
	// Update list model
	if w.selectedSSID == "" && !w.connecting && !w.connected {
		var cmd tea.Cmd
		w.list, cmd = w.list.Update(msg)
		cmds = append(cmds, cmd)
	}
	
	// Update password input if relevant
	if w.showPassword && !w.connecting && !w.connected {
		var cmd tea.Cmd
		w.passwordInput, cmd = w.passwordInput.Update(msg)
		cmds = append(cmds, cmd)
	}
	
	return w, tea.Batch(cmds...)
}

// View renders the Wi-Fi list page
func (w WiFiList) View() string {
	title := ui.TitleStyle.Copy().Render("WiFi Networks")
	
	var content string
	
	if w.error != "" {
		errorBox := ui.BoxStyle.Copy().
			BorderForeground(ui.ColorError).
			Render("Error: " + w.error)
			
		refreshHelp := ui.HelpStyle.Render("Press 'r' to refresh")
		
		content = lipgloss.JoinVertical(lipgloss.Left, 
			title, 
			errorBox,
			refreshHelp)
	} else if w.connected {
		connectedMsg := fmt.Sprintf("Connected to: %s", ui.SelectedItemStyle.Render(w.selectedSSID))
		statusBox := ui.BoxStyle.Copy().
			BorderForeground(ui.ColorSuccess).
			Render(connectedMsg)
			
		refreshHelp := ui.HelpStyle.Render("Press 'r' to refresh • Esc: Back to list")
		
		content = lipgloss.JoinVertical(lipgloss.Left, 
			title, 
			statusBox,
			refreshHelp)
	} else if w.connecting {
		connectingMsg := fmt.Sprintf("Connecting to: %s", w.selectedSSID)
		statusBox := ui.BoxStyle.Copy().
			BorderForeground(ui.ColorPrimary).
			Render(connectingMsg)
			
		content = lipgloss.JoinVertical(lipgloss.Left, 
			title, 
			statusBox)
	} else if w.selectedSSID != "" && w.showPassword {
		// Password entry screen
		networkInfo := fmt.Sprintf("Enter password for: %s", ui.SelectedItemStyle.Render(w.selectedSSID))
		
		passwordPrompt := lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.NewStyle().Width(20).Render("Password:"),
			w.passwordInput.View())
			
		helpText := ui.HelpStyle.Render("Enter: Connect • Esc: Back to list")
		
		content = lipgloss.JoinVertical(lipgloss.Left,
			title,
			lipgloss.NewStyle().PaddingTop(1).Render(networkInfo),
			lipgloss.NewStyle().PaddingTop(1).Render(passwordPrompt),
			lipgloss.NewStyle().PaddingTop(1).Render(helpText))
	} else {
		if len(w.networks) == 0 {
			emptyMsg := "No WiFi networks found"
			emptyBox := ui.BoxStyle.Copy().
				Render(emptyMsg)
				
			refreshHelp := ui.HelpStyle.Render("Press 'r' to refresh")
			
			content = lipgloss.JoinVertical(lipgloss.Left, 
				title, 
				emptyBox,
				refreshHelp)
		} else {
			listView := w.list.View()
			
			helpText := ui.HelpStyle.Render("Enter: Select network • r: Refresh list")
			
			content = lipgloss.JoinVertical(lipgloss.Left,
				title,
				listView,
				helpText)
		}
	}
	
	return content
}
