package network

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"os/exec"
	"strings"
	"net"
	"datadike-tui/i18n"
)

// Definição de cores utilizadas na interface (configuração da UI)
var (
	borderColor      = tcell.ColorDodgerBlue  // Cor da borda
	backgroundColor  = tcell.ColorBlack      // Cor de fundo da interface
	fieldTextColor   = tcell.ColorWhite      // Cor do texto nos campos de entrada
	labelColor       = tcell.ColorYellow     // Cor dos rótulos
	buttonTextColor  = tcell.ColorBlack      // Cor do texto nos botões
	buttonBgColor    = tcell.ColorRoyalBlue  // Cor de fundo dos botões
	fieldBgColor     = tcell.ColorDarkSlateGray // Cor de fundo dos campos de entrada
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
	for _, line := range lines {
		if len(line) > 0 && line[0] != ' ' { // Ignora linhas que não contêm nome da interface
			fields := strings.Fields(line)
			if len(fields) > 0 {
				interfaces = append(interfaces, fields[0])
			}
		}
	}
	if len(interfaces) == 0 {
		return nil, fmt.Errorf("nenhuma interface de rede encontrada")
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

// Função que configura a rede a partir de uma interface TUI (Terminal User Interface)
func ConfigureNetwork(app *tview.Application) {
	// Cria o formulário de configuração de rede
	form := tview.NewForm()
	form.SetBorder(true).
		SetTitle(" 🛠️ " + i18n.T("network_title") + " 🛠️ ").
		SetTitleAlign(tview.AlignCenter).
		SetBorderColor(borderColor).
		SetBackgroundColor(backgroundColor).
		SetBorderPadding(1, 1, 2, 2)

	// Obtém as interfaces de rede disponíveis
	interfaces, err := GetNetworkConnections()
	if err != nil {
		// Exibe uma mensagem de erro se não encontrar interfaces de rede
		showMessage(app, i18n.T("error_title"), err.Error())
		return
	}

	// Adiciona a opção de selecionar a interface de rede ao formulário
	form.AddDropDown(i18n.T("network_interface"), interfaces, 0, nil).
		SetFieldTextColor(fieldTextColor).
		SetLabelColor(labelColor).
		SetFieldBackgroundColor(fieldBgColor)

	// Adiciona as configurações de rede ao formulário
	addNetworkSettings(form)

	// Botão de confirmação para aplicar as configurações
	form.AddButton(i18n.T("network_save"), func() {
		// Valida e aplica as configurações
		if err := applyNetworkSettings(form); err != nil {
			showMessage(app, i18n.T("error_title"), err.Error())
		} else {
			showMessage(app, i18n.T("success_title"), i18n.T("success_message"))

			//  lida com a mudança do nome do host
			newHostName := form.GetFormItemByLabel("Hostname").(*tview.InputField).GetText()
			if newHostName != "" {
				// Aplicando a mudança de hostname
				if err := setHostName(newHostName); err != nil {
					showMessage(app, "Error", "Falha ao definir o hostname: "+err.Error())
					return
				}

				// Verificando se a mudança foi bem-sucedida
				hostname, err := getHostName()
				if err != nil {
					showMessage(app, "Error", "Falha ao recuperar o hostname: "+err.Error())
					return
				}

				// Confirmando se o hostname foi alterado corretamente
				if hostname != newHostName {
					showMessage(app, "Error", "O hostname não foi atualizado corretamente.")
					return
				}
			}
		}
	}).SetButtonTextColor(buttonTextColor).SetButtonBackgroundColor(buttonBgColor)

	// Botão de voltar que encerra a aplicação
	form.AddButton(i18n.T("network_cancel"), func() {
		app.Stop()
	}).SetButtonTextColor(buttonTextColor).SetButtonBackgroundColor(buttonBgColor)

	// Configura o layout do formulário na tela
	form.SetRect(10, 5, 60, 20)
	form.SetFieldBackgroundColor(fieldBgColor)

	// Define o foco inicial para o primeiro campo
	form.SetFocus(0)
	// Exibe o formulário na aplicação
	app.SetRoot(form, true)
}

// Função que adiciona campos de configuração ao formulário
func addNetworkSettings(form *tview.Form) {
	// Adiciona a opção de configurar DNS
	form.AddDropDown(i18n.T("network_dhcp"), []string{"Yes", "No"}, 0, nil).
		SetFieldTextColor(fieldTextColor).
		SetLabelColor(labelColor).
		SetFieldBackgroundColor(fieldBgColor)

	// Adiciona o campo de seleção de modo IPv6 (Auto ou Manual)
	form.AddDropDown(("network_ipv6_mode"), []string{IPv6ModeAuto, IPv6ModeManual}, 0, func(option string, index int) {
		// Se a opção selecionada for "Manual", exibe o campo de entrada para o endereço IPv6
		if option == IPv6ModeManual {
			form.AddInputField("IPv6 Address", "", 30, nil, nil).
				SetFieldTextColor(fieldTextColor).
				SetLabelColor(labelColor).
				SetFieldBackgroundColor(fieldBgColor)
		}
	}).
		SetFieldTextColor(fieldTextColor).
		SetLabelColor(labelColor).
		SetFieldBackgroundColor(fieldBgColor)

	// Adiciona o campo de seleção de modo IPv4 (Auto ou Manual)
	form.AddDropDown(("network_ipv4_mode"), []string{IPv4ModeAuto, IPv4ModeManual}, 0, func(option string, index int) {
		// Se a opção selecionada for "Manual", exibe os campos de entrada para o endereço IP estático e Gateway
		if option == IPv4ModeManual {
			form.AddInputField("Static IP Address", "", 30, nil, nil).
				SetFieldTextColor(fieldTextColor).
				SetLabelColor(labelColor).
				SetFieldBackgroundColor(fieldBgColor)
			form.AddInputField("Gateway", "", 30, nil, nil).
				SetFieldTextColor(fieldTextColor).
				SetLabelColor(labelColor).
				SetFieldBackgroundColor(fieldBgColor)
		}
	}).
		SetFieldTextColor(fieldTextColor).
		SetLabelColor(labelColor).
		SetFieldBackgroundColor(fieldBgColor)

	// Adiciona campos para configuração de DNS
	form.AddInputField(i18n.T("network_dns1"), "", 30, nil, nil).
		SetFieldTextColor(fieldTextColor).
		SetLabelColor(labelColor).
		SetFieldBackgroundColor(fieldBgColor)

	form.AddInputField(i18n.T("network_dns2"), "", 30, nil, nil).
		SetFieldTextColor(fieldTextColor).
		SetLabelColor(labelColor).
		SetFieldBackgroundColor(fieldBgColor)

	// Adiciona o campo para configurar o nome do host
	form.AddInputField("Hostname", "", 30, nil, nil).
		SetFieldTextColor(fieldTextColor).
		SetLabelColor(labelColor).
		SetFieldBackgroundColor(fieldBgColor)
}

// Função para aplicar a configuração do nome do host
func setHostName(newHostName string) error {
	cmd := exec.Command("nmcli", "general", "hostname", newHostName)
	return cmd.Run()
}

func getHostName() (string, error) {
	cmd := exec.Command("nmcli", "general", "hostname")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// Função para aplicar as configurações de rede baseadas nas opções selecionadas
func applyNetworkSettings(form *tview.Form) error {
	// Obtém o nome da conexão ativa
	connectionName, err := GetActiveConnection()
	if err != nil {
		return fmt.Errorf("erro ao obter a conexão ativa: %w", err)
	}

	// Obtém o modo de configuração do IPv6 selecionado (Auto ou Manual)
	ipv6Mode, _ := form.GetFormItemByLabel(("network_ipv6_mode")).(*tview.DropDown).GetCurrentOption()

	// Se o modo IPv6 for manual, aplica a configuração de endereço IPv6
	if ipv6Mode == 1 { // Índice 1 é o modo "Manual"
		ipv6Address := form.GetFormItemByLabel("IPv6 Address").(*tview.InputField).GetText()
		if ipv6Address != "" {
			// Valida se o endereço IPv6 fornecido é válido
			if !isValidIP(ipv6Address) {
				return fmt.Errorf("endereço IPv6 inválido")
			}
			// Aplica a configuração IPv6 usando o comando `nmcli`
			err := configureIPv6(ipv6Address, connectionName)
			if err != nil {
				return fmt.Errorf("erro ao aplicar configuração IPv6 manual: %w", err)
			}
			fmt.Println(i18n.T("success_message"))
		}
	}

	// Obtém o modo de configuração do IPv4 selecionado (Auto ou Manual)
	ipv4Mode, _ := form.GetFormItemByLabel(("network_ipv4_mode")).(*tview.DropDown).GetCurrentOption()

	// Se o modo IPv4 for manual, aplica a configuração de endereço IPv4
	if ipv4Mode == 1 { // Índice 1 é o modo "Manual"
		ipv4Address := form.GetFormItemByLabel("Static IP Address").(*tview.InputField).GetText()
		if ipv4Address != "" {
			// Valida se o endereço IPv4 fornecido é válido
			if !isValidIP(ipv4Address) {
				return fmt.Errorf("endereço IPv4 inválido")
			}
			// Aplica a configuração IPv4 usando o comando `nmcli`
			err := configureIPv4(ipv4Address, connectionName)
			if err != nil {
				return fmt.Errorf("erro ao aplicar configuração IPv4 manual: %w", err)
			}
			fmt.Println(i18n.T("success_message"))
		}
	}

	// Aplique outras configurações (DNS, etc) conforme necessário

	return nil
}

// Função que valida se um endereço IP fornecido é válido
func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// Função que aplica a configuração do IPv6 manualmente usando o comando `nmcli`
func configureIPv6(ipv6Address, connectionName string) error {
	cmd := exec.Command("nmcli", "connection", "modify", connectionName, "ipv6.addresses", ipv6Address)
	return cmd.Run()
}

// Função que aplica a configuração do IPv4 manualmente usando o comando `nmcli`
func configureIPv4(ipv4Address, connectionName string) error {
	cmd := exec.Command("nmcli", "connection", "modify", connectionName, "ipv4.addresses", ipv4Address)
	return cmd.Run()
}

// Função que exibe uma mensagem de erro ou sucesso em uma caixa de diálogo
func showMessage(app *tview.Application, title, message string) {
	dialog := tview.NewModal().
		SetText(message).
		AddButtons([]string{i18n.T("menu_exit")}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.Stop()
		})

	dialog.SetTitle(title)
	app.SetRoot(dialog, true)
}
