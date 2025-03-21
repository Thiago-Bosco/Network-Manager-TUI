package sysinfo

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func GetSystemInfo() string {
	now := time.Now().Format("Mon Jan 2 15:04:05 MST 2006")
	cores, err := countCPUCores()
	if err != nil {
		cores = runtime.NumCPU()
	}

	cpuModel := getCPUModel()

	loadAvg, err := getLoadAverage()
	if err != nil {
		loadAvg = 0.0
	}

	memInfo, memPercent, err := getMemoryInfo()
	if err != nil {
		memInfo = "N/A"
		memPercent = 0.0
	}

	diskInfo, diskPercent, err := getDiskInfo()
	if err != nil {
		diskInfo = "N/A"
		diskPercent = 0.0
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	uptime := getUptime()
	kernelVer := getKernelVersion()

	width := 65 // Total width for centering
	centerText := func(text string) string {
		textLen := len(stripColor(text))
		if textLen >= width {
			return text
		}
		padding := (width - textLen) / 2
		if padding < 0 {
			padding = 0
		}
		return strings.Repeat(" ", padding) + text
	}

	output := fmt.Sprintf("[yellow]%s[white]\n", centerText("╭─────────────────────────────────────────────────────────────────╮"))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText("│[cyan]                   SYSTEM INFORMATION DASHBOARD                   [yellow]│"))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText("├─────────────────────────────────────────────────────────────────┤"))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText(fmt.Sprintf("│[white] [green]🕒 Date & Time:[white] %-47s [yellow]│", now)))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText(fmt.Sprintf("│[white] [green]🖥️  Hostname:[white]   %-47s [yellow]│", hostname)))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText("├─────────────────────────────────────────────────────────────────┤"))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText("│[cyan]                      SYSTEM SPECIFICATIONS                      [yellow]│"))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText("├─────────────────────────────────────────────────────────────────┤"))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText(fmt.Sprintf("│[white] [green]🐧 OS:[white]          %-47s [yellow]│", runtime.GOOS)))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText(fmt.Sprintf("│[white] [green]🔄 Kernel:[white]      %-47s [yellow]│", kernelVer)))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText(fmt.Sprintf("│[white] [green]⚙️  Architecture:[white] %-47s [yellow]│", runtime.GOARCH)))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText(fmt.Sprintf("│[white] [green]⏱️  Uptime:[white]      %-47s [yellow]│", uptime)))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText("├─────────────────────────────────────────────────────────────────┤"))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText("│[cyan]                        HARDWARE STATUS                         [yellow]│"))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText("├─────────────────────────────────────────────────────────────────┤"))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText(fmt.Sprintf("│[white] [green]🧠 CPU Model:[white]   %-47s [yellow]│", cpuModel)))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText(fmt.Sprintf("│[white] [green]📊 CPU Cores:[white]   %-47d [yellow]│", cores)))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText(fmt.Sprintf("│[white] [green]📈 Load Average:[white] %-47.2f [yellow]│", loadAvg)))

	memBar := generateProgressBar(memPercent, 40)
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText(fmt.Sprintf("│[white] [green]🧮 Memory:[white]      %s [%5.1f%%] [yellow]│", memBar, memPercent)))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText(fmt.Sprintf("│[white]                %s [yellow]│", memInfo)))

	diskBar := generateProgressBar(diskPercent, 40)
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText(fmt.Sprintf("│[white] [green]💾 Disk:[white]        %s [%5.1f%%] [yellow]│", diskBar, diskPercent)))
	output += fmt.Sprintf("[yellow]%s[white]\n", centerText(fmt.Sprintf("│[white]                %s [yellow]│", diskInfo)))

	output += fmt.Sprintf("[yellow]%s[white]\n", centerText("╰─────────────────────────────────────────────────────────────────╯"))

	return output
}

// Helper function to strip color codes for proper centering
func stripColor(text string) string {
	re := regexp.MustCompile(`\[[^\]]*\]`)
	return re.ReplaceAllString(text, "")
}

