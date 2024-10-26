package config

import (
	"os"
	"testing"
)

func TestGetDNSPort(t *testing.T) {
	// Save current env and defer restore
	oldPort := os.Getenv("DNS_PORT")
	defer os.Setenv("DNS_PORT", oldPort)

	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "default port",
			envValue: "",
			want:     DefaultPort,
		},
		{
			name:     "custom port",
			envValue: "5353",
			want:     "5353",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv("DNS_PORT", tt.envValue)
			} else {
				os.Unsetenv("DNS_PORT")
			}

			if got := GetDNSPort(); got != tt.want {
				t.Errorf("GetDNSPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsServiceRecord(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "service prefix",
			value: "service:webapp",
			want:  true,
		},
		{
			name:  "no service prefix",
			value: "192.168.1.1",
			want:  false,
		},
		{
			name:  "empty string",
			value: "",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsServiceRecord(tt.value); got != tt.want {
				t.Errorf("IsServiceRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceName(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "with service prefix",
			value: "service:webapp",
			want:  "webapp",
		},
		{
			name:  "without service prefix",
			value: "webapp",
			want:  "webapp",
		},
		{
			name:  "empty string",
			value: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceName(tt.value); got != tt.want {
				t.Errorf("GetServiceName() = %v, want %v", got, tt.want)
			}
		})
	}
}
