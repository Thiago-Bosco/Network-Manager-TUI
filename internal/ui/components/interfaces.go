package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/user/networkmanager-tui/internal/network"
	"github.com/user/networkmanager-tui/internal/ui"
)

// InterfacesList displays network interfaces in a table
type InterfacesList struct {
	table     table.Model
	interfaces []network.Interface
	width     int
	height    int
}

// NewInterfacesList creates a new interfaces list component
func NewInterfacesList() InterfacesList {
	columns := []table.Column{
		{Title: "Name", Width: 10},
		{Title: "Type", Width: 10},
		{Title: "Status", Width: 10},
		{Title: "IP Address", Width: 15},
		{Title: "MAC Address", Width: 20},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(ui.ColorBorder).
		BorderBottom(true).
		Bold(true).
		Foreground(ui.ColorPrimary)
	s.Selected = s.Selected.
		Foreground(ui.ColorBg).
		Background(ui.ColorPrimary).
		Bold(true)
	t.SetStyles(s)

	return InterfacesList{
		table:      t,
		interfaces: []network.Interface{},
		width:      80,
		height:     10,
	}
}

// SetSize sets the size of the interfaces list
func (i *InterfacesList) SetSize(width, height int) {
	i.width = width
	i.height = height
	i.table.SetHeight(min(height, 20))
	i.table.SetWidth(min(width, 80))
}

// SetInterfaces updates the interfaces in the list
func (i *InterfacesList) SetInterfaces(interfaces []network.Interface) {
	i.interfaces = interfaces
	
	rows := []table.Row{}
	for _, iface := range interfaces {
		status := "Inactive"
		if iface.IsActive {
			status = "Active"
		}
		
		rows = append(rows, table.Row{
			iface.Name,
			iface.Type,
			status,
			iface.IPAddress,
			iface.MACAddress,
		})
	}
	
	i.table.SetRows(rows)
}

// Update handles input for the interfaces list
func (i *InterfacesList) Update(msg table.Msg) (InterfacesList, table.Cmd) {
	var cmd table.Cmd
	i.table, cmd = i.table.Update(msg)
	return *i, cmd
}

// SelectedInterface returns the selected interface, if any
func (i *InterfacesList) SelectedInterface() *network.Interface {
	selected := i.table.SelectedRow()
	if selected == nil {
		return nil
	}
	
	// Find the interface that matches the selected row
	for idx, iface := range i.interfaces {
		if iface.Name == selected[0] {
			return &i.interfaces[idx]
		}
	}
	
	return nil
}

// View renders the interfaces list
func (i InterfacesList) View() string {
	if len(i.interfaces) == 0 {
		return ui.BoxStyle.Copy().
			Width(i.width - 4).
			Align(lipgloss.Center).
			Render("No network interfaces found")
	}
	
	title := ui.TitleStyle.Copy().Render("Network Interfaces")
	tableView := i.table.View()
	
	return lipgloss.JoinVertical(lipgloss.Left, title, tableView)
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
