package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/user/networkmanager-tui/internal/ui/components"
	"github.com/user/networkmanager-tui/internal/ui/pages"
)

// keyMap defines the key mappings for the application
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Select key.Binding
	Back   key.Binding
	Quit   key.Binding
	Help   key.Binding
}

// Disable mouse events
func (m Model) MouseScrollLeft()  tea.Msg { return nil }
func (m Model) MouseScrollRight() tea.Msg { return nil }
func (m Model) MouseScrollUp()    tea.Msg { return nil }
func (m Model) MouseScrollDown()  tea.Msg { return nil }

// ShortHelp returns keybindings to be shown in the mini help view.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Select, k.Back, k.Quit, k.Help}
}

// FullHelp returns keybindings for the expanded help view.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Select, k.Back, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c", "q"),
		key.WithHelp("ctrl+c/q", "quit"),
	),
}

// Page represents a screen in the application
type Page int

const (
	MainMenuPage Page = iota
	EditConnectionPage
	AddConnectionPage
	WiFiListPage
	ActivateConnectionPage
)

// App is the main application model
type App struct {
	header  components.Header
	footer  components.Footer
	help    help.Model
	
	// Pages
	mainMenu           pages.MainMenu
	editConnection     pages.EditConnection
	addConnection      pages.AddConnection
	wifiList           pages.WiFiList
	activateConnection pages.ActivateConnection
	
	currentPage Page
	showHelp    bool
	width       int
	height      int
	ready       bool
}

// NewApp creates a new application model
func NewApp() *App {
	helpModel := help.New()
	helpModel.Width = 80

	return &App{
		header:  components.NewHeader(),
		footer:  components.NewFooter(),
		help:    helpModel,
		
		mainMenu:           pages.NewMainMenu(),
		editConnection:     pages.NewEditConnection(),
		addConnection:      pages.NewAddConnection(),
		wifiList:           pages.NewWiFiList(),
		activateConnection: pages.NewActivateConnection(),
		
		currentPage: MainMenuPage,
		showHelp:    false,
	}
}

// Init initializes the application
func (a App) Init() tea.Cmd {
	return tea.Batch(
		a.mainMenu.Init(),
	)
}

// Update handles user input and updates the application state
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return a, tea.Quit

		case key.Matches(msg, keys.Help):
			a.showHelp = !a.showHelp
			return a, nil
		}

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.ready = true

		headerHeight := a.header.Height()
		footerHeight := a.footer.Height()
		helpHeight := 0
		if a.showHelp {
			helpHeight = 4 // Approximate height for help
		}

		// Calculate the available height for the content
		contentHeight := a.height - headerHeight - footerHeight - helpHeight

		// Update header and footer with new dimensions
		a.header.SetWidth(msg.Width)
		a.footer.SetWidth(msg.Width)
		a.help.Width = msg.Width

		// Update pages with new dimensions
		a.mainMenu.SetSize(msg.Width, contentHeight)
		a.editConnection.SetSize(msg.Width, contentHeight)
		a.addConnection.SetSize(msg.Width, contentHeight)
		a.wifiList.SetSize(msg.Width, contentHeight)
		a.activateConnection.SetSize(msg.Width, contentHeight)
	}

	// Handle page-specific updates
	switch a.currentPage {
	case MainMenuPage:
		newModel, newCmd := a.mainMenu.Update(msg)
		a.mainMenu = newModel
		cmds = append(cmds, newCmd)

		// Check if we need to navigate to another page
		if a.mainMenu.Selected() == 0 {
			a.currentPage = EditConnectionPage
			cmds = append(cmds, a.editConnection.Init())
		} else if a.mainMenu.Selected() == 1 {
			a.currentPage = AddConnectionPage
			cmds = append(cmds, a.addConnection.Init())
		} else if a.mainMenu.Selected() == 2 {
			a.currentPage = WiFiListPage
			cmds = append(cmds, a.wifiList.Init())
		} else if a.mainMenu.Selected() == 3 {
			a.currentPage = ActivateConnectionPage
			cmds = append(cmds, a.activateConnection.Init())
		}

	case EditConnectionPage:
		if key.Matches(msg, keys.Back) {
			a.currentPage = MainMenuPage
			a.mainMenu.ResetSelection()
			return a, nil
		}
		newModel, newCmd := a.editConnection.Update(msg)
		a.editConnection = newModel
		cmds = append(cmds, newCmd)

	case AddConnectionPage:
		if key.Matches(msg, keys.Back) {
			a.currentPage = MainMenuPage
			a.mainMenu.ResetSelection()
			return a, nil
		}
		newModel, newCmd := a.addConnection.Update(msg)
		a.addConnection = newModel
		cmds = append(cmds, newCmd)

	case WiFiListPage:
		if key.Matches(msg, keys.Back) {
			a.currentPage = MainMenuPage
			a.mainMenu.ResetSelection()
			return a, nil
		}
		newModel, newCmd := a.wifiList.Update(msg)
		a.wifiList = newModel
		cmds = append(cmds, newCmd)

	case ActivateConnectionPage:
		if key.Matches(msg, keys.Back) {
			a.currentPage = MainMenuPage
			a.mainMenu.ResetSelection()
			return a, nil
		}
		newModel, newCmd := a.activateConnection.Update(msg)
		a.activateConnection = newModel
		cmds = append(cmds, newCmd)
	}

	return a, tea.Batch(cmds...)
}

// View renders the application
func (a App) View() string {
	if !a.ready {
		return "Initializing..."
	}

	// Get the content based on the current page
	var content string
	switch a.currentPage {
	case MainMenuPage:
		content = a.mainMenu.View()
	case EditConnectionPage:
		content = a.editConnection.View()
	case AddConnectionPage:
		content = a.addConnection.View()
	case WiFiListPage:
		content = a.wifiList.View()
	case ActivateConnectionPage:
		content = a.activateConnection.View()
	}

	// Combine the header, content, help (if shown), and footer
	view := a.header.View()
	view += "\n" + content

	if a.showHelp {
		view += "\n" + a.help.View(keys)
	}

	view += "\n" + a.footer.View()

	// Center everything on the screen
	return lipgloss.Place(
		a.width,
		a.height,
		lipgloss.Center,
		lipgloss.Center,
		view,
	)
}
