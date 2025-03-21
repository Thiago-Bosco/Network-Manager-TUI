package pages

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	
	"github.com/user/networkmanager-tui/internal/ui"
)

// MainMenuItem represents an item in the main menu
type MainMenuItem struct {
	title       string
	description string
}

// Title returns the title of the menu item
func (i MainMenuItem) Title() string { return i.title }

// Description returns the description of the menu item
func (i MainMenuItem) Description() string { return i.description }

// FilterValue returns the title of the menu item
func (i MainMenuItem) FilterValue() string { return i.title }

// MainMenu is the main menu of the application
type MainMenu struct {
	list          list.Model
	selectedIndex int
	hasSelected   bool
	width         int
	height        int
}

// NewMainMenu creates a new main menu
func NewMainMenu() MainMenu {
	items := []list.Item{
		MainMenuItem{
			title:       "Edit Connection",
			description: "Edit an existing network connection",
		},
		MainMenuItem{
			title:       "Add Connection",
			description: "Create a new network connection",
		},
		MainMenuItem{
			title:       "WiFi Networks",
			description: "View and connect to available WiFi networks",
		},
		MainMenuItem{
			title:       "Activate Connection",
			description: "Activate or deactivate network connections",
		},
		MainMenuItem{
			title:       "Quit",
			description: "Exit Network Manager TUI",
		},
	}

	// Setup the list
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Network Manager TUI"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = ui.TitleStyle

	return MainMenu{
		list:          l,
		selectedIndex: -1,
		hasSelected:   false,
		width:         80,
		height:        20,
	}
}

// Init initializes the main menu
func (m MainMenu) Init() tea.Cmd {
	return nil
}

// SetSize sets the size of the main menu
func (m *MainMenu) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width, height)
}

// Selected returns the selected index, or -1 if nothing is selected
func (m MainMenu) Selected() int {
	if !m.hasSelected {
		return -1
	}
	return m.selectedIndex
}

// ResetSelection resets the selection state
func (m *MainMenu) ResetSelection() {
	m.selectedIndex = -1
	m.hasSelected = false
}

// Update handles user input and updates the menu state
func (m MainMenu) Update(msg tea.Msg) (MainMenu, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "enter" {
			m.selectedIndex = m.list.Index()
			m.hasSelected = true
			
			// If Quit is selected, return quit command
			if m.selectedIndex == 4 {
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the main menu
func (m MainMenu) View() string {
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(m.list.View())
}
