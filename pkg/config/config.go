package config

import (
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
		config.Enabled = true
		// Split and clean nameserver addresses
		rawServers := strings.Split(servers, ",")
		config.Nameservers = make([]string, 0, len(rawServers))
		for _, server := range rawServers {
			server = strings.TrimSpace(server)
			if server == "" {
				continue
			}
			config.Nameservers = append(config.Nameservers, server)
		}

		// Disable relay if no valid nameservers
		if len(config.Nameservers) == 0 {
			config.Enabled = false
		}
	}

	return config
}
