package dns

import (
	"os"
	"testing"
)

func TestParseRecord(t *testing.T) {
	tests := []struct {
		name        string
		key         string
		value       string
		wantRecord  DNSRecord
		wantErr     bool
		errContains string
	}{
		{
			name:  "valid A record",
			key:   "A_REC1",
			value: "example.com|192.168.1.1|300",
			wantRecord: DNSRecord{
				Domain:     "example.com.",
				Value:      "192.168.1.1",
				TTL:        300,
				RecordType: ARecord,
				IsService:  false,
			},
			wantErr: false,
		},
		{
			name:  "valid A record with service",
			key:   "A_REC1",
			value: "example.com|service:webapp",
			wantRecord: DNSRecord{
				Domain:     "example.com.",
				Value:      "webapp",
				TTL:        60,
				RecordType: ARecord,
				IsService:  true,
			},
			wantErr: false,
		},
		{
			name:  "valid CNAME record",
			key:   "CNAME_REC1",
			value: "www.example.com|example.com|600",
			wantRecord: DNSRecord{
				Domain:     "www.example.com.",
				Value:      "example.com",
				TTL:        600,
				RecordType: CNAMERecord,
			},
			wantErr: false,
		},
		{
			name:  "valid MX record",
			key:   "MX_REC1",
			value: "example.com|10|mail.example.com|300",
			wantRecord: DNSRecord{
				Domain:     "example.com.",
				Value:      "mail.example.com",
				TTL:        300,
				RecordType: MXRecord,
				Priority:   10,
			},
			wantErr: false,
		},
		{
			name:  "valid TXT record",
			key:   "TXT_REC1",
			value: "example.com|v=spf1 include:_spf.example.com ~all|300",
			wantRecord: DNSRecord{
				Domain:     "example.com.",
				Value:      "v=spf1 include:_spf.example.com ~all",
				TTL:        300,
				RecordType: TXTRecord,
			},
			wantErr: false,
		},
		{
			name:        "invalid format",
			key:         "A_REC1",
			value:       "example.com",
			wantErr:     true,
			errContains: "invalid format",
		},
		{
			name:        "invalid MX priority",
			key:         "MX_REC1",
			value:       "example.com|invalid|mail.example.com",
			wantErr:     true,
			errContains: "invalid MX priority",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRecord(tt.key, tt.value)
			if tt.wantErr {
				if err == nil {
					t.Errorf("parseRecord() error = nil, want error containing %q", tt.errContains)
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("parseRecord() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("parseRecord() error = %v, want nil", err)
				return
			}
			if !recordsEqual(got, tt.wantRecord) {
				t.Errorf("parseRecord() = %v, want %v", got, tt.wantRecord)
			}
		})
	}
}

func TestLoadRecords(t *testing.T) {
	// Save current env and defer restore
	oldEnv := os.Environ()
	defer func() {
		os.Clearenv()
		for _, pair := range oldEnv {
			parts := splitEnv(pair)
			os.Setenv(parts[0], parts[1])
		}
	}()

	// Set up test environment
	os.Clearenv()
	testEnv := map[string]string{
		"A_REC1":     "app.example.com|192.168.1.1|300",
		"CNAME_REC1": "www.example.com|app.example.com|600",
		"MX_REC1":    "example.com|10|mail.example.com|300",
		"TXT_REC1":   "example.com|v=spf1 include:_spf.example.com ~all|300",
	}

	for k, v := range testEnv {
		os.Setenv(k, v)
	}

	// Run test
	got := LoadRecords()

	// Verify records
	tests := []struct {
		domain string
		count  int
	}{
		{"app.example.com.", 1},
		{"www.example.com.", 1},
		{"example.com.", 2}, // MX and TXT records
	}

	for _, tt := range tests {
		t.Run(tt.domain, func(t *testing.T) {
			records := got[tt.domain]
			if len(records) != tt.count {
				t.Errorf("LoadRecords() got %d records for %s, want %d", len(records), tt.domain, tt.count)
			}
		})
	}
}

// Helper functions
func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[:len(substr)] == substr
}

func recordsEqual(a, b DNSRecord) bool {
	return a.Domain == b.Domain &&
		a.Value == b.Value &&
		a.TTL == b.TTL &&
		a.RecordType == b.RecordType &&
		a.IsService == b.IsService &&
		a.Priority == b.Priority
}

func splitEnv(env string) [2]string {
	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			return [2]string{env[:i], env[i+1:]}
		}
	}
	return [2]string{env, ""}
}
