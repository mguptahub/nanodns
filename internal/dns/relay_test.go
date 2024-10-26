package dns

import (
	"testing"
	"time"

	"github.com/mguptahub/nanodns/pkg/config"
	"github.com/miekg/dns"
)

func TestRelayClient_Relay(t *testing.T) {
	tests := []struct {
		name        string
		config      config.RelayConfig
		query       string
		qtype       uint16
		shouldError bool
	}{
		{
			name: "Valid nameserver",
			config: config.RelayConfig{
				Enabled:     true,
				Nameservers: []string{"8.8.8.8:53"},
				Timeout:     5 * time.Second,
			},
			query:       "google.com.",
			qtype:       dns.TypeA,
			shouldError: false,
		},
		{
			name: "Invalid nameserver",
			config: config.RelayConfig{
				Enabled:     true,
				Nameservers: []string{"invalid.nameserver:53"},
				Timeout:     1 * time.Second,
			},
			query:       "example.com.",
			qtype:       dns.TypeA,
			shouldError: true,
		},
		{
			name: "Multiple nameservers with first failing",
			config: config.RelayConfig{
				Enabled:     true,
				Nameservers: []string{"invalid.nameserver:53", "8.8.8.8:53"},
				Timeout:     1 * time.Second,
			},
			query:       "google.com.",
			qtype:       dns.TypeA,
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewRelayClient(tt.config)

			m := new(dns.Msg)
			m.SetQuestion(tt.query, tt.qtype)

			resp, err := client.Relay(m)

			if tt.shouldError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
				if resp == nil {
					t.Error("Expected response but got nil")
				}
			}
		})
	}
}
