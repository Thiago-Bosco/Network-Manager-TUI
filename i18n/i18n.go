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
                "menu_sysinfo":      "System Information",
                "menu_reboot":       "Reboot",
                "menu_shutdown":     "ShutDown",
                "menu_exit":         "Exit",
                "menu_language":     "Change Language",
                "network_title":     "Configure Network",
                "network_interface": "Network Interface (e.g., eth0):",
                "network_dhcp":      "DHCP:",
                "network_ipv4":      "IPv4 Address:",
                "network_netmask":   "Netmask:",
                "network_gateway":   "Gateway:",
                "network_dns1":      "Primary DNS (e.g., 8.8.8.8):",
                "network_dns2":      "Secondary DNS (optional):",
                "network_save":      "Save",
                "network_cancel":    "Cancel",
                "sysinfo_title":     "System Information",
                "reboot_title":      "Reboot",
                "reboot_message":    "The system will now reboot!",
                "shutdown_title":    "ShutDown",
                "shutdown_message":  "The system will now shut down!",
                "success_title":     "Success",
                "success_message":   "Network configuration applied successfully!",
                "error_title":       "Error",
                "error_root_required": "This application requires root privileges.",
                "error_empty_fields": "Please fill all required fields",
                "error_invalid_ip": "Invalid IP address format",
                "error_invalid_netmask": "Invalid netmask format (use CIDR number or full mask)",
                "error_invalid_gateway": "Invalid gateway IP address format",
                "error_invalid_dns1": "Invalid primary DNS format",
                "error_invalid_dns2": "Invalid secondary DNS format",
                "returned_to_main":  "Returned to main menu. Press Esc to exit.",
        },
        "pt": {
                "menu_title":        "Gerenciador de Rede TUI",
                "menu_configure":    "Configurar Rede",
                "menu_sysinfo":      "Informações do Sistema",
                "menu_reboot":       "Reiniciar",
                "menu_shutdown":     "Desligar",
                "menu_exit":         "Sair",
                "menu_language":     "Mudar Idioma",
                "network_title":     "Configurar Rede",
                "network_interface": "Interface de Rede (ex.: eth0):",
                "network_dhcp":      "DHCP:",
                "network_ipv4":      "Endereço IPv4:",
                "network_netmask":   "Máscara de Rede:",
                "network_gateway":   "Gateway:",
                "network_dns1":      "DNS Primário (ex.: 8.8.8.8):",
                "network_dns2":      "DNS Secundário (opcional):",
                "network_save":      "Salvar",
                "network_cancel":    "Cancelar",
                "sysinfo_title":     "Informações do Sistema",
                "reboot_title":      "Reiniciar",
                "reboot_message":    "O sistema será reiniciado agora!",
                "shutdown_title":    "Desligar",
                "shutdown_message":  "O sistema será desligado agora!",
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
                "returned_to_main":  "Retornou ao menu principal. Pressione Esc para sair.",
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