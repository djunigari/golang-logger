package logger

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

// LogError registra logs de erro em um arquivo separado para o dia atual
func LogError(errorMessage string, errorDetails string) {
	logFolderPath := os.Getenv("LOGGER_FOLDER_PATH")

	err := os.MkdirAll(logFolderPath, 0755)
	if err != nil {
		log.Fatalf("Falha ao criar o diretório de logs: %v", err)
	}

	// Criar o caminho completo para o arquivo de log de erro
	errLogFileName := filepath.Join(logFolderPath, "error_"+time.Now().Format("2006-01-02")+".log")

	// Abrir ou criar o arquivo de log de erro
	errLogFile, err := os.OpenFile(errLogFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Falha ao abrir o arquivo de log de erro: %v", err)
	}
	defer errLogFile.Close()

	// Estrutura para representar o log de erro como JSON
	type ErrorLog struct {
		Timestamp    string `json:"timestamp"`
		ErrorMessage string `json:"error_message"`
		ErrorDetails string `json:"error_details"`
	}

	// Criar uma instância da estrutura ErrorLog
	logEntry := ErrorLog{
		Timestamp:    time.Now().Format(time.RFC3339),
		ErrorMessage: errorMessage,
		ErrorDetails: errorDetails,
	}

	// Serializar o log de erro como JSON
	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		log.Fatalf("Erro ao serializar o log de erro como JSON: %v", err)
	}

	// Escrever o log de erro no arquivo
	_, err = errLogFile.Write(jsonData)
	if err != nil {
		log.Fatalf("Erro ao escrever no arquivo de log de erro: %v", err)
	}

	// Imprimir no console com cor
	red := color.New(color.FgRed).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	log.SetPrefix("[APP] ")
	log.Printf("%s %s: %s\n", red("ERROR"), cyan("Message"), errorMessage)

	if errorDetails != "" {
		log.Printf("%s %s: %s\n", red("ERROR"), cyan("Details"), errorDetails)
	}
}
