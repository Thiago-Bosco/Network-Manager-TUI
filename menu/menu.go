package menu

import (
        "flag"
        "fmt"
        "os/exec"

        "github.com/gdamore/tcell/v2"
        "github.com/rivo/tview"

        "networkmanager-tui/i18n"
        "networkmanager-tui/network"
        "networkmanager-tui/sysinfo"
)

// Cores padrão para UI - Paleta melhorada
var (
        borderColor      = tcell.ColorDeepSkyBlue    // Cor da borda
        backgroundColor  = tcell.ColorBlack          // Cor de fundo
        primaryTextColor = tcell.ColorWhite          // Cor do texto principal
        secondaryColor   = tcell.ColorAqua           // Cor secundária
        accentColor      = tcell.ColorTurquoise      // Cor de destaque
        titleColor       = tcell.ColorTurquoise      // Cor do título
        buttonBgColor    = tcell.ColorTurquoise      // Cor de fundo dos botões
        buttonTextColor  = tcell.ColorBlack          // Cor do texto nos botões
        successColor     = tcell.ColorPaleGreen      // Cor para mensagens de sucesso
        errorColor       = tcell.ColorSalmon         // Cor para mensagens de erro
        fieldBgColor     = tcell.ColorMidnightBlue   // Cor de fundo dos campos de entrada
        headerColor      = tcell.ColorDodgerBlue     // Cor dos cabeçalhos
        infoColor        = tcell.ColorLightSkyBlue   // Cor para mensagens informativas
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
                AddItem("📡 "+i18n.T("menu_status"), "Show network connections status", '2', func() {
                        showNetworkStatus(app)
                }).
                AddItem("📶 "+i18n.T("menu_ping_test"), "Test network connectivity", '3', func() {
                        showPingTest(app)
                }).
                AddItem("📊 "+i18n.T("menu_sysinfo"), "Display system information", '4', func() {
                        showSystemInfo(app)
                }).
                AddItem("ℹ️ "+i18n.T("menu_help"), "Show help information", '5', func() {
                        showHelp(app)
                }).
                AddItem("🔄 "+i18n.T("menu_reboot"), "Reboot the system", '6', func() {
                        confirmAndExecute(app, i18n.T("reboot_title"), i18n.T("reboot_message"), rebootSystem)
                }).
                AddItem("⏻ "+i18n.T("menu_shutdown"), "Shutdown the system", '7', func() {
                        confirmAndExecute(app, i18n.T("shutdown_title"), i18n.T("shutdown_message"), shutdownSystem)
                }).
                AddItem("🌐 "+i18n.T("menu_language"), "Change language", '8', func() {
                        changeLanguage(app)
                }).
                AddItem("❌ "+i18n.T("menu_exit"), "Exit the application", '9', func() {
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
        
        // Define a cor do título com base no tipo de mensagem
        var msgTitleColor tcell.Color
        if title == i18n.T("success_title") {
                msgTitleColor = successColor
        } else if title == i18n.T("error_title") {
                msgTitleColor = errorColor
        } else {
                msgTitleColor = titleColor
        }
        
        modal.SetBorder(true).
                SetTitle(" "+title+" ").
                SetTitleAlign(tview.AlignCenter).
                SetTitleColor(msgTitleColor).
                SetBorderColor(msgTitleColor).
                SetBackgroundColor(backgroundColor)
        
        app.SetRoot(modal, true)
}

// Mostra o status atual das conexões de rede
func showNetworkStatus(app *tview.Application) {
        flex := network.ShowNetworkStatus(app)
        app.SetRoot(flex, true)
}

// Testa conectividade de rede (ping)
func showPingTest(app *tview.Application) {
        // Cria o formulário de teste de ping
        form := tview.NewForm()
        form.SetBorder(true).
                SetTitle(" 📶 "+i18n.T("ping_title")+" 📶 ").
                SetTitleAlign(tview.AlignCenter).
                SetTitleColor(titleColor).
                SetBorderColor(borderColor).
                SetBackgroundColor(backgroundColor).
                SetBorderPadding(2, 2, 3, 3)
                
        // Configurando cores dos campos do formulário
        form.SetFieldBackgroundColor(fieldBgColor)
        form.SetFieldTextColor(primaryTextColor)
        form.SetLabelColor(secondaryColor)
        form.SetButtonBackgroundColor(buttonBgColor)
        form.SetButtonTextColor(buttonTextColor)
        
        // Campos para o teste de ping
        form.AddInputField(i18n.T("ping_target"), "8.8.8.8", 30, nil, nil)
        form.AddInputField(i18n.T("ping_count"), "4", 10, nil, nil)
        
        // Área de resultados
        resultsTextView := tview.NewTextView().
                SetDynamicColors(true).
                SetRegions(true).
                SetWordWrap(true).
                SetTextAlign(tview.AlignLeft).
                SetText(i18n.T("ping_results") + ":\n\n")
                
        resultsTextView.SetBorder(true).
                SetTitle(" "+i18n.T("ping_results")+" ").
                SetTitleAlign(tview.AlignCenter).
                SetBorderColor(borderColor)
        
        // Botões
        form.AddButton(i18n.T("ping_start"), func() {
                // Obter os valores dos campos
                targetHost := form.GetFormItemByLabel(i18n.T("ping_target")).(*tview.InputField).GetText()
                countStr := form.GetFormItemByLabel(i18n.T("ping_count")).(*tview.InputField).GetText()
                
                // Validação básica
                if targetHost == "" {
                        resultsTextView.SetText(i18n.T("ping_results") + ":\n\n" + 
                                "[red]" + i18n.T("error_empty_fields") + "[white]")
                        return
                }
                
                // Se count for vazio, usar valor padrão
                pingArgs := []string{"-c", "4", targetHost}
                if countStr != "" {
                        pingArgs = []string{"-c", countStr, targetHost}
                }
                
                // Executar o ping
                resultsTextView.SetText(i18n.T("ping_results") + ":\n\n" + 
                        "[yellow]" + i18n.T("ping_start") + " ping " + targetHost + "...[white]\n")
                
                // Verificamos se estamos no modo de desenvolvimento
                devMode := false
                flag.Visit(func(f *flag.Flag) {
                        if f.Name == "dev" && f.Value.String() == "true" {
                                devMode = true
                        }
                })
                
                if devMode {
                        // Em modo de desenvolvimento, simulamos a saída do ping
                        simulatedOutput := fmt.Sprintf("PING %s (%s) 56(84) bytes of data.\n", targetHost, targetHost)
                        simulatedOutput += "64 bytes from 8.8.8.8: icmp_seq=1 ttl=128 time=15.6 ms\n"
                        simulatedOutput += "64 bytes from 8.8.8.8: icmp_seq=2 ttl=128 time=14.2 ms\n"
                        simulatedOutput += "64 bytes from 8.8.8.8: icmp_seq=3 ttl=128 time=16.8 ms\n"
                        simulatedOutput += "64 bytes from 8.8.8.8: icmp_seq=4 ttl=128 time=13.9 ms\n\n"
                        simulatedOutput += "--- 8.8.8.8 ping statistics ---\n"
                        simulatedOutput += "4 packets transmitted, 4 received, 0% packet loss, time 3005ms\n"
                        simulatedOutput += "rtt min/avg/max/mdev = 13.921/15.137/16.821/1.154 ms\n"
                        
                        resultsTextView.SetText(i18n.T("ping_results") + ":\n\n" + 
                                "[green]" + simulatedOutput + "[white]")
                        return
                }
                
                // Executa o comando ping
                cmd := exec.Command("ping", pingArgs...)
                output, err := cmd.CombinedOutput()
                
                if err != nil {
                        resultsTextView.SetText(i18n.T("ping_results") + ":\n\n" + 
                                "[red]" + string(output) + "\n" + err.Error() + "[white]")
                        return
                }
                
                resultsTextView.SetText(i18n.T("ping_results") + ":\n\n" + 
                        "[green]" + string(output) + "[white]")
        })
        
        form.AddButton(i18n.T("network_back"), func() {
                app.SetRoot(createMainMenu(app), true)
        })
        
        // Layout principal
        flex := tview.NewFlex().
                SetDirection(tview.FlexRow).
                AddItem(form, 10, 0, true).
                AddItem(resultsTextView, 0, 1, false)
        
        app.SetRoot(flex, true)
}

// Mostra a tela de ajuda
func showHelp(app *tview.Application) {
        // Cria a área de texto para exibir a ajuda
        textView := tview.NewTextView().
                SetText(i18n.T("help_description")).
                SetDynamicColors(true).
                SetRegions(true).
                SetWordWrap(true).
                SetTextAlign(tview.AlignLeft)
        
        textView.SetBorder(true).
                SetTitle(" ℹ️ "+i18n.T("help_title")+" ℹ️ ").
                SetTitleAlign(tview.AlignCenter).
                SetTitleColor(titleColor).
                SetBorderColor(borderColor)
        
        // Botão para voltar ao menu principal
        form := tview.NewForm()
        form.SetBackgroundColor(backgroundColor)
        form.AddButton(i18n.T("network_back"), func() {
                app.SetRoot(createMainMenu(app), true)
        })
        
        // Layout principal
        flex := tview.NewFlex().
                SetDirection(tview.FlexRow).
                AddItem(textView, 0, 1, true).
                AddItem(form, 3, 0, false)
        
        app.SetRoot(flex, true)
}