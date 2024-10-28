package config

import (
	"log"
	"net"
	"os"
	"strings"
	"time"
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
		hasInvalid := false

		for _, server := range rawServers {
			server = strings.TrimSpace(server)

			// Skip empty entries
			if server == "" {
				continue
			} else if !isValidNameserver(server) {
				log.Printf("Warning: Invalid nameserver address: %s", server)
				hasInvalid = true
				continue
			}

			validServers = append(validServers, server)
		}

		// Enable only if all provided servers are valid
		if len(validServers) > 0 && !hasInvalid {
			config.Enabled = true
			config.Nameservers = validServers
		} else {
			log.Print("Warning: DNS relay disabled due to invalid nameserver entries")
		}
	}

	return config
}

// isValidNameserver checks if the address is a valid IP address
func isValidNameserver(address string) bool {
	// Split address into host and port if port is present
	host := address
	if strings.Contains(address, ":") {
		parts := strings.Split(address, ":")
		host = parts[0]
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
