package network

import (
        "fmt"
        "github.com/gdamore/tcell/v2"
        "github.com/rivo/tview"
        "os/exec"
        "regexp"
        "strings"
        "networkmanager-tui/i18n"
        "flag"
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

// Valores padrão para configuração manual
const (
        DefaultIP       = "192.168.1.100"  // Endereço IP padrão
        DefaultNetmask  = "24"             // Máscara de rede padrão (corresponde a 255.255.255.0)
        DefaultGateway  = "192.168.1.1"    // Gateway padrão
        DefaultDNS1     = "8.8.8.8"        // Servidor DNS primário (Google)
        DefaultDNS2     = "8.8.4.4"        // Servidor DNS secundário (Google)
)

// Valida um endereço IPv4
func validateIPv4(ip string) bool {
        // Padrão para validar endereços IPv4 (xxx.xxx.xxx.xxx onde xxx é um número de 0 a 255)
        ipPattern := `^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
        match, _ := regexp.MatchString(ipPattern, ip)
        return match
}

// Valida uma máscara de rede
func validateNetmask(netmask string) bool {
        // A máscara pode ser um número (CIDR) de 1 a 32
        cidrPattern := `^([1-9]|[12][0-9]|3[0-2])$`
        if match, _ := regexp.MatchString(cidrPattern, netmask); match {
                return true
        }
        
        // Ou uma máscara completa no formato xxx.xxx.xxx.xxx
        ipPattern := `^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
        match, _ := regexp.MatchString(ipPattern, netmask)
        return match
}

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
                        // Encerra a aplicação após mostrar a mensagem
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
        
        // Adicionamos o dropdown com o callback correto
        form.AddDropDown(i18n.T("network_dhcp"), mode_options, ipv4Mode, func(option string, index int) {
                ipv4Mode = index
                
                // Precisamos adicionar primeiro os campos para depois obter referências e depois definir seus estados
                if ipInput != nil {
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
                }
        })

        // Campos para configuração manual de IPv4 com valores padrão
        form.AddInputField(i18n.T("network_ipv4"), DefaultIP, 20, nil, nil)
        form.AddInputField(i18n.T("network_netmask"), DefaultNetmask, 20, nil, nil)
        form.AddInputField(i18n.T("network_gateway"), DefaultGateway, 20, nil, nil)
        form.AddInputField(i18n.T("network_dns1"), DefaultDNS1, 20, nil, nil)
        form.AddInputField(i18n.T("network_dns2"), DefaultDNS2, 20, nil, nil)
        
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

        // Verificamos se estamos no modo de desenvolvimento
        devMode := false
        flag.Visit(func(f *flag.Flag) {
                if f.Name == "dev" && f.Value.String() == "true" {
                        devMode = true
                }
        })

        // Botões
        form.AddButton(i18n.T("network_save"), func() {
                // Aqui implementaria a lógica para salvar as configurações
                index, _ := form.GetFormItemByLabel(i18n.T("network_interface")).(*tview.DropDown).GetCurrentOption()
                interfaceName := interfaces[index]
                
                if ipv4Mode == 0 { // DHCP
                        if devMode {
                                // Em modo de desenvolvimento, simulamos o comando
                                cmdStr := fmt.Sprintf("Simulando: nmcli con mod %s ipv4.method auto\n", interfaceName)
                                fmt.Print(cmdStr)
                                // Imprimimos informações para debug
                                fmt.Println("[DEBUG] Configurando com DHCP no modo de desenvolvimento")
                                fmt.Printf("[DEBUG] Interface: %s\n", interfaceName)
                                showMessage(app, i18n.T("success_title"), i18n.T("success_message")+" (DEV MODE)")
                                return
                        }
                        
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
                                showMessage(app, i18n.T("error_title"), i18n.T("error_empty_fields"))
                                return
                        }
                        
                        // Validação adicional de formato
                        if !validateIPv4(ip) {
                                showMessage(app, i18n.T("error_title"), i18n.T("error_invalid_ip"))
                                return
                        }
                        
                        if !validateNetmask(netmask) {
                                showMessage(app, i18n.T("error_title"), i18n.T("error_invalid_netmask"))
                                return
                        }
                        
                        if !validateIPv4(gateway) {
                                showMessage(app, i18n.T("error_title"), i18n.T("error_invalid_gateway"))
                                return
                        }
                        
                        if !validateIPv4(dns1) {
                                showMessage(app, i18n.T("error_title"), i18n.T("error_invalid_dns1"))
                                return
                        }
                        
                        if dns2 != "" && !validateIPv4(dns2) {
                                showMessage(app, i18n.T("error_title"), i18n.T("error_invalid_dns2"))
                                return
                        }
                        
                        if devMode {
                                // Em modo de desenvolvimento, simulamos o comando
                                cmdStr := fmt.Sprintf("Simulando: nmcli con mod %s ipv4.method manual ipv4.addresses %s/%s ipv4.gateway %s",
                                        interfaceName, ip, netmask, gateway)
                                
                                if dns2 != "" {
                                        cmdStr += fmt.Sprintf(" ipv4.dns %s,%s", dns1, dns2)
                                } else {
                                        cmdStr += fmt.Sprintf(" ipv4.dns %s", dns1)
                                }
                                
                                fmt.Println(cmdStr)
                                
                                // Logs detalhados para debug
                                fmt.Println("[DEBUG] Configurando com IP Manual no modo de desenvolvimento")
                                fmt.Printf("[DEBUG] Interface: %s\n", interfaceName)
                                fmt.Printf("[DEBUG] IP: %s\n", ip)
                                fmt.Printf("[DEBUG] Netmask: %s\n", netmask)
                                fmt.Printf("[DEBUG] Gateway: %s\n", gateway)
                                fmt.Printf("[DEBUG] DNS1: %s\n", dns1)
                                fmt.Printf("[DEBUG] DNS2: %s\n", dns2)
                                
                                showMessage(app, i18n.T("success_title"), i18n.T("success_message")+" (DEV MODE)")
                                return
                        }
                        
                        // Constrói o comando para configuração manual
                        var cmd *exec.Cmd
                        if dns2 != "" {
                                cmd = exec.Command("nmcli", "con", "mod", interfaceName,
                                        "ipv4.method", "manual",
                                        "ipv4.addresses", fmt.Sprintf("%s/%s", ip, netmask),
                                        "ipv4.gateway", gateway,
                                        "ipv4.dns", fmt.Sprintf("%s,%s", dns1, dns2))
                        } else {
                                cmd = exec.Command("nmcli", "con", "mod", interfaceName,
                                        "ipv4.method", "manual",
                                        "ipv4.addresses", fmt.Sprintf("%s/%s", ip, netmask),
                                        "ipv4.gateway", gateway,
                                        "ipv4.dns", dns1)
                        }
                        
                        err := cmd.Run()
                        if err != nil {
                                showMessage(app, i18n.T("error_title"), err.Error())
                                return
                        }
                }
                
                // Se não estamos no modo de desenvolvimento, ativamos a conexão
                if !devMode {
                        // Ativa a conexão
                        activateCmd := exec.Command("nmcli", "con", "up", interfaceName)
                        err := activateCmd.Run()
                        if err != nil {
                                showMessage(app, i18n.T("error_title"), err.Error())
                                return
                        }
                        
                        showMessage(app, i18n.T("success_title"), i18n.T("success_message"))
                }
        })
        
        form.AddButton(i18n.T("network_cancel"), func() {
                // Encerra a aplicação - ela será reiniciada pelo workflow
                app.Stop()
        })

        return form
}