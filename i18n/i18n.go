package i18n

import (
        "sync"
)

var (
        language = "en" // Idioma padrão: inglês
        mu       sync.Mutex
)

// Strings para cada idioma
var translations = map[string]map[string]string{
        "en": {
                "menu_title":        "Network Manager TUI",
                "menu_configure":    "Configure Network",
                "menu_status":       "Network Status",
                "menu_ping_test":    "Ping Test",
                "menu_sysinfo":      "System Information",
                "menu_reboot":       "Reboot",
                "menu_shutdown":     "ShutDown",
                "menu_exit":         "Exit",
                "menu_language":     "Change Language",
                "menu_help":         "Help",
                
                "network_title":     "Configure Network",
                "network_status":    "Network Status",
                "network_interface": "Network Interface (e.g., eth0):",
                "network_dhcp":      "DHCP:",
                "network_ipv4":      "IPv4 Address:",
                "network_ipv6":      "IPv6 Address:",
                "network_mac":       "MAC Address:",
                "network_netmask":   "Netmask:",
                "network_gateway":   "Gateway:",
                "network_dns":       "DNS:",
                "network_dns1":      "Primary DNS (e.g., 8.8.8.8):",
                "network_dns2":      "Secondary DNS (optional):",
                "network_save":      "Save",
                "network_cancel":    "Cancel",
                "network_back":      "Back",
                "network_refresh":   "Refresh",
                "network_device":    "Device",
                "network_type":      "Type",
                "network_state":     "Status",
                "network_name":      "Name",
                "network_ipv4_config": "IPv4 Configuration",
                "network_ipv6_config": "IPv6 Configuration",
                "network_ipv4_mode": "IPv4 Mode:",
                "network_ipv6_mode": "IPv6 Mode:",
                "network_ipv4_address": "IPv4 Address:",
                "network_ipv4_netmask": "Netmask:",
                "network_ipv4_gateway": "IPv4 Gateway:",
                "network_ipv4_dns1": "IPv4 Primary DNS:",
                "network_ipv4_dns2": "IPv4 Secondary DNS:",
                "network_ipv6_address": "IPv6 Address:",
                "network_ipv6_prefix": "IPv6 Prefix Length:",
                "network_ipv6_gateway": "IPv6 Gateway:",
                "network_ipv6_dns1": "IPv6 Primary DNS:",
                "network_ipv6_dns2": "IPv6 Secondary DNS:",
                
                "ping_title":        "Ping Test",
                "ping_target":       "Target Host (e.g., 8.8.8.8):",
                "ping_count":        "Count (optional):",
                "ping_start":        "Start",
                "ping_results":      "Ping Results",
                
                "sysinfo_title":     "System Information",
                "reboot_title":      "Reboot",
                "reboot_message":    "The system will now reboot!",
                "shutdown_title":    "ShutDown",
                "shutdown_message":  "The system will now shut down!",
                
                "help_title":        "Help",
                "help_description":  "This application allows you to configure and monitor network interfaces.\n\nUse arrow keys to navigate, Enter to select, and Esc to go back.",
                
                "success_title":     "Success",
                "success_message":   "Network configuration applied successfully!",
                "error_title":       "Error",
                "error_root_required": "This application requires root privileges.",
                "error_empty_fields": "Please fill all required fields",
                "error_invalid_ip":    "Invalid IP address format",
                "error_invalid_netmask": "Invalid netmask format (use CIDR number or full mask)",
                "error_invalid_gateway": "Invalid gateway IP address format",
                "error_invalid_dns1": "Invalid primary DNS format",
                "error_invalid_dns2": "Invalid secondary DNS format",
                "error_network_info": "Failed to get network information",
                
                "returned_to_main":  "Returned to main menu. Press Esc to exit.",
                "press_esc_return":  "Press ESC to return to main menu",
                
                "refresh":           "Refresh",
                "back":              "Back",
        },
        "pt": {
                "menu_title":        "Gerenciador de Rede TUI",
                "menu_configure":    "Configurar Rede",
                "menu_status":       "Status da Rede",
                "menu_ping_test":    "Teste de Ping",
                "menu_sysinfo":      "Informações do Sistema",
                "menu_reboot":       "Reiniciar",
                "menu_shutdown":     "Desligar",
                "menu_exit":         "Sair",
                "menu_language":     "Mudar Idioma",
                "menu_help":         "Ajuda",
                
                "network_title":     "Configurar Rede",
                "network_status":    "Status da Rede",
                "network_interface": "Interface de Rede (ex.: eth0):",
                "network_dhcp":      "DHCP:",
                "network_ipv4":      "Endereço IPv4:",
                "network_ipv6":      "Endereço IPv6:",
                "network_mac":       "Endereço MAC:",
                "network_netmask":   "Máscara de Rede:",
                "network_gateway":   "Gateway:",
                "network_dns":       "DNS:",
                "network_dns1":      "DNS Primário (ex.: 8.8.8.8):",
                "network_dns2":      "DNS Secundário (opcional):",
                "network_save":      "Salvar",
                "network_cancel":    "Cancelar",
                "network_back":      "Voltar",
                "network_refresh":   "Atualizar",
                "network_device":    "Dispositivo",
                "network_type":      "Tipo",
                "network_state":     "Status",
                "network_name":      "Nome",
                "network_ipv4_config": "Configuração IPv4",
                "network_ipv6_config": "Configuração IPv6",
                "network_ipv4_mode": "Modo IPv4:",
                "network_ipv6_mode": "Modo IPv6:",
                "network_ipv4_address": "Endereço IPv4:",
                "network_ipv4_netmask": "Máscara de Rede:",
                "network_ipv4_gateway": "Gateway IPv4:",
                "network_ipv4_dns1": "DNS Primário IPv4:",
                "network_ipv4_dns2": "DNS Secundário IPv4:",
                "network_ipv6_address": "Endereço IPv6:",
                "network_ipv6_prefix": "Tamanho do Prefixo IPv6:",
                "network_ipv6_gateway": "Gateway IPv6:",
                "network_ipv6_dns1": "DNS Primário IPv6:",
                "network_ipv6_dns2": "DNS Secundário IPv6:",
                
                "ping_title":        "Teste de Ping",
                "ping_target":       "Host Alvo (ex.: 8.8.8.8):",
                "ping_count":        "Contagem (opcional):",
                "ping_start":        "Iniciar",
                "ping_results":      "Resultados do Ping",
                
                "sysinfo_title":     "Informações do Sistema",
                "reboot_title":      "Reiniciar",
                "reboot_message":    "O sistema será reiniciado agora!",
                "shutdown_title":    "Desligar",
                "shutdown_message":  "O sistema será desligado agora!",
                
                "help_title":        "Ajuda",
                "help_description":  "Este aplicativo permite configurar e monitorar interfaces de rede.\n\nUse as teclas de seta para navegar, Enter para selecionar e Esc para voltar.",
                
                "success_title":     "Sucesso",
                "success_message":   "Configuração de rede aplicada com sucesso!",
                "error_title":       "Erro",
                "error_root_required": "Este aplicativo requer privilégios de root.",
                "error_empty_fields": "Por favor, preencha todos os campos obrigatórios",
                "error_invalid_ip": "Formato de endereço IP inválido",
                "error_invalid_netmask": "Formato de máscara de rede inválido (use número CIDR ou máscara completa)",
                "error_invalid_gateway": "Formato de gateway inválido",
                "error_invalid_dns1": "Formato de DNS primário inválido",
                "error_invalid_dns2": "Formato de DNS secundário inválido",
                "error_network_info": "Falha ao obter informações de rede",
                
                "returned_to_main":  "Retornou ao menu principal. Pressione Esc para sair.",
                "press_esc_return":  "Pressione ESC para voltar ao menu principal",
                
                "refresh":           "Atualizar",
                "back":              "Voltar",
        },
}

// Retorna a string traduzida
func T(key string) string {
        mu.Lock()
        defer mu.Unlock()
        
        // Se a tradução não existir, retorna a chave
        if val, ok := translations[language][key]; ok {
                return val
        }
        return key
}

// Define o idioma atual
func SetLanguage(lang string) {
        mu.Lock()
        defer mu.Unlock()
        language = lang
}

// Retorna o idioma atual
func GetLanguage() string {
        mu.Lock()
        defer mu.Unlock()
        return language
}