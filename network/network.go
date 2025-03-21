package network

import (
        "fmt"
        "github.com/gdamore/tcell/v2"
        "github.com/rivo/tview"
        "os/exec"
        "strings"
        "networkmanager-tui/i18n"
)

// Definição de cores utilizadas na interface
var (
        borderColor      = tcell.ColorDodgerBlue      // Cor da borda
        backgroundColor  = tcell.ColorBlack          // Cor de fundo da interface
        fieldTextColor   = tcell.ColorWhite          // Cor do texto nos campos de entrada
        labelColor       = tcell.ColorYellow         // Cor dos rótulos
        buttonTextColor  = tcell.ColorBlack          // Cor do texto nos botões
        buttonBgColor    = tcell.ColorRoyalBlue      // Cor de fundo dos botões
        fieldBgColor     = tcell.ColorDarkSlateGray  // Cor de fundo dos campos de entrada
)

// Constantes para as opções de configuração de IPv4 e IPv6
const (
        IPv6ModeAuto   = "Auto"   // IPv6 configurado automaticamente
        IPv6ModeManual = "Manual" // IPv6 configurado manualmente
        IPv4ModeAuto   = "Auto"   // IPv4 configurado automaticamente
        IPv4ModeManual = "Manual" // IPv4 configurado manualmente
)

// Função que obtém as conexões de rede disponíveis usando o comando `nmcli`
func GetNetworkConnections() ([]string, error) {
        cmd := exec.Command("nmcli", "device", "status")
        out, err := cmd.Output()
        if err != nil {
                return nil, fmt.Errorf("erro ao obter conexões de rede: %w", err)
        }

        // Processa a saída do comando, extraindo os nomes das interfaces de rede
        lines := strings.Split(string(out), "\n")
        var interfaces []string
        for i, line := range lines {
                if i == 0 { // Pula o cabeçalho
                        continue
                }
                if len(line) > 0 { 
                        fields := strings.Fields(line)
                        if len(fields) > 0 {
                                interfaces = append(interfaces, fields[0])
                        }
                }
        }
        if len(interfaces) == 0 {
                // Fallback para interfaces comuns se o nmcli não retornar nada
                return []string{"eth0", "wlan0"}, nil
        }
        return interfaces, nil
}

// Função para obter o nome da conexão ativa
func GetActiveConnection() (string, error) {
        cmd := exec.Command("nmcli", "-t", "-f", "NAME", "connection", "show", "--active")
        output, err := cmd.Output()
        if err != nil {
                return "", fmt.Errorf("erro ao obter a conexão ativa: %w", err)
        }

        // Retorna o nome da conexão ativa
        return strings.TrimSpace(string(output)), nil
}

// Exibe uma mensagem de erro/sucesso
func showMessage(app *tview.Application, title, message string) {
        modal := tview.NewModal().
                SetText(message).
                AddButtons([]string{"OK"}).
                SetDoneFunc(func(buttonIndex int, buttonLabel string) {
                        app.Stop()
                })

        modal.SetBorder(true).
                SetTitle(" " + title + " ").
                SetTitleAlign(tview.AlignCenter).
                SetBackgroundColor(backgroundColor)

        app.SetRoot(modal, true)
}

// Função que configura a rede a partir de uma interface TUI
func ConfigureNetwork(app *tview.Application) *tview.Form {
        // Cria o formulário de configuração de rede
        form := tview.NewForm()
        form.SetBorder(true).
                SetTitle(" 🛠️ " + i18n.T("network_title") + " 🛠️ ").
                SetTitleAlign(tview.AlignCenter).
                SetBorderColor(borderColor).
                SetBackgroundColor(backgroundColor).
                SetBorderPadding(1, 1, 2, 2)
                
        // Configurando cores dos campos do formulário
        form.SetFieldBackgroundColor(fieldBgColor)
        form.SetFieldTextColor(fieldTextColor)
        form.SetLabelColor(labelColor)
        form.SetButtonBackgroundColor(buttonBgColor)
        form.SetButtonTextColor(buttonTextColor)

        // Obtém as interfaces de rede disponíveis
        interfaces, err := GetNetworkConnections()
        if err != nil {
                // Fallback para interfaces comuns
                interfaces = []string{"eth0", "wlan0"}
        }

        // Adiciona a opção de selecionar a interface de rede ao formulário
        form.AddDropDown(i18n.T("network_interface"), interfaces, 0, nil)

        // Adiciona opção para DHCP ou configuração manual
        mode_options := []string{IPv4ModeAuto, IPv4ModeManual}
        ipv4Mode := 0 // Padrão: Auto
        
        // Campos para configuração manual de IPv4
        var ipInput, netmaskInput, gatewayInput, dns1Input, dns2Input *tview.InputField
        
        // Adiciona dropdown para modo de configuração (DHCP ou Manual)
        form.AddDropDown(i18n.T("network_dhcp"), mode_options, ipv4Mode, func(option string, index int) {
                ipv4Mode = index
                if ipv4Mode == 1 { // Manual
                        // Mostra campos de configuração manual
                        ipInput.SetDisabled(false)
                        netmaskInput.SetDisabled(false)
                        gatewayInput.SetDisabled(false)
                        dns1Input.SetDisabled(false)
                        dns2Input.SetDisabled(false)
                } else { // Auto (DHCP)
                        // Desabilita campos de configuração manual
                        ipInput.SetDisabled(true)
                        netmaskInput.SetDisabled(true)
                        gatewayInput.SetDisabled(true)
                        dns1Input.SetDisabled(true)
                        dns2Input.SetDisabled(true)
                }
        })

        // Campos para configuração manual de IPv4
        form.AddInputField(i18n.T("network_ipv4"), "", 20, nil, nil)
        form.AddInputField(i18n.T("network_netmask"), "", 20, nil, nil)
        form.AddInputField(i18n.T("network_gateway"), "", 20, nil, nil)
        form.AddInputField(i18n.T("network_dns1"), "", 20, nil, nil)
        form.AddInputField(i18n.T("network_dns2"), "", 20, nil, nil)
        
        // Obtém referências aos campos de entrada
        ipInput = form.GetFormItemByLabel(i18n.T("network_ipv4")).(*tview.InputField)
        netmaskInput = form.GetFormItemByLabel(i18n.T("network_netmask")).(*tview.InputField)
        gatewayInput = form.GetFormItemByLabel(i18n.T("network_gateway")).(*tview.InputField)
        dns1Input = form.GetFormItemByLabel(i18n.T("network_dns1")).(*tview.InputField)
        dns2Input = form.GetFormItemByLabel(i18n.T("network_dns2")).(*tview.InputField)

        // Por padrão, desabilita os campos de configuração manual
        ipInput.SetDisabled(true)
        netmaskInput.SetDisabled(true)
        gatewayInput.SetDisabled(true)
        dns1Input.SetDisabled(true)
        dns2Input.SetDisabled(true)

        // Botões
        form.AddButton(i18n.T("network_save"), func() {
                // Aqui implementaria a lógica para salvar as configurações
                index, _ := form.GetFormItemByLabel(i18n.T("network_interface")).(*tview.DropDown).GetCurrentOption()
                interfaceName := interfaces[index]
                
                if ipv4Mode == 0 { // DHCP
                        // Configura a interface para usar DHCP
                        cmd := exec.Command("nmcli", "con", "mod", interfaceName, "ipv4.method", "auto")
                        err := cmd.Run()
                        if err != nil {
                                showMessage(app, i18n.T("error_title"), err.Error())
                                return
                        }
                } else { // Manual
                        // Configura a interface com IP estático
                        ip := ipInput.GetText()
                        netmask := netmaskInput.GetText()
                        gateway := gatewayInput.GetText()
                        dns1 := dns1Input.GetText()
                        dns2 := dns2Input.GetText()
                        
                        // Verifica se os campos obrigatórios estão preenchidos
                        if ip == "" || netmask == "" || gateway == "" || dns1 == "" {
                                showMessage(app, i18n.T("error_title"), "Please fill all required fields")
                                return
                        }
                        
                        // Constrói o comando para configuração manual
                        cmd := exec.Command("nmcli", "con", "mod", interfaceName,
                                "ipv4.method", "manual",
                                "ipv4.addresses", fmt.Sprintf("%s/%s", ip, netmask),
                                "ipv4.gateway", gateway,
                                "ipv4.dns", dns1)
                        
                        if dns2 != "" {
                                cmd = exec.Command("nmcli", "con", "mod", interfaceName,
                                        "ipv4.method", "manual",
                                        "ipv4.addresses", fmt.Sprintf("%s/%s", ip, netmask),
                                        "ipv4.gateway", gateway,
                                        "ipv4.dns", fmt.Sprintf("%s,%s", dns1, dns2))
                        }
                        
                        err := cmd.Run()
                        if err != nil {
                                showMessage(app, i18n.T("error_title"), err.Error())
                                return
                        }
                }
                
                // Ativa a conexão
                activateCmd := exec.Command("nmcli", "con", "up", interfaceName)
                err := activateCmd.Run()
                if err != nil {
                        showMessage(app, i18n.T("error_title"), err.Error())
                        return
                }
                
                showMessage(app, i18n.T("success_title"), i18n.T("success_message"))
        })
        
        form.AddButton(i18n.T("network_cancel"), func() {
                app.Stop()
        })

        return form
}