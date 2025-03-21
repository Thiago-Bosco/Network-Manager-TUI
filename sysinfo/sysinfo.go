package sysinfo

import (
        "fmt"
        "os"
        "runtime"
        "strconv"
        "strings"
        "syscall"
        "time"
        "regexp"
)

// Coleta informações detalhadas do sistema
func GetSystemInfo() string {
        // Data e hora atuais
        now := time.Now().Format("Mon Jan 2 15:04:05 MST 2006")

        // Coletar informações de hardware e do sistema
        cores, err := countCPUCores()
        if err != nil {
                cores = runtime.NumCPU() // Usar fallback seguro
        }

        cpuModel := getCPUModel()
        
        loadAvg, err := getLoadAverage()
        if err != nil {
                loadAvg = 0.0 // Valor padrão se não conseguir ler
        }

        memInfo, memPercent, err := getMemoryInfo()
        if err != nil {
                memInfo = "N/A" // Não disponível
                memPercent = 0.0
        }

        diskInfo, diskPercent, err := getDiskInfo()
        if err != nil {
                diskInfo = "N/A" // Não disponível
                diskPercent = 0.0
        }

        hostname, err := os.Hostname()
        if err != nil {
                hostname = "Unknown"
        }
        
        uptime := getUptime()
        kernelVer := getKernelVersion()
        
        // Formatar a saída com cores e melhor formatação
        output := fmt.Sprintf("[yellow]╭─────────────────────────────────────────────────────────────────╮[white]\n")
        output += fmt.Sprintf("[yellow]│[cyan]                   SYSTEM INFORMATION DASHBOARD                   [yellow]│[white]\n")
        output += fmt.Sprintf("[yellow]├─────────────────────────────────────────────────────────────────┤[white]\n")
        output += fmt.Sprintf("[yellow]│[white] [green]🕒 Date & Time:[white] %-47s [yellow]│[white]\n", now)
        output += fmt.Sprintf("[yellow]│[white] [green]🖥️  Hostname:[white]   %-47s [yellow]│[white]\n", hostname)
        output += fmt.Sprintf("[yellow]├─────────────────────────────────────────────────────────────────┤[white]\n")
        output += fmt.Sprintf("[yellow]│[cyan]                      SYSTEM SPECIFICATIONS                      [yellow]│[white]\n")
        output += fmt.Sprintf("[yellow]├─────────────────────────────────────────────────────────────────┤[white]\n")
        output += fmt.Sprintf("[yellow]│[white] [green]🐧 OS:[white]          %-47s [yellow]│[white]\n", runtime.GOOS)
        output += fmt.Sprintf("[yellow]│[white] [green]🔄 Kernel:[white]      %-47s [yellow]│[white]\n", kernelVer)
        output += fmt.Sprintf("[yellow]│[white] [green]⚙️  Architecture:[white] %-47s [yellow]│[white]\n", runtime.GOARCH)
        output += fmt.Sprintf("[yellow]│[white] [green]⏱️  Uptime:[white]      %-47s [yellow]│[white]\n", uptime)
        output += fmt.Sprintf("[yellow]├─────────────────────────────────────────────────────────────────┤[white]\n")
        output += fmt.Sprintf("[yellow]│[cyan]                        HARDWARE STATUS                         [yellow]│[white]\n")
        output += fmt.Sprintf("[yellow]├─────────────────────────────────────────────────────────────────┤[white]\n")
        output += fmt.Sprintf("[yellow]│[white] [green]🧠 CPU Model:[white]   %-47s [yellow]│[white]\n", cpuModel)
        output += fmt.Sprintf("[yellow]│[white] [green]📊 CPU Cores:[white]   %-47d [yellow]│[white]\n", cores)
        output += fmt.Sprintf("[yellow]│[white] [green]📈 Load Average:[white] %-47.2f [yellow]│[white]\n", loadAvg)
        
        // Barra de progresso para memória
        memBar := generateProgressBar(memPercent, 40)
        output += fmt.Sprintf("[yellow]│[white] [green]🧮 Memory:[white]      %s [%5.1f%%] [yellow]│[white]\n", memBar, memPercent)
        output += fmt.Sprintf("[yellow]│[white]                %s [yellow]│[white]\n", memInfo)
        
        // Barra de progresso para disco
        diskBar := generateProgressBar(diskPercent, 40)
        output += fmt.Sprintf("[yellow]│[white] [green]💾 Disk:[white]        %s [%5.1f%%] [yellow]│[white]\n", diskBar, diskPercent)
        output += fmt.Sprintf("[yellow]│[white]                %s [yellow]│[white]\n", diskInfo)
        
        output += fmt.Sprintf("[yellow]╰─────────────────────────────────────────────────────────────────╯[white]\n")

        return output
}

// Gera uma barra de progresso colorida
func generateProgressBar(percent float64, width int) string {
    filledWidth := int(percent / 100.0 * float64(width))
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