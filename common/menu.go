
package common

import "github.com/rivo/tview"

// CreateMainMenuFunc is a type for the menu creation function
type CreateMainMenuFunc func(app *tview.Application) *tview.Flex

// MainMenuCreator holds the function to create the main menu
var MainMenuCreator CreateMainMenuFunc
