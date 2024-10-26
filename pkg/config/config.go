package config

import (
	"os"
	"strings"
)

const (
	DefaultTTL    = 60
	DefaultPort   = "53"
	ServicePrefix = "service:"
)

// GetDNSPort returns the configured DNS port or default
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
