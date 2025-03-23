package network

import (
        "context"
        "fmt"
        "github.com/gdamore/tcell/v2"
        "github.com/rivo/tview"
        "os/exec"
        "regexp"
        "strings"
        "time"
        "networkmanager-tui/i18n"
        "flag"
)

// Definição de cores utilizadas na interface - Paleta melhorada

// Timeout para comandos (5 segundos)
const commandTimeout = 5 * time.Second

// Executa comando com timeout
func runCommandWithTimeout(name string, args ...string) ([]byte, error) {
        ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
        defer cancel()
        
        cmd := exec.CommandContext(ctx, name, args...)
        out, err := cmd.Output()
        if ctx.Err() == context.DeadlineExceeded {
                return nil, fmt.Errorf("timeout ao executar comando %s", name)
        }
        return out, err
}

var (
        borderColor      = tcell.ColorDeepSkyBlue     // Cor da borda
        backgroundColor  = tcell.ColorBlack           // Cor de fundo da interface
        fieldTextColor   = tcell.ColorWhite           // Cor do texto nos campos de entrada
        labelColor       = tcell.ColorAqua            // Cor dos rótulos
        buttonTextColor  = tcell.ColorBlack           // Cor do texto nos botões
        buttonBgColor    = tcell.ColorTurquoise       // Cor de fundo dos botões
        fieldBgColor     = tcell.ColorMidnightBlue    // Cor de fundo dos campos de entrada
        titleColor       = tcell.ColorTurquoise       // Cor do título
        headerColor      = tcell.ColorDodgerBlue      // Cor dos cabeçalhos
        successColor     = tcell.ColorPaleGreen       // Cor para mensagens de sucesso
        errorColor       = tcell.ColorSalmon          // Cor para mensagens de erro
        infoColor        = tcell.ColorLightSkyBlue    // Cor para mensagens informativas
)

// Constantes para as opções de configuração de IPv4 e IPv6
const (
        IPv6ModeAuto     = "Auto"     // IPv6 configurado automaticamente
        IPv6ModeManual   = "Manual"   // IPv6 configurado manualmente
        IPv6ModeDisabled = "Disabled" // IPv6 desabilitado
        IPv4ModeAuto     = "Auto"     // IPv4 configurado automaticamente
        IPv4ModeManual   = "Manual"   // IPv4 configurado manualmente
)

// Valores padrão para configuração manual IPv4
const (
        DefaultIP       = "192.168.1.100"  // Endereço IP padrão
        DefaultNetmask  = "24"             // Máscara de rede padrão (corresponde a 255.255.255.0)
        DefaultGateway  = "192.168.1.1"    // Gateway padrão
        DefaultDNS1     = "8.8.8.8"        // Servidor DNS primário (Google)
        DefaultDNS2     = "8.8.4.4"        // Servidor DNS secundário (Google)
)

