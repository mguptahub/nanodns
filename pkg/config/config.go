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

func GetRelayConfig() RelayConfig {
	config := RelayConfig{
		Timeout: DefaultTimeout,
	}

	if servers := os.Getenv("DNS_RELAY_SERVERS"); servers != "" {
		config.Enabled = true
		config.Nameservers = strings.Split(servers, ",")
	}

	return config
}
