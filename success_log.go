package logging

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

// Structure to represent the action log as JSON
type ActionLog struct {
	Timestamp  string      `json:"timestamp"`
	ActionType string      `json:"action_type"`
	EntityName string      `json:"entity_name"`
	Attributes interface{} `json:"attributes"`
}

// ActionType defines possible action types
type ActionType int

const (
	Create ActionType = iota
	Update
	Delete
)

// LogCreatedSuccess registers successful creation logs
func LogCreatedSuccess(entityName string, attributes interface{}) {
	logFolderPath := os.Getenv("CREATED_SUCCESS_FOLDER_PATH")
	logActionSuccess("CREATE", entityName, attributes, logFolderPath)
}

// LogUpdatedSuccess registers successful update logs
func LogUpdatedSuccess(entityName string, attributes interface{}) {
	logFolderPath := os.Getenv("UPDATED_SUCCESS_FOLDER_PATH")
	logActionSuccess("UPDATE", entityName, attributes, logFolderPath)
}

// LogDeletedSuccess registers successful deletion logs
func LogDeletedSuccess(entityName string, attributes interface{}) {
	logFolderPath := os.Getenv("DELETED_SUCCESS_FOLDER_PATH")
	logActionSuccess("DELETE", entityName, attributes, logFolderPath)
}

func logActionSuccess(actionType string, entityName string, attributes interface{}, logFolderPath string) {
	// Check if the log directory exists. If not, create it
	err := os.MkdirAll(logFolderPath, 0755)
	if err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// Create the full path for the action log file
	actionLogFileName := filepath.Join(logFolderPath, "action_"+time.Now().Format("2006-01-02")+".log")

	// Open or create the action log file
	actionLogFile, err := os.OpenFile(actionLogFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open action log file: %v", err)
	}
	defer actionLogFile.Close()

	// Create an instance of the ActionLog structure
	logEntry := ActionLog{
		Timestamp:  time.Now().Format(time.RFC3339),
		ActionType: actionType,
		EntityName: entityName,
		Attributes: attributes,
	}

	// Serialize the action log as JSON
	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		log.Fatalf("Error serializing action log to JSON: %v", err)
	}

	// Write the action log to the file
	_, err = actionLogFile.Write(jsonData)
	if err != nil {
		log.Fatalf("Error writing to action log file: %v", err)
	}

	// Convert the object to JSON for printing to console
	attributesJSON, err := json.Marshal(attributes)
	if err != nil {
		log.Fatalf("Error serializing attributes to JSON: %v", err)
	}

	cyan := color.New(color.FgWhite, color.BgHiCyan).SprintFunc()
	log.SetPrefix("[APP] ")
	log.Printf("%s %s: %s\n", withColor(actionType), cyan(" "+entityName+" "), string(attributesJSON))
}

func withColor(actionType string) string {
	var actionColor func(a ...interface{}) string
	switch actionType {
	case "CREATE":
		actionColor = color.New(color.FgWhite, color.BgGreen).SprintFunc()
	case "UPDATE":
		actionColor = color.New(color.FgWhite, color.BgYellow).SprintFunc()
	case "DELETE":
		actionColor = color.New(color.FgWhite, color.BgRed).SprintFunc()
	}

	return actionColor(" " + actionType + " ")
}
