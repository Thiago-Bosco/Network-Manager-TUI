package menu

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rivo/tview"

	"datadike-tui/i18n"
	network "datadike-tui/internal"
	"datadike-tui/sysinfo"
)

// Cores padrão
const (
	colorPrimary   = "#00ff00"
	colorSecondary = "#ffffff"
	colorError     = "#ff0000"
)

// ======================================
// Função Principal - Inicia o menu
func Start(app *tview.Application) {
	if !isRootUser() {
		fmt.Println(i18n.T("error_root_required"))
		os.Exit(1)
	}

	app.SetRoot(createMainMenu(app), true).Run()
}

// ======================================
// Verifica se o usuário é root
func isRootUser() bool {
	return os.Geteuid() == 0
}

// ======================================
// Criação do Menu Principal
func createMainMenu(app *tview.Application) *tview.Flex {
	menu := tview.NewList().
		AddItem("Network", "", '1', func() { ConfigureNetwork(app) }).
		AddItem(i18n.T("menu_sysinfo"), "", '2', func() { showSystemInfo(app) }).
		AddItem(i18n.T("menu_reboot"), "", '3', func() { confirmAndExecute(app, i18n.T("reboot_title"), i18n.T("reboot_message"), rebootSystem) }).
		AddItem(i18n.T("menu_shutdown"), "", '4', func() { confirmAndExecute(app, i18n.T("shutdown_title"), i18n.T("shutdown_message"), shutdownSystem) }).
		AddItem(i18n.T("menu_exit"), "", '5', func() { app.Stop() }).
		AddItem("Change Language", "", '6', func() { changeLanguage(app) })

	menu.SetBorder(true).SetTitle(i18n.T("menu_title")).SetTitleAlign(tview.AlignLeft)

	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false). // Espaço flexível acima
		AddItem(tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(menu, 0, 1, true).
			AddItem(nil, 0, 1, false), 0, 1, true).
		AddItem(nil, 0, 1, false) // Espaço flexível abaixo
}

// ======================================
// Mostra as Informações do Sistema
func showSystemInfo(app *tview.Application) {
	output := sysinfo.GetSystemInfo()
	showMessageAndReturn(app, i18n.T("sysinfo_title"), output)
}

// ======================================
// Reinicia o sistema
func rebootSystem() error {
	return exec.Command("sudo", "reboot").Run()
}

// ======================================
// Desliga o sistema
func shutdownSystem() error {
	return exec.Command("sudo", "shutdown", "-h", "now").Run()
}

// ======================================
// Confirmação antes de executar uma ação
func confirmAndExecute(app *tview.Application, title, message string, action func() error) {
	confirmAction(app, title, message, func() {
		if err := action(); err != nil {
			showMessage(app, i18n.T("error_title"), fmt.Sprintf("Failed: %v", err))
		}
	})
}

// ======================================
// Alteração de idioma
func changeLanguage(app *tview.Application) {
	modal := tview.NewModal().
		SetText("Select language:").
		AddButtons([]string{"English", "Português"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "English" {
				i18n.SetLanguage("en")
			} else {
				i18n.SetLanguage("pt")
			}
			updateMenuTitle(app)
		})

	modal.SetBorder(true).SetTitle("Change Language").SetTitleAlign(tview.AlignCenter)
	app.SetRoot(modal, true)
}

// ======================================
// Atualiza o título do menu após troca de idioma
func updateMenuTitle(app *tview.Application) {
	app.SetRoot(createMainMenu(app), true)
}

// ======================================
// Exibe uma mensagem e retorna ao menu
func showMessageAndReturn(app *tview.Application, title, message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			updateMenuTitle(app)
		})

	modal.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignCenter)
	app.SetRoot(modal, true)
}

// ======================================
// Exibe uma mensagem simples
func showMessage(app *tview.Application, title, message string) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(tview.NewBox(), true)
		})

	modal.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignCenter)
	app.SetRoot(modal, true)
}

// ======================================
// Confirmação de ação com botões "Yes" ou "No"
func confirmAction(app *tview.Application, title, message string, action func()) {
	modal := tview.NewModal().
		SetText(message).
		AddButtons([]string{"Yes", "No"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
				action()
			}
			app.SetRoot(tview.NewBox(), true)
		})

	modal.SetBorder(true).SetTitle(title).SetTitleAlign(tview.AlignCenter)
	app.SetRoot(modal, true)
}




// ======================================
// Seleciona a ferramenta de rede
func ConfigureNetwork(app *tview.Application) {
	network.ConfigureNetwork(app)
}
