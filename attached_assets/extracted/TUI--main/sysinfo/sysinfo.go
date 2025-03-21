package sysinfo

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Coleta informações detalhadas do sistema
func GetSystemInfo() string {
	// Data e hora atuais
	now := time.Now().Format("Mon Jan 2 15:04:05 MST 2006")

	// Coletar informações de hardware e do sistema
	cores, err := countCPUCores()
	if err != nil {
		return "Error: Failed to count CPU cores."
	}

	loadAvg, err := getLoadAverage()
	if err != nil {
		return "Error: Failed to get load average."
	}

	memInfo, err := getMemoryInfo()
	if err != nil {
		return "Error: Failed to get memory info."
	}

	diskInfo, err := getDiskInfo()
	if err != nil {
		return "Error: Failed to get disk info."
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown"
	}

	// saída
	output := fmt.Sprintf("\nSystem Information as of \n%s\n\n", now)
	output += fmt.Sprintf("|----------------------------------------------|\n")
	output += fmt.Sprintf("|%-15s |%-30s\n", "Host", hostname)
	output += fmt.Sprintf("|%-15s |%-30s\n", "OS", runtime.GOOS)
	output += fmt.Sprintf("|%-15s |%-30s\n", "Architecture", runtime.GOARCH)
	output += fmt.Sprintf("|%-15s |%-30d\n", "CPU Cores", cores)
	output += fmt.Sprintf("|%-15s |%-30.2f\n", "Load Average", loadAvg)
	output += fmt.Sprintf("|%-15s |%-30s\n", "Memory Usage", memInfo)
	output += fmt.Sprintf("|%-15s |%-15s\n", "Disk Usage", diskInfo)
	output += fmt.Sprintf("|----------------------------------------------|\n")

	return output
}

// Conta o número de núcleos da CPU
func countCPUCores() (int, error) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(data), "\n")
	count := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "processor") {
			count++
		}
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
func getMemoryInfo() (string, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	var total, free int
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			total, err = strconv.Atoi(strings.Fields(line)[1])
			if err != nil {
				return "", err
			}
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			free, err = strconv.Atoi(strings.Fields(line)[1])
			if err != nil {
				return "", err
			}
		}
	}
	used := total - free
	memUsage := fmt.Sprintf("%.2f GB / %.2f GB", float64(used)/1048576, float64(total)/1048576)
	return memUsage, nil
}

// Obtém informações sobre o uso do disco
func getDiskInfo() (string, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		return "", err
	}

	total := stat.Blocks * uint64(stat.Bsize) // total em bytes
	free := stat.Bfree * uint64(stat.Bsize)   // livre em bytes
	used := total - free                      // usado em bytes

	// Convertendo para GB com duas casas decimais
	totalGB := float64(total) / (1024 * 1024 * 1024)
	usedGB := float64(used) / (1024 * 1024 * 1024)
	freeGB := float64(free) / (1024 * 1024 * 1024)

	diskUsage := fmt.Sprintf("(Used:%.2fGB)(Total:%.2fGB)\n(Free:%.2fGB)", usedGB, totalGB, freeGB)
	return diskUsage, nil
}
