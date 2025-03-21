package pages

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/user/networkmanager-tui/internal/network"
	"github.com/user/networkmanager-tui/internal/ui"
)

// ConnectionItem represents a connection in the list
type ConnectionItem struct {
	Name     string
	Type     string
	Active   bool
}

// Title returns the name of the connection
func (c ConnectionItem) Title() string { 
	if c.Active {
		return fmt.Sprintf("%s ✓", c.Name)
	}
	return c.Name 
}

// Description returns the type of the connection
func (c ConnectionItem) Description() string { 
	return fmt.Sprintf("Type: %s", c.Type)
}

// FilterValue returns the name of the connection for filtering
func (c ConnectionItem) FilterValue() string { return c.Name }

// ActivateConnection is the page for activating/deactivating network connections
type ActivateConnection struct {
	list           list.Model
	width          int
	height         int
	loading        bool
	error          string
	statusMessage  string
}

// NewActivateConnection creates a new activate connection page
func NewActivateConnection() ActivateConnection {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Network Connections"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = ui.TitleStyle
	
	return ActivateConnection{
		list:           l,
		width:          80,
		height:         20,
		loading:        false,
		error:          "",
		statusMessage:  "",
	}
}

// Init initializes the activate connection page
func (a ActivateConnection) Init() tea.Cmd {
	return func() tea.Msg {
		connections, err := network.GetConnections()
		if err != nil {
			return network.ErrorMsg{Err: err}
		}
		return network.ConnectionsMsg{Connections: connections}
	}
}

// SetSize sets the size of the activate connection page
func (a *ActivateConnection) SetSize(width, height int) {
	a.width = width
	a.height = height
	a.list.SetSize(width, height-2) // Adjust for status message
}

// Update handles user input and updates the page state
func (a ActivateConnection) Update(msg tea.Msg) (ActivateConnection, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case network.ConnectionsMsg:
		a.loading = false
		items := make([]list.Item, len(msg.Connections))
		for i, conn := range msg.Connections {
			items[i] = ConnectionItem{
				Name:  conn.Name,
				Type:  conn.Type,
				Active: conn.Active,
			}
		}
		a.list.SetItems(items)
		
	case network.ErrorMsg:
		a.loading = false
		a.error = msg.Err.Error()
		
	case network.ConnectionStatusChangedMsg:
		a.loading = false
		if msg.Active {
			a.statusMessage = fmt.Sprintf("Connection '%s' activated successfully", msg.Name)
		} else {
			a.statusMessage = fmt.Sprintf("Connection '%s' deactivated successfully", msg.Name)
		}
		
		// Refresh the connection list
		cmds = append(cmds, func() tea.Msg {
			connections, err := network.GetConnections()
			if err != nil {
				return network.ErrorMsg{Err: err}
			}
			return network.ConnectionsMsg{Connections: connections}
		})
		
	case tea.KeyMsg:
		if msg.String() == "enter" && !a.loading && a.list.SelectedItem() != nil {
			selectedConnection := a.list.SelectedItem().(ConnectionItem)
			a.loading = true
			
			if selectedConnection.Active {
				// Deactivate the connection
				cmds = append(cmds, func() tea.Msg {
					err := network.DeactivateConnection(selectedConnection.Name)
					if err != nil {
						return network.ErrorMsg{Err: err}
					}
					return network.ConnectionStatusChangedMsg{
						Name:   selectedConnection.Name,
						Active: false,
					}
				})
			} else {
				// Activate the connection
				cmds = append(cmds, func() tea.Msg {
					err := network.ActivateConnection(selectedConnection.Name)
					if err != nil {
						return network.ErrorMsg{Err: err}
					}
					return network.ConnectionStatusChangedMsg{
						Name:   selectedConnection.Name,
						Active: true,
					}
				})
			}
		} else if msg.String() == "r" {
			// Refresh the list
			a.loading = true
			a.statusMessage = ""
			a.error = ""
			cmds = append(cmds, func() tea.Msg {
				connections, err := network.GetConnections()
				if err != nil {
					return network.ErrorMsg{Err: err}
				}
				return network.ConnectionsMsg{Connections: connections}
			})
		}
	}
	
	// Update list model
	if !a.loading {
		var cmd tea.Cmd
		a.list, cmd = a.list.Update(msg)
		cmds = append(cmds, cmd)
	}
	
	return a, tea.Batch(cmds...)
}

// View renders the activate connection page
func (a ActivateConnection) View() string {
	title := ui.TitleStyle.Copy().Render("Activate Connection")
	
	var content string
	
	if a.loading {
		loadingMsg := "Processing..."
		loadingBox := ui.BoxStyle.Copy().
			BorderForeground(ui.ColorPrimary).
			Render(loadingMsg)
			
		content = lipgloss.JoinVertical(lipgloss.Left, 
			title, 
			loadingBox)
	} else if a.error != "" {
		errorBox := ui.BoxStyle.Copy().
			BorderForeground(ui.ColorError).
			Render("Error: " + a.error)
			
		refreshHelp := ui.HelpStyle.Render("Press 'r' to refresh")
		
		content = lipgloss.JoinVertical(lipgloss.Left, 
			title, 
			errorBox,
			refreshHelp)
	} else {
		helpText := "Enter: Toggle connection status • r: Refresh list"
		
		var statusText string
		if a.statusMessage != "" {
			statusText = ui.StatusActiveStyle.Render(a.statusMessage)
		}
		
		content = lipgloss.JoinVertical(lipgloss.Left,
			title,
			a.list.View(),
			statusText,
			ui.HelpStyle.Render(helpText))
	}
	
	return content
}
