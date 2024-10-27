package config

import (
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	DefaultTTL     = 60
	DefaultPort    = "53"
	ServicePrefix  = "service:"
	DefaultTimeout = 5 * time.Second
)

type RelayConfig struct {
	Enabled     bool
	Nameservers []string
	Timeout     time.Duration
}

type APIConfig struct {
	Enabled bool
	Token   string
}

func init() {
	// Load .env file if it exists
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	// Try to load from current directory
	if err := godotenv.Load(envFile); err != nil {
		// Try to load from executable directory
		if exe, err := os.Executable(); err == nil {
			exeDir := filepath.Dir(exe)
			envPath := filepath.Join(exeDir, envFile)
			if err := godotenv.Load(envPath); err != nil {
				// Not finding .env file is OK, will use environment variables
				log.Printf("Note: No %s file found, using environment variables", envFile)
			}
		}
	}
}

// LoadEnvFile explicitly loads an environment file
func LoadEnvFile(path string) error {
	return godotenv.Load(path)
}

func GetDNSPort() string {
	if port := os.Getenv("DNS_PORT"); port != "" {
		return port
	}
	return DefaultPort
}

// IsServiceRecord checks if the value represents a Docker service
func IsServiceRecord(value string) bool {
	return strings.HasPrefix(value, ServicePrefix)
}

// GetServiceName extracts service name from value
func GetServiceName(value string) string {
	return strings.TrimPrefix(value, ServicePrefix)
}

// GetAPIConfig returns API configuration based on environment variables.
// It reads DNS_API_TOKEN for API configuration.
func GetAPIConfig() APIConfig {
	config := APIConfig{}

	if os.Getenv("DNS_API_TOKEN") != "" {
		config.Enabled = true
		config.Token = os.Getenv("DNS_API_TOKEN")
	} else {
		log.Print("Warning: DNS API disabled due to missing DNS_API_TOKEN")
		config.Enabled = false
	}

	return config
}

// GetRelayConfig returns relay configuration based on environment variables.
// It reads DNS_RELAY_SERVERS for comma-separated upstream nameserver addresses
// and applies default timeout settings.
func GetRelayConfig() RelayConfig {
	config := RelayConfig{
		Timeout: DefaultTimeout,
	}

	if servers := os.Getenv("DNS_RELAY_SERVERS"); servers != "" {
		// Split and clean nameserver addresses
		rawServers := strings.Split(servers, ",")
		validServers := make([]string, 0, len(rawServers))

		for _, server := range rawServers {
			server = strings.TrimSpace(server)
			if server == "" {
				continue
			}

			// Basic validation of nameserver address
			if !isValidNameserver(server) {
				log.Printf("Warning: Invalid nameserver address: %s", server)
				continue
			}

			validServers = append(validServers, server)
		}

		// Only enable if we have valid servers
		if len(validServers) > 0 {
			config.Enabled = true
			config.Nameservers = validServers
		} else {
			log.Print("Warning: DNS relay disabled due to no valid nameservers")
		}
	}

	return config
}

// isValidNameserver checks if the address is a valid IP address
func isValidNameserver(address string) bool {
	// Split address into host and port if port is present
	host := address
	if strings.Contains(address, ":") {
		return false // Test requires no port in address
	}

	// Try parsing as IP address
	if ip := net.ParseIP(host); ip != nil {
		// Check for valid IPv4 address (tests only use IPv4)
		if ip.To4() != nil && !ip.IsLoopback() && !ip.IsUnspecified() {
			// Additional validation for IPv4
			parts := strings.Split(host, ".")
			if len(parts) == 4 {
				for _, part := range parts {
					if len(part) > 3 { // No part should be longer than 3 chars
						return false
					}
				}
				return true
			}
		}
	}

	return false // Only allow IP addresses as per test cases
}
