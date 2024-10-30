// internal/logging/logging.go
package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

// Configuration constants with default values
const (
	// Default log settings if env vars are not set
	DefaultLogDir        = "/tmp/log/nanodns"
	DefaultServiceLog    = "service.log"
	DefaultActionLog     = "actions.log"
	DefaultMaxLogSize    = 10 * 1024 * 1024 // 10MB
	DefaultMaxLogBackups = 5
)

// Config holds the logging configuration
type Config struct {
	LogDir         string
	ServiceLogFile string
	ActionLogFile  string
	MaxLogSize     int64
	MaxLogBackups  int
}

var (
	serviceLogger *log.Logger
	actionLogger  *log.Logger
	config        Config
	serviceFile   *os.File // Keep track of the service log file
	actionFile    *os.File // Keep track of the action log file
	fileMutex     sync.Mutex
)

// loadConfig loads configuration from environment variables with fallbacks to defaults
func loadConfig() Config {
	return Config{
		LogDir:         getEnv("LOG_DIR", DefaultLogDir),
		ServiceLogFile: getEnv("SERVICE_LOG", DefaultServiceLog),
		ActionLogFile:  getEnv("ACTION_LOG", DefaultActionLog),
		MaxLogSize:     getEnvInt64("MAX_LOG_SIZE", DefaultMaxLogSize),
		MaxLogBackups:  getEnvInt("MAX_LOG_BACKUPS", DefaultMaxLogBackups),
	}
}

// Helper function to get environment variables with default fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper function to get integer environment variables
func getEnvInt(key string, fallback int) int {
	strValue, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	intValue, err := strconv.Atoi(strValue)
	if err != nil {
		return fallback
	}
	return intValue
}

// Helper function to get int64 environment variables
func getEnvInt64(key string, fallback int64) int64 {
	strValue, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}

	int64Value, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return fallback
	}
	return int64Value
}

// Init initializes the logging system
func Init() error {
	// Load configuration from environment
	config = loadConfig()

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// Initialize service logger
	svcFile, err := GetServiceLogFile()
	if err != nil {
		return fmt.Errorf("failed to initialize service log: %v", err)
	}
	if err != nil {
		return fmt.Errorf("failed to open service log file: %v", err)
	}

	// Initialize action logger
	actionLog, err := os.OpenFile(
		filepath.Join(config.LogDir, config.ActionLogFile),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return fmt.Errorf("failed to open action log file: %v", err)
	}

	actionFile = actionLog

	// Create loggers with timestamps and prefixes
	serviceLogger = log.New(svcFile, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	actionLogger = log.New(actionLog, "", log.Ldate|log.Ltime|log.Lmicroseconds)

	return nil
}

// GetLogPaths returns the current log file paths
func GetLogPaths() (service, action string) {
	return filepath.Join(config.LogDir, config.ServiceLogFile),
		filepath.Join(config.LogDir, config.ActionLogFile)
}

// LogAction logs administrative actions
func LogAction(action, details string) {
	if actionLogger != nil {
		actionLogger.Printf("%s - %s", action, details)
	}
}

// LogService logs DNS service operations
func LogService(message string) {
	if serviceLogger != nil {
		serviceLogger.Print(message)
	}
}

// GetActionLogs retrieves action logs with optional filtering
func GetActionLogs(since time.Duration) ([]string, error) {
	return readLogFile(filepath.Join(config.LogDir, config.ActionLogFile), since)
}

// GetServiceLogs retrieves service logs with optional filtering
func GetServiceLogs(since time.Duration) ([]string, error) {
	return readLogFile(filepath.Join(config.LogDir, config.ServiceLogFile), since)
}

// Helper function to read and filter log files
func readLogFile(filepath string, since time.Duration) ([]string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %v", err)
	}

	// TODO: Implement log filtering based on the 'since' parameter
	return []string{string(content)}, nil
}

// RotateLogs checks log sizes and rotates if necessary
func RotateLogs() error {
	files := []string{config.ServiceLogFile, config.ActionLogFile}

	for _, file := range files {
		path := filepath.Join(config.LogDir, file)

		info, err := os.Stat(path)
		if err != nil {
			continue // Skip if file doesn't exist
		}

		if info.Size() > config.MaxLogSize {
			if err := rotateLogFile(path); err != nil {
				return fmt.Errorf("failed to rotate %s: %v", file, err)
			}
		}
	}

	return nil
}

// Helper function to rotate a single log file
func rotateLogFile(filepath string) error {
	// Remove the oldest backup if it exists
	oldestBackup := fmt.Sprintf("%s.%d", filepath, config.MaxLogBackups)
	_ = os.Remove(oldestBackup)

	// Shift existing backups
	for i := config.MaxLogBackups - 1; i > 0; i-- {
		oldPath := fmt.Sprintf("%s.%d", filepath, i)
		newPath := fmt.Sprintf("%s.%d", filepath, i+1)
		_ = os.Rename(oldPath, newPath)
	}

	// Rename current log to .1
	return os.Rename(filepath, filepath+".1")
}

func GetServiceLogFile() (*os.File, error) {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	// If we already have an open file, return it
	if serviceFile != nil {
		// Check if file is still valid
		if _, err := serviceFile.Stat(); err == nil {
			return serviceFile, nil
		}
		// If file is not valid, close it and we'll reopen
		serviceFile.Close()
		serviceFile = nil
	}

	// Open a new file handle
	logPath := filepath.Join(config.LogDir, config.ServiceLogFile)

	// Ensure directory exists
	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	// Open the file with append mode
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open service log file: %v", err)
	}

	// Store the file handle
	serviceFile = file

	return serviceFile, nil
}

// CloseServiceLogFile closes the service log file if it's open
func CloseServiceLogFile() error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	if serviceFile != nil {
		err := serviceFile.Close()
		serviceFile = nil
		return err
	}
	return nil
}

func Cleanup() {
	CloseServiceLogFile()
	if actionFile != nil {
		actionFile.Close()
		actionFile = nil
	}
}
