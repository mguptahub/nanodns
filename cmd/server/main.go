package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mguptahub/nanodns/internal/dns"
	"github.com/mguptahub/nanodns/internal/logging"
	"github.com/mguptahub/nanodns/pkg/config"
	externaldns "github.com/miekg/dns"
)

const pidFilePath = "/tmp/nanodns.pid"

var logDuration time.Duration

var (
	version      = "dev" // Default version; can be overridden at build time
	startService bool
	stopService  bool
	showVersion  bool
	showHelp     bool
	showStatus   bool
	showLogs     bool
)

func init() {
	flag.DurationVar(&logDuration, "duration", 24*time.Hour, "Duration of logs to show (e.g., 72h, 7d)")

	flag.BoolVar(&startService, "start", false, "Run the binary as a daemon")
	flag.BoolVar(&stopService, "stop", false, "Stop the running daemon service")
	flag.BoolVar(&showVersion, "version", false, "Show the binary version")
	flag.BoolVar(&showVersion, "v", false, "Show the binary version (short)")
	flag.BoolVar(&showHelp, "help", false, "Shows help information")
	flag.BoolVar(&showStatus, "status", false, "Show service status")

	flag.BoolVar(&showLogs, "logs", false, "Show the logs")

	// Custom usage function
	flag.Usage = func() {
		fmt.Println("")
		fmt.Println("Usage: nanodns [command | options]")
		fmt.Println("")
		fmt.Println("commands:")
		fmt.Println("  start                              Run the binary as a daemon")
		fmt.Println("  stop                               Stop the running daemon service")
		fmt.Println("  status                             Show service status")
		fmt.Println("  logs                               Show service logs")
		fmt.Println("")
		fmt.Println("options:")
		fmt.Println("  -v | --version                     Show the binary version")
		fmt.Println("  -a | --action-logs                 Show the action logs. This works with the logs command")
	}

	config.Initialize()

	// Initialize logging system
	if err := logging.Init(); err != nil {
		log.Fatalf("Failed to initialize logging: %v", err)
	}
}

func main() {
	flag.Parse()
	if len(os.Args) > 1 {
		switch flag.Arg(0) {
		case "start":
			startDaemon()
			return
		case "stop":
			stopDaemon()
			return
		case "status":
			checkServiceStatus()
			return
		case "logs", "log":
			showSelectiveLogs()
			return
		case "help":
			flag.Usage()
			return
		default:
			flag.Usage()
			return
		}
	}

	if showHelp {
		flag.Usage()
		return
	}

	if showVersion {
		fmt.Printf("NanoDNS Version: %s\n", version)
		return
	}

	if startService {
		if checkIfRunning() {
			fmt.Println("")
			fmt.Println("NanoDNS is already running.")
			fmt.Println("")
			return
		}
		fmt.Println("")
		fmt.Println("Running as daemon...")
		fmt.Println("")
		startDaemon()
		return
	}

	if stopService {
		stopDaemon()
		return
	}

	if showStatus {
		checkServiceStatus()
		return
	}

	// Regular server startup if no flags are provided
	startDNSServer()
}

func startDNSServer() {
	logging.LogService("Initializing DNS server")

	// Load records from environment variables
	records := dns.LoadRecords()
	logging.LogService(fmt.Sprintf("Loaded %d DNS records", len(records)))

	// Get relay configuration
	relayConfig := config.GetRelayConfig()
	if relayConfig.Enabled {
		logging.LogService(fmt.Sprintf("DNS relay enabled, using nameservers: %v", relayConfig.Nameservers))
	}

	// Create DNS handler
	handler, err := dns.NewHandler(records, relayConfig)
	if err != nil {
		logging.LogService(fmt.Sprintf("Failed to create DNS handler: %v", err))
		log.Fatalf("Failed to create DNS handler: %v", err)
	}
	externaldns.HandleFunc(".", handler.ServeDNS)

	// Configure server
	port := config.GetDNSPort()
	server := &externaldns.Server{
		Addr: ":" + port,
		Net:  "udp",
	}

	logging.LogService(fmt.Sprintf("Starting DNS server on port %s", port))
	if err := server.ListenAndServe(); err != nil {
		logging.LogService(fmt.Sprintf("Failed to start server: %v", err))
		log.Fatalf("Failed to start server: %v", err)
	}

	defer func() {
		if err := server.Shutdown(); err != nil {
			logging.LogService(fmt.Sprintf("Error during server shutdown: %v", err))
		}
	}()
}

func startDaemon() {
	if checkIfRunning() {
		logging.LogAction("START_ATTEMPT", "Server already running")
		fmt.Println("NanoDNS is already running.")
		return
	}

	// Get service log file for output redirection
	logFile, err := logging.GetServiceLogFile()
	if err != nil {
		logging.LogAction("START_FAILED", fmt.Sprintf("Failed to open service log: %v", err))
		return // Just exit the function
	}

	// cmd := exec.Command(os.Args[0])
	cmd := exec.Command(os.Args[0])
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}
	// cmd.Dir = config.WorkingDir

	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		logging.LogAction("START_FAILED", fmt.Sprintf("Failed to start daemon: %v", err))
		return // Just exit the function
	}

	if err := writePID(cmd.Process.Pid); err != nil {
		if err := cmd.Process.Kill(); err != nil {
			logging.LogAction("PROCESS_KILL_FAILED", fmt.Sprintf("Failed to kill process: %v", err))
		}
		log.Fatalf("Failed to write PID: %v", err)
	}

	logging.LogAction("START_SUCCESS", fmt.Sprintf("Server started with PID %d", cmd.Process.Pid))
	fmt.Printf("Server running in background with PID %d\n", cmd.Process.Pid)
}

func writePID(pid int) error {
	return os.WriteFile(pidFilePath, []byte(strconv.Itoa(pid)), 0644)
}

func checkIfRunning() bool {
	if _, err := os.Stat(pidFilePath); os.IsNotExist(err) {
		return false // PID file does not exist
	}

	data, err := os.ReadFile(pidFilePath)
	if err != nil {
		log.Fatalf("Failed to read PID file: %v", err)
	}
	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		log.Fatalf("Invalid PID in file: %v", err)
	}

	// Check if the process is still running
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Try to signal the process with zero (doesn't send a signal, just checks if it's valid)
	err = process.Signal(syscall.Signal(0)) // Just checking if the process exists
	return err == nil
}

func checkServiceStatus() {
	fmt.Println("")
	if checkIfRunning() {
		fmt.Println("The DNS server is currently running.")
	} else {
		fmt.Println("The DNS server is not running.")
	}
	fmt.Println("")
}

func stopDaemon() {
	pidData, err := os.ReadFile(pidFilePath)
	if err != nil {
		logging.LogAction("STOP_ATTEMPT", "No PID file found")
		fmt.Println("No PID file found. NanoDNS is not running.")
		return
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(pidData)))
	if err != nil {
		logging.LogAction("STOP_FAILED", fmt.Sprintf("Invalid PID in file: %v", err))
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		logging.LogAction("STOP_FAILED", fmt.Sprintf("Failed to find process with PID %d: %v", pid, err))
		return
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		logging.LogAction("STOP_FAILED", fmt.Sprintf("Failed to stop process: %v", err))
		log.Fatalf("Failed to stop NanoDNS: %v", err)
	}

	logging.LogAction("STOP_SUCCESS", fmt.Sprintf("Server stopped (PID: %d)", pid))
	fmt.Printf("NanoDNS stopped successfully (PID: %d)\n", pid)

	if err := os.Remove(pidFilePath); err != nil {
		logging.LogAction("PID_REMOVE_FAILED", fmt.Sprintf("Failed to remove PID file: %v", err))
	}
}

func showSelectiveLogs() {
	if len(os.Args) > 2 {
		switch os.Args[2] {
		case "-a", "--action-logs":
			showActionLogs()
			return
		default:
			showServiceLogs()
			return
		}
	} else {
		showServiceLogs()
		return
	}
}

func showServiceLogs() {
	// Use the new configuration-aware logging package
	logs, err := logging.GetServiceLogs(logDuration)
	if err != nil {
		fmt.Printf("Failed to read service logs: %v\n", err)
		return
	}

	for _, entry := range logs {
		fmt.Println(entry)
	}
}
func showActionLogs() {
	logs, err := logging.GetActionLogs(logDuration)
	if err != nil {
		fmt.Printf("Failed to read action logs: %v\n", err)
		return
	}

	for _, entry := range logs {
		fmt.Println(entry)
	}
}
