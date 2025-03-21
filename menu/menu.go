package menu

import (
        "fmt"
        "os/exec"

        "github.com/gdamore/tcell/v2"
        "github.com/rivo/tview"

        "networkmanager-tui/i18n"
        "networkmanager-tui/network"
        "networkmanager-tui/sysinfo"
)

// Cores padrão para UI
var (
        borderColor      = tcell.ColorDodgerBlue      // Cor da borda
        backgroundColor  = tcell.ColorBlack          // Cor de fundo
        primaryTextColor = tcell.ColorWhite          // Cor do texto principal
        secondaryColor   = tcell.ColorYellow         // Cor secundária
        accentColor      = tcell.ColorGreen          // Cor de destaque
)

// StartMenu inicia o menu principal da aplicação
func StartMenu(app *tview.Application) {
        mainFlex := createMainMenu(app)
        app.SetRoot(mainFlex, true)
}

// Cria o menu principal
func createMainMenu(app *tview.Application) *tview.Flex {
        // Lista com as opções do menu
        list := tview.NewList().
                AddItem("🔌 "+i18n.T("menu_configure"), "Configure network interfaces", '1', func() {
                        configureNetworkMenu(app)
                }).
                AddItem("📊 "+i18n.T("menu_sysinfo"), "Display system information", '2', func() {
                        showSystemInfo(app)
                }).
                AddItem("🔄 "+i18n.T("menu_reboot"), "Reboot the system", '3', func() {
                        confirmAndExecute(app, i18n.T("reboot_title"), i18n.T("reboot_message"), rebootSystem)
                }).
                AddItem("⏻ "+i18n.T("menu_shutdown"), "Shutdown the system", '4', func() {
                        confirmAndExecute(app, i18n.T("shutdown_title"), i18n.T("shutdown_message"), shutdownSystem)
                }).
                AddItem("🌐 "+i18n.T("menu_language"), "Change language", '5', func() {
                        changeLanguage(app)
                }).
                AddItem("❌ "+i18n.T("menu_exit"), "Exit the application", '6', func() {
                        app.Stop()
                })

        // Estiliza a lista
        list.SetBorder(true).
                SetTitle(" 🖥️ "+i18n.T("menu_title")+" 🖥️ ").
                SetTitleAlign(tview.AlignCenter).
                SetBorderColor(borderColor).
                SetBackgroundColor(backgroundColor)
        
        // Configura as cores do texto
        list.SetMainTextColor(primaryTextColor)
        list.SetSecondaryTextColor(secondaryColor)
        list.SetSelectedTextColor(backgroundColor)
        list.SetSelectedBackgroundColor(accentColor)

        // Cria um layout flexível para centralizar o menu
        flex := tview.NewFlex().
                SetDirection(tview.FlexRow).
                AddItem(nil, 0, 1, false). // Espaço em branco superior
                AddItem(tview.NewFlex().
                        AddItem(nil, 0, 1, false). // Espaço em branco à esquerda
                        AddItem(list, 40, 1, true). // Lista centralizada com largura fixa
                        AddItem(nil, 0, 1, false), // Espaço em branco à direita
                        10, 1, true). // Altura do menu
                AddItem(nil, 0, 1, false) // Espaço em branco inferior

        return flex
}

// Abre a tela de configuração de rede
func configureNetworkMenu(app *tview.Application) {
        form := network.ConfigureNetwork(app)
        app.SetRoot(form, true)
}

// Mostra as informações do sistema
func showSystemInfo(app *tview.Application) {
        output := sysinfo.GetSystemInfo()
        
        textView := tview.NewTextView().
                SetText(output).
                SetDynamicColors(true).
                SetRegions(true).
                SetWordWrap(true).
                SetTextAlign(tview.AlignLeft)
        
        textView.SetBorder(true).
                SetTitle(" 📊 "+i18n.T("sysinfo_title")+" 📊 ").
                SetTitleAlign(tview.AlignCenter).
                SetBorderColor(borderColor)
        
        // Botão para voltar ao menu principal
        backButton := tview.NewButton("Back").
                SetSelectedFunc(func() {
                        app.SetRoot(createMainMenu(app), true)
                })
        
        backButton.SetBackgroundColor(tcell.ColorRoyalBlue)
        backButton.SetLabelColor(tcell.ColorBlack)
        
        // Layout para a tela de informações do sistema
        flex := tview.NewFlex().
                SetDirection(tview.FlexRow).
                AddItem(textView, 0, 1, true).
                AddItem(backButton, 1, 0, false)
        
        app.SetRoot(flex, true)
}

// Reinicia o sistema
func rebootSystem() error {
        return exec.Command("reboot").Run()
}

// Desliga o sistema
func shutdownSystem() error {
        return exec.Command("shutdown", "-h", "now").Run()
}

// Confirmação antes de executar uma ação
func confirmAndExecute(app *tview.Application, title, message string, action func() error) {
        modal := tview.NewModal().
                SetText(message).
                AddButtons([]string{"Yes", "No"}).
                SetDoneFunc(func(buttonIndex int, buttonLabel string) {
                        if buttonIndex == 0 { // "Yes"
                                if err := action(); err != nil {
                                        showMessage(app, i18n.T("error_title"), fmt.Sprintf("Error: %v", err))
                                }
                        } else {
                                app.SetRoot(createMainMenu(app), true)
                        }
                })
        
        modal.SetBorder(true).
                SetTitle(" "+title+" ").
                SetTitleAlign(tview.AlignCenter).
                SetBackgroundColor(backgroundColor)
        
        app.SetRoot(modal, true)
}

// Alteração de idioma
func changeLanguage(app *tview.Application) {
        modal := tview.NewModal().
                SetText("Select language / Selecionar idioma:").
                AddButtons([]string{"English", "Português"}).
                SetDoneFunc(func(buttonIndex int, buttonLabel string) {
                        if buttonLabel == "English" {
                                i18n.SetLanguage("en")
                        } else {
                                i18n.SetLanguage("pt")
                        }
                        app.SetRoot(createMainMenu(app), true)
                })
        
        modal.SetBorder(true).
                SetTitle(" 🌐 Language / Idioma 🌐 ").
                SetTitleAlign(tview.AlignCenter).
                SetBackgroundColor(backgroundColor)
        
        app.SetRoot(modal, true)
}

// Exibe uma mensagem
func showMessage(app *tview.Application, title, message string) {
        modal := tview.NewModal().
                SetText(message).
                AddButtons([]string{"OK"}).
                SetDoneFunc(func(buttonIndex int, buttonLabel string) {
                        app.SetRoot(createMainMenu(app), true)
                })
        
        modal.SetBorder(true).
                SetTitle(" "+title+" ").
                SetTitleAlign(tview.AlignCenter).
                SetBackgroundColor(backgroundColor)
        
        app.SetRoot(modal, true)
}