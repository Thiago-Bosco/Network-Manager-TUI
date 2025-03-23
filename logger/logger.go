
package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	logFile     *os.File
)

func cleanOldLogs() error {
	logDir := "logs"
	files, err := os.ReadDir(logDir)
	if err != nil {
		return fmt.Errorf("erro ao ler diretório de logs: %v", err)
	}

	cutoffDate := time.Now().AddDate(0, 0, -90)
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(logDir, file.Name())
			fileInfo, err := file.Info()
			if err != nil {
				continue
			}
			
			if fileInfo.ModTime().Before(cutoffDate) {
				os.Remove(filePath)
				LogInfo("Log antigo removido: %s", file.Name())
			}
		}
	}
	return nil
}

func Init() error {
	// Cria diretório de logs se não existir
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório de logs: %v", err)
	}

	// Limpa logs antigos
	if err := cleanOldLogs(); err != nil {
		LogError("Erro ao limpar logs antigos: %v", err)
	}

	// Nome do arquivo de log com timestamp
	logFileName := filepath.Join(logDir, fmt.Sprintf("system_%s.log", time.Now().Format("2006-01-02")))
	
	// Abre arquivo de log
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de log: %v", err)
	}

	logFile = file
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime)
	
	return nil
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

func LogInfo(format string, v ...interface{}) {
	if InfoLogger != nil {
		InfoLogger.Printf(format, v...)
	}
}

func LogError(format string, v ...interface{}) {
	if ErrorLogger != nil {
		ErrorLogger.Printf(format, v...)
	}
}