// Gera uma barra de progresso colorida
func generateProgressBar(percent float64, width int) string {
	filledWidth := int(percent/100.0*float64(width))
	emptyWidth := width - filledWidth

	var color string
	if percent < 60 {
		color = "[green]"
	} else if percent < 85 {
		color = "[yellow]"
	} else {
		color = "[red]"
	}

	bar := color + strings.Repeat("█", filledWidth) + "[white]" + strings.Repeat("░", emptyWidth)
	return bar
}

// Obtém o modelo da CPU
func getCPUModel() string {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "Unknown CPU"
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "model name") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}

	return "Unknown CPU"
}

// Obtém a versão do kernel
func getKernelVersion() string {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return runtime.Version() // Fallback para a versão do Go runtime
	}

	// Extrai apenas a versão do kernel
	re := regexp.MustCompile(`Linux version ([^ ]+)`)
	matches := re.FindStringSubmatch(string(data))
	if len(matches) > 1 {
		return matches[1]
	}

	return string(data)[:30] + "..." // Trunca se for muito longo
}

// Obtém o tempo de atividade do sistema
func getUptime() string {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "Unknown"
	}

	parts := strings.Split(string(data), " ")
	if len(parts) < 1 {
		return "Unknown"
	}

	uptime, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return "Unknown"
	}

	// Converte segundos em um formato mais legível
	days := int(uptime / 86400)
	hours := int((uptime - float64(days)*86400) / 3600)
	minutes := int((uptime - float64(days)*86400 - float64(hours)*3600) / 60)

	result := ""
	if days > 0 {
		result += fmt.Sprintf("%d dias, ", days)
	}
	result += fmt.Sprintf("%d horas, %d minutos", hours, minutes)

	return result
}

// Conta o número de núcleos da CPU
func countCPUCores() (int, error) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return runtime.NumCPU(), nil // Fallback para runtime.NumCPU() se não conseguir ler o arquivo
	}
	lines := strings.Split(string(data), "\n")
	count := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "processor") {
			count++
		}
	}
	if count == 0 {
		return runtime.NumCPU(), nil
	}
	return count, nil
}

// Obtém a média de carga do sistema
func getLoadAverage() (float64, error) {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return 0, err
	}
	loadAvg := strings.Split(string(data), " ")[0]
	return strconv.ParseFloat(loadAvg, 64)
}

// Obtém informações sobre a memória
func getMemoryInfo() (string, float64, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "", 0.0, err
	}

	lines := strings.Split(string(data), "\n")
	var total, free int
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			total, err = strconv.Atoi(strings.Fields(line)[1])
			if err != nil {
				return "", 0.0, err
			}
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			free, err = strconv.Atoi(strings.Fields(line)[1])
			if err != nil {
				return "", 0.0, err
			}
		}
	}
	used := total - free
	usedPercent := float64(used) * 100.0 / float64(total)
	memUsage := fmt.Sprintf("%.2f GB / %.2f GB", float64(used)/1048576, float64(total)/1048576)
	return memUsage, usedPercent, nil
}

// Obtém informações sobre o uso do disco
func getDiskInfo() (string, float64, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		return "", 0.0, err
	}

	total := stat.Blocks * uint64(stat.Bsize) // total em bytes
	free := stat.Bfree * uint64(stat.Bsize)   // livre em bytes
	used := total - free                      // usado em bytes

	// Calculando a porcentagem de uso
	usedPercent := float64(used) * 100.0 / float64(total)

	// Convertendo para GB com duas casas decimais
	totalGB := float64(total) / (1024 * 1024 * 1024)
	usedGB := float64(used) / (1024 * 1024 * 1024)
	freeGB := float64(free) / (1024 * 1024 * 1024)

	diskUsage := fmt.Sprintf("%.2f GB / %.2f GB (%.2f%% free)",
		usedGB, totalGB, (freeGB/totalGB)*100)
	return diskUsage, usedPercent, nil
}