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
		"menu_title":        "F-Safer",
		"menu_configure":    "Configure Network",
		"menu_sysinfo":      "System Information",
		"menu_reboot":       "Reboot",
		"menu_shutdown":     "ShutDown",
		"menu_exit":         "Exit",
		"network_title":     "Configure Network",
		"network_interface": "Network Interface (e.g., eth0):",
		"network_dhcp":      "DHCP:",
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
		
	},
	"pt": {
		"menu_title":        "F-Safer",
		"menu_configure":    "Configurar Rede",
		"menu_sysinfo":      "Informações do Sistema",
		"menu_reboot":       "Reiniciar",
		"menu_shutdown":     "Desligar",
		"menu_exit":         "Sair",
		"network_dhcp":      "DHCP:",

		"network_title":     "Configurar Rede",
		"network_interface": "Interface de Rede (ex.: eth0):",
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
	},
}

// Retorna a string traduzida
func T(key string) string {
	mu.Lock()
	defer mu.Unlock()
	return translations[language][key]
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
