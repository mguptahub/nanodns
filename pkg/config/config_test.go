package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestEnvFile tests the .env file loading
func TestEnvFile(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()

	// Create test .env file
	envContent := `
DNS_PORT=15353
DNS_API_TOKEN=test-token
DNS_RELAY_SERVERS=8.8.8.8
`
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Test loading the env file
	if err := LoadEnvFile(envPath); err != nil {
		t.Fatalf("Failed to load test env file: %v", err)
	}

	// Test configurations
	tests := []struct {
		name     string
		getFunc  func() string
		expected string
	}{
		{
			name:     "DNS_PORT from env file",
			getFunc:  GetDNSPort,
			expected: "15353",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.getFunc(); got != tt.expected {
				t.Errorf("Got %v, want %v", got, tt.expected)
			}
		})
	}
}

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
			envValue: "15353",
			want:     "15353",
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

func TestGetRelayConfig(t *testing.T) {
	// Save current env and defer restore
	oldServers := os.Getenv("DNS_RELAY_SERVERS")
	defer os.Setenv("DNS_RELAY_SERVERS", oldServers)

	tests := []struct {
		name     string
		envValue string
		want     RelayConfig
	}{
		{
			name:     "no servers",
			envValue: "",
			want: RelayConfig{
				Enabled: false,
				Timeout: DefaultTimeout,
			},
		},
		{
			name:     "single server",
			envValue: "8.8.8.8",
			want: RelayConfig{
				Enabled:     true,
				Nameservers: []string{"8.8.8.8"},
				Timeout:     DefaultTimeout,
			},
		},
		{
			name:     "multiple servers",
			envValue: "8.8.8.8,1.1.1.1",
			want: RelayConfig{
				Enabled:     true,
				Nameservers: []string{"8.8.8.8", "1.1.1.1"},
				Timeout:     DefaultTimeout,
			},
		},
		{
			name:     "with whitespace",
			envValue: " 8.8.8.8 , 1.1.1.1 ",
			want: RelayConfig{
				Enabled:     true,
				Nameservers: []string{"8.8.8.8", "1.1.1.1"},
				Timeout:     DefaultTimeout,
			},
		},
		{
			name:     "empty entries",
			envValue: "8.8.8.8,,1.1.1.1,",
			want: RelayConfig{
				Enabled:     true,
				Nameservers: []string{"8.8.8.8", "1.1.1.1"},
				Timeout:     DefaultTimeout,
			},
		},
		{
			name:     "invalid ip address",
			envValue: "256.256.256.256",
			want: RelayConfig{
				Enabled: false,
				Timeout: DefaultTimeout,
			},
		},
		{
			name:     "malformed input",
			envValue: "8.8.8.8:53,not.an.ip",
			want: RelayConfig{
				Enabled: false,
				Timeout: DefaultTimeout,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DNS_RELAY_SERVERS", tt.envValue)
			got := GetRelayConfig()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRelayConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAPIConfig(t *testing.T) {
	// Save current env and defer restore
	oldToken := os.Getenv("DNS_API_TOKEN")
	defer func() {
		os.Setenv("DNS_API_TOKEN", oldToken)
	}()

	tests := []struct {
		name        string
		enabled     string
		token       string
		wantEnabled bool
		wantToken   string
	}{
		{
			name:        "disabled by default",
			enabled:     "",
			token:       "",
			wantEnabled: false,
			wantToken:   "",
		},
		{
			name:        "enabled with token",
			enabled:     "true",
			token:       "secret-token",
			wantEnabled: true,
			wantToken:   "secret-token",
		},
		{
			name:        "enabled without token",
			enabled:     "true",
			token:       "",
			wantEnabled: false,
			wantToken:   "",
		},
		{
			name:        "explicitly disabled",
			enabled:     "false",
			token:       "secret-token",
			wantEnabled: false,
			wantToken:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("DNS_API_TOKEN", tt.token)

			got := GetAPIConfig()
			if got.Enabled != tt.wantEnabled {
				t.Errorf("GetAPIConfig().Enabled = %v, want %v", got.Enabled, tt.wantEnabled)
			}
			if got.Token != tt.wantToken {
				t.Errorf("GetAPIConfig().Token = %v, want %v", got.Token, tt.wantToken)
			}
		})
	}
}