// Valores padrão para configuração manual IPv6
const (
        DefaultIPv6       = "2001:db8::1"       // Endereço IPv6 padrão de exemplo
        DefaultIPv6Prefix = "64"                // Prefixo padrão para IPv6
        DefaultIPv6Gateway= "2001:db8::1"       // Gateway IPv6 padrão
        DefaultIPv6DNS1   = "2001:4860:4860::8888" // Servidor DNS IPv6 primário (Google)
        DefaultIPv6DNS2   = "2001:4860:4860::8844" // Servidor DNS IPv6 secundário (Google)
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

// Valida um endereço IPv6
func validateIPv6(ip string) bool {
        // Formato simplificado para validação de IPv6
        // Esta expressão regular verifica o formato básico de um endereço IPv6
        ipv6Pattern := `^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$|^([0-9a-fA-F]{1,4}:){1,7}:|^([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}$|^([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}$|^([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}$|^([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}$|^([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}$|^[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})$|^:((:[0-9a-fA-F]{1,4}){1,7}|:)$`
        match, _ := regexp.MatchString(ipv6Pattern, ip)
        return match
}

// Valida um prefixo IPv6 (um número entre 1 e 128)
func validateIPv6Prefix(prefix string) bool {
        // O prefixo pode ser um número de 1 a 128
        prefixPattern := `^([1-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8])$`
        match, _ := regexp.MatchString(prefixPattern, prefix)
        return match
}

// Função que obtém as conexões de rede disponíveis usando o comando `nmcli`
func GetNetworkConnections() ([]string, error) {
        cmd := exec.Command("nmcli", "device", "status")
        out, err := cmd.Output()
        if err != nil {
                if exitErr, ok := err.(*exec.ExitError); ok {
                        return nil, fmt.Errorf("erro ao executar nmcli (código %d): %s", exitErr.ExitCode(), string(exitErr.Stderr))
                }
                return nil, fmt.Errorf("erro ao obter conexões de rede: %w", err)
        }
        if len(out) == 0 {
                return nil, fmt.Errorf("nenhuma saída do comando nmcli")
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

// GetActiveConnection retorna o nome da conexão ativa
func GetActiveConnection() (string, error) {
        cmd := exec.Command("nmcli", "-t", "-f", "NAME", "connection", "show", "--active")
        output, err := cmd.Output()
        if err != nil {
                return "", fmt.Errorf("erro ao obter a conexão ativa: %w", err)
        }

        return strings.TrimSpace(string(output)), nil
}

// Estrutura para armazenar informações detalhadas de uma conexão de rede
type NetworkConnectionInfo struct {
        Name      string // Nome da conexão
        Type      string // Tipo de conexão (wifi, ethernet, etc)
        Device    string // Dispositivo associado
        State     string // Estado da conexão (conectado, desconectado, etc)
        IPv4      string // Endereço IPv4
        IPv6      string // Endereço IPv6
        MAC       string // Endereço MAC
        Gateway   string // Gateway padrão
        DNS       string // Servidores DNS
}

// Obtém informações detalhadas das conexões de rede ativas
func GetNetworkConnectionsInfo() ([]NetworkConnectionInfo, error) {
        // Obtém o estado dos dispositivos
        cmd := exec.Command("nmcli", "-t", "device", "status")
        devOutput, err := cmd.Output()
        if err != nil {
                return nil, fmt.Errorf("erro ao obter status dos dispositivos: %w", err)
        }

        // Processa a saída
        lines := strings.Split(string(devOutput), "\n")
        var connections []NetworkConnectionInfo
        
        for _, line := range lines {
                if line == "" {
                        continue
                }
                
                fields := strings.Split(line, ":")
                if len(fields) >= 4 {
                        connInfo := NetworkConnectionInfo{
                                Device: fields[0],
                                Type:   fields[1],
                                State:  fields[2],
                                Name:   fields[3],
                        }
                        
                        // Se o dispositivo estiver conectado, obtém mais informações
                        if connInfo.State == "connected" {
                                // Obtém informações de IP para esta conexão
                                ipCmd := exec.Command("nmcli", "-t", "device", "show", connInfo.Device)
                                ipOutput, err := ipCmd.Output()
                                if err == nil {
                                        ipLines := strings.Split(string(ipOutput), "\n")
                                        for _, ipLine := range ipLines {
                                                if ipLine == "" {
                                                        continue
                                                }
                                                
                                                ipFields := strings.Split(ipLine, ":")
                                                if len(ipFields) >= 2 {
                                                        key := ipFields[0]
                                                        value := ipFields[1]
                                                        
                                                        switch {
                                                        case strings.Contains(key, "IP4.ADDRESS"):
                                                                connInfo.IPv4 = value
                                                        case strings.Contains(key, "IP6.ADDRESS"):
                                                                connInfo.IPv6 = value
                                                        case strings.Contains(key, "GENERAL.HWADDR"):
                                                                connInfo.MAC = value
                                                        case strings.Contains(key, "IP4.GATEWAY"):
                                                                connInfo.Gateway = value
                                                        case strings.Contains(key, "IP4.DNS"):
                                                                if connInfo.DNS == "" {
                                                                        connInfo.DNS = value
                                                                } else {
                                                                        connInfo.DNS += ", " + value
                                                                }
                                                        }
                                                }
                                        }
                                }
                        }
                        
                        connections = append(connections, connInfo)
                }
        }
        
        return connections, nil
}

// Exibe uma mensagem de erro/sucesso com cores apropriadas
func showMessage(app *tview.Application, title, message string) {
        modal := tview.NewModal().
                SetText(message).
                AddButtons([]string{"OK"}).
                SetDoneFunc(func(buttonIndex int, buttonLabel string) {
                        // Encerra a aplicação após mostrar a mensagem
                        app.Stop()
                })

        // Define cores com base no tipo de mensagem
        var titleColor tcell.Color
        if title == i18n.T("success_title") {
                titleColor = successColor
        } else if title == i18n.T("error_title") {
                titleColor = errorColor
        } else {
                titleColor = infoColor
        }

        modal.SetBorder(true).
                SetTitle(" " + title + " ").
                SetTitleAlign(tview.AlignCenter).
                SetTitleColor(titleColor).
                SetBorderColor(titleColor).
                SetBackgroundColor(backgroundColor)

        app.SetRoot(modal, true)
}

// Função que exibe o status atual das conexões de rede
func ShowNetworkStatus(app *tview.Application) *tview.Flex {
        // Obtém informações detalhadas das conexões
        connections, err := GetNetworkConnectionsInfo()
        
        // Flex container principal
        flex := tview.NewFlex().SetDirection(tview.FlexRow)
        
        // Cria uma tabela para exibir as informações
        table := tview.NewTable()
        table.SetBorders(true)
        table.SetBorderColor(borderColor)
        table.SetTitle(" 📊 " + i18n.T("network_status") + " 📊 ")
        table.SetTitleColor(titleColor)
        table.SetTitleAlign(tview.AlignCenter)
        table.Select(0, 0)
        table.SetSelectable(true, false)
        table.SetFixed(1, 0)
                
        table.SetBackgroundColor(backgroundColor)
                
        // Define títulos das colunas
        headers := []string{
                i18n.T("network_device"),    // Dispositivo
                i18n.T("network_type"),      // Tipo
                i18n.T("network_state"),     // Status
                i18n.T("network_name"),      // Nome
                i18n.T("network_ipv4"),      // IPv4
                i18n.T("network_ipv6"),      // IPv6
                i18n.T("network_gateway"),   // Gateway
                i18n.T("network_dns"),       // DNS
        }
        
        // Adiciona cabeçalho
        for col, header := range headers {
                cell := tview.NewTableCell(header).
                        SetTextColor(headerColor).
                        SetAlign(tview.AlignCenter).
                        SetSelectable(false)
                table.SetCell(0, col, cell)
        }
        
        // Se não conseguiu obter os dados, mostra mensagem de erro
        if err != nil {
                // Mostra mensagem na tabela
                errorCell := tview.NewTableCell(i18n.T("error_network_info") + ": " + err.Error()).
                        SetTextColor(errorColor).
                        SetAlign(tview.AlignCenter).
                        SetSelectable(false).
                        SetExpansion(1)
                table.SetCell(1, 0, errorCell)
                table.SetCell(1, 1, tview.NewTableCell("").SetSelectable(false))
                table.SetCell(1, 2, tview.NewTableCell("").SetSelectable(false))
                table.SetCell(1, 3, tview.NewTableCell("").SetSelectable(false))
                table.SetCell(1, 4, tview.NewTableCell("").SetSelectable(false))
                table.SetCell(1, 5, tview.NewTableCell("").SetSelectable(false))
                table.SetCell(1, 6, tview.NewTableCell("").SetSelectable(false))
                table.SetCell(1, 7, tview.NewTableCell("").SetSelectable(false))
                
                // Em desenvolvimento simulamos dados para testes
                devMode := false
                flag.Visit(func(f *flag.Flag) {
                        if f.Name == "dev" && f.Value.String() == "true" {
                                devMode = true
                        }
                })
                
                if devMode {
                        // Gera dados simulados para testes
                        connections = []NetworkConnectionInfo{
                                {
                                        Device:  "eth0",
                                        Type:    "ethernet",
                                        State:   "connected",
                                        Name:    "Ethernet Connection",
                                        IPv4:    "192.168.1.100/24",
                                        IPv6:    "fe80::1234:5678:abcd:ef12/64",
                                        Gateway: "192.168.1.1",
                                        DNS:     "8.8.8.8, 8.8.4.4",
                                },
                                {
                                        Device:  "wlan0",
                                        Type:    "wifi",
                                        State:   "disconnected",
                                        Name:    "Wi-Fi Network",
                                },
                                {
                                        Device:  "tun0",
                                        Type:    "tun",
                                        State:   "connected",
                                        Name:    "VPN Connection",
                                        IPv4:    "10.8.0.2/24",
                                        Gateway: "10.8.0.1",
                                        DNS:     "10.8.0.1",
                                },
                        }
                } else {
                        // Em caso de erro e não estando em modo dev, mantem a tela com a mensagem de erro
                        flex.AddItem(table, 0, 1, true)
                        
                        // Adiciona botões de ação
                        buttonsForm := tview.NewForm()
                        buttonsForm.SetBackgroundColor(backgroundColor)
                        
                        buttonsForm.AddButton(i18n.T("network_back"), func() {
                                app.Stop() // Retorna ao menu principal
                        })
                        
                        buttonsForm.AddButton(i18n.T("network_refresh"), func() {
                                // Recria a tela com dados atualizados
                                app.SetRoot(ShowNetworkStatus(app), true)
                        })
                        
                        flex.AddItem(buttonsForm, 3, 0, false)
                        
                        return flex
                }
        }
        
        // Preenche a tabela com os dados obtidos
        for row, conn := range connections {
                // Define a cor baseada no estado da conexão
                var stateColor tcell.Color
                switch conn.State {
                case "connected":
                        stateColor = successColor
                case "disconnected", "unavailable":
                        stateColor = errorColor
                default:
                        stateColor = fieldTextColor
                }
                
                // Dados a serem exibidos
                rowData := []struct {
                        text  string
                        color tcell.Color
                }{
                        {conn.Device, fieldTextColor},
                        {conn.Type, fieldTextColor},
                        {conn.State, stateColor},
                        {conn.Name, fieldTextColor},
                        {conn.IPv4, fieldTextColor},
                        {conn.IPv6, fieldTextColor},
                        {conn.Gateway, fieldTextColor},
                        {conn.DNS, fieldTextColor},
                }
                
                // Adiciona os dados à tabela
                for col, data := range rowData {
                        cell := tview.NewTableCell(data.text).
                                SetTextColor(data.color).
                                SetAlign(tview.AlignLeft)
                        table.SetCell(row+1, col, cell)
                }
        }
        
        flex.AddItem(table, 0, 1, true)
        
        // Adiciona botões de ação
        buttonsForm := tview.NewForm()
        buttonsForm.SetBackgroundColor(backgroundColor)
        
        buttonsForm.AddButton(i18n.T("network_back"), func() {
                app.Stop() // Retorna ao menu principal
        })
        
        buttonsForm.AddButton(i18n.T("network_refresh"), func() {
                // Recria a tela com dados atualizados
                app.SetRoot(ShowNetworkStatus(app), true)
        })
        
        // Adicionando texto de ajuda para mostrar a tecla Esc
        helpText := tview.NewTextView()
        helpText.SetTextAlign(tview.AlignCenter)
        helpText.SetDynamicColors(true)
        helpText.SetText("[yellow]" + i18n.T("press_esc_return") + "[white]")
        
        flex.AddItem(buttonsForm, 3, 0, false)
        flex.AddItem(helpText, 1, 0, false)
        
        return flex
}

// Função que configura a rede a partir de uma interface TUI
func ConfigureNetwork(app *tview.Application) *tview.Flex {
        // Cria o formulário de configuração de rede com um visual mais bonito
        form := tview.NewForm()
        form.SetBorder(true).
                SetTitle(" 🛠️  " + i18n.T("network_title") + " 🛠️  ").
                SetTitleAlign(tview.AlignCenter).
                SetTitleColor(titleColor).
                SetBorderColor(borderColor).
                SetBackgroundColor(backgroundColor).
                SetBorderPadding(2, 2, 3, 3)
                
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

        // === Configuração IPv4 ===
        form.AddTextView("", "=== "+i18n.T("network_ipv4_config")+" ===", 20, 1, true, false)

        // Adiciona opção para DHCP ou configuração manual IPv4
        mode_options := []string{IPv4ModeAuto, IPv4ModeManual}
        ipv4Mode := 0 // Padrão: Auto
        
        // Campos para configuração manual de IPv4
        var ipInput, netmaskInput, gatewayInput, dns1Input, dns2Input *tview.InputField
        
        // Adicionamos o dropdown com o callback correto
        form.AddDropDown(i18n.T("network_ipv4_mode"), mode_options, ipv4Mode, func(option string, index int) {
                if ipInput != nil {
                    if option == IPv4ModeManual {
                        // Mostra campos de configuração manual IPv4
                        ipInput.SetDisabled(false)
                        netmaskInput.SetDisabled(false)
                        gatewayInput.SetDisabled(false)
                        dns1Input.SetDisabled(false)
                        dns2Input.SetDisabled(false)
                    } else { // Auto/DHCP
                        // Desabilita campos de configuração manual IPv4
                        ipInput.SetDisabled(true)
                        netmaskInput.SetDisabled(true)
                        gatewayInput.SetDisabled(true)
                        dns1Input.SetDisabled(true)
                        dns2Input.SetDisabled(true)
                    }
                }
        })

        // Campos para configuração manual de IPv4 com valores padrão
        form.AddInputField(i18n.T("network_ipv4_address"), DefaultIP, 20, nil, nil)
        form.AddInputField(i18n.T("network_ipv4_netmask"), DefaultNetmask, 20, nil, nil)
        form.AddInputField(i18n.T("network_ipv4_gateway"), DefaultGateway, 20, nil, nil)
        form.AddInputField(i18n.T("network_ipv4_dns1"), DefaultDNS1, 20, nil, nil)
        form.AddInputField(i18n.T("network_ipv4_dns2"), DefaultDNS2, 20, nil, nil)
        
        // Obtém referências aos campos de entrada IPv4
        ipInput = form.GetFormItemByLabel(i18n.T("network_ipv4_address")).(*tview.InputField)
        netmaskInput = form.GetFormItemByLabel(i18n.T("network_ipv4_netmask")).(*tview.InputField)
        gatewayInput = form.GetFormItemByLabel(i18n.T("network_ipv4_gateway")).(*tview.InputField)
        dns1Input = form.GetFormItemByLabel(i18n.T("network_ipv4_dns1")).(*tview.InputField)
        dns2Input = form.GetFormItemByLabel(i18n.T("network_ipv4_dns2")).(*tview.InputField)

        // Por padrão, desabilita os campos de configuração manual IPv4
        ipInput.SetDisabled(true)
        netmaskInput.SetDisabled(true)
        gatewayInput.SetDisabled(true)
        dns1Input.SetDisabled(true)
        dns2Input.SetDisabled(true)
        
        // === Configuração IPv6 ===
        form.AddTextView("", "=== "+i18n.T("network_ipv6_config")+" ===", 20, 1, true, false)
        
        // Adiciona opções para IPv6 (Auto, Manual, Desabilitado)
        ipv6_options := []string{IPv6ModeAuto, IPv6ModeManual, IPv6ModeDisabled}
        ipv6Mode := 0 // Padrão: Auto
        
        // Campos para configuração manual de IPv6
        var ipv6Input, ipv6PrefixInput, ipv6GatewayInput, ipv6DNS1Input, ipv6DNS2Input *tview.InputField
        
        // Dropdown para modo IPv6 com callback
        form.AddDropDown(i18n.T("network_ipv6_mode"), ipv6_options, ipv6Mode, func(option string, index int) {
                ipv6Mode = index
                
                if ipv6Input != nil {
                    if ipv6Mode == 1 { // Manual
                        // Habilita campos de configuração manual IPv6
                        ipv6Input.SetDisabled(false)
                        ipv6PrefixInput.SetDisabled(false)
                        ipv6GatewayInput.SetDisabled(false)
                        ipv6DNS1Input.SetDisabled(false)
                        ipv6DNS2Input.SetDisabled(false)
                    } else { // Auto ou Desabilitado
                        // Desabilita campos de configuração manual IPv6
                        ipv6Input.SetDisabled(true)
                        ipv6PrefixInput.SetDisabled(true)
                        ipv6GatewayInput.SetDisabled(true)
                        ipv6DNS1Input.SetDisabled(true)
                        ipv6DNS2Input.SetDisabled(true)
                    }
                }
        })
        
        // Campos para configuração manual de IPv6 com valores padrão
        form.AddInputField(i18n.T("network_ipv6_address"), DefaultIPv6, 40, nil, nil)
        form.AddInputField(i18n.T("network_ipv6_prefix"), DefaultIPv6Prefix, 20, nil, nil)
        form.AddInputField(i18n.T("network_ipv6_gateway"), DefaultIPv6Gateway, 40, nil, nil)
        form.AddInputField(i18n.T("network_ipv6_dns1"), DefaultIPv6DNS1, 40, nil, nil)
        form.AddInputField(i18n.T("network_ipv6_dns2"), DefaultIPv6DNS2, 40, nil, nil)
        
        // Obtém referências aos campos de entrada IPv6
        ipv6Input = form.GetFormItemByLabel(i18n.T("network_ipv6_address")).(*tview.InputField)
        ipv6PrefixInput = form.GetFormItemByLabel(i18n.T("network_ipv6_prefix")).(*tview.InputField)
        ipv6GatewayInput = form.GetFormItemByLabel(i18n.T("network_ipv6_gateway")).(*tview.InputField)
        ipv6DNS1Input = form.GetFormItemByLabel(i18n.T("network_ipv6_dns1")).(*tview.InputField)
        ipv6DNS2Input = form.GetFormItemByLabel(i18n.T("network_ipv6_dns2")).(*tview.InputField)
        
        // Por padrão, desabilita os campos de configuração manual IPv6 (exceto no modo manual)
        ipv6Input.SetDisabled(true)
        ipv6PrefixInput.SetDisabled(true)
        ipv6GatewayInput.SetDisabled(true)
        ipv6DNS1Input.SetDisabled(true)
        ipv6DNS2Input.SetDisabled(true)

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
        
        // Adicionando texto de ajuda para mostrar a tecla Esc
        helpText := tview.NewTextView()
        helpText.SetTextAlign(tview.AlignCenter)
        helpText.SetDynamicColors(true)
        helpText.SetText("[yellow]" + i18n.T("press_esc_return") + "[white]")
        
        // Criando um flex para adicionar o texto de ajuda abaixo do formulário
        flex := tview.NewFlex().
                SetDirection(tview.FlexRow).
                AddItem(form, 0, 1, true).
                AddItem(helpText, 1, 0, false)
        
        return flex
}