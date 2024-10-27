package dns

import (
	"net"
	"testing"
	"time"

	"github.com/mguptahub/nanodns/pkg/config"
	"github.com/miekg/dns"
)

type mockResponseWriter struct {
	msgs []*dns.Msg
}

func (m *mockResponseWriter) LocalAddr() net.Addr         { return nil }
func (m *mockResponseWriter) RemoteAddr() net.Addr        { return nil }
func (m *mockResponseWriter) WriteMsg(msg *dns.Msg) error { m.msgs = append(m.msgs, msg); return nil }
func (m *mockResponseWriter) Write([]byte) (int, error)   { return 0, nil }
func (m *mockResponseWriter) Close() error                { return nil }
func (m *mockResponseWriter) TsigStatus() error           { return nil }
func (m *mockResponseWriter) TsigTimersOnly(bool)         {}
func (m *mockResponseWriter) Hijack()                     {}

func TestHandler_ServeDNS(t *testing.T) {
	// Test records
	records := map[string][]DNSRecord{
		"example.com.": {
			{
				Domain:     "example.com.",
				Value:      "192.168.1.1",
				TTL:        300,
				RecordType: ARecord,
			},
			{
				Domain:     "example.com.",
				Value:      "mail.example.com",
				TTL:        300,
				RecordType: MXRecord,
				Priority:   10,
			},
		},
		"www.example.com.": {
			{
				Domain:     "www.example.com.",
				Value:      "example.com",
				TTL:        300,
				RecordType: CNAMERecord,
			},
		},
	}

	// Create relay config for testing
	relayConfig := config.RelayConfig{
		Enabled:     true,
		Nameservers: []string{"8.8.8.8:53"},
		Timeout:     5 * time.Second,
	}

	handler := NewHandler(records, relayConfig)

	tests := []struct {
		name           string
		question       dns.Question
		expectedRcode  int
		expectedCount  int
		expectedType   uint16
		expectedAnswer string
		expectRelay    bool
	}{
		{
			name: "A record query - local",
			question: dns.Question{
				Name:   "example.com.",
				Qtype:  dns.TypeA,
				Qclass: dns.ClassINET,
			},
			expectedRcode:  dns.RcodeSuccess,
			expectedCount:  1,
			expectedType:   dns.TypeA,
			expectedAnswer: "192.168.1.1",
			expectRelay:    false,
		},
		{
			name: "CNAME record query - local",
			question: dns.Question{
				Name:   "www.example.com.",
				Qtype:  dns.TypeCNAME,
				Qclass: dns.ClassINET,
			},
			expectedRcode:  dns.RcodeSuccess,
			expectedCount:  1,
			expectedType:   dns.TypeCNAME,
			expectedAnswer: "example.com",
			expectRelay:    false,
		},
		{
			name: "MX record query - local",
			question: dns.Question{
				Name:   "example.com.",
				Qtype:  dns.TypeMX,
				Qclass: dns.ClassINET,
			},
			expectedRcode:  dns.RcodeSuccess,
			expectedCount:  1,
			expectedType:   dns.TypeMX,
			expectedAnswer: "mail.example.com",
			expectRelay:    false,
		},
		{
			name: "Non-existent domain",
			question: dns.Question{
				Name:   "nonexistent.com.",
				Qtype:  dns.TypeA,
				Qclass: dns.ClassINET,
			},
			expectedRcode: dns.RcodeSuccess,
			expectedCount: 0,
			expectRelay:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &mockResponseWriter{msgs: make([]*dns.Msg, 0)}
			r := new(dns.Msg)
			r.Question = []dns.Question{tt.question}

			handler.ServeDNS(w, r)

			if len(w.msgs) != 1 {
				t.Fatalf("Expected 1 message, got %d", len(w.msgs))
			}

			msg := w.msgs[0]
			if msg.Rcode != tt.expectedRcode {
				t.Errorf("Expected Rcode %d, got %d", tt.expectedRcode, msg.Rcode)
			}

			if len(msg.Answer) != tt.expectedCount {
				t.Errorf("Expected %d answers, got %d", tt.expectedCount, len(msg.Answer))
			}

			if tt.expectedCount > 0 {
				ans := msg.Answer[0]
				if ans.Header().Rrtype != tt.expectedType {
					t.Errorf("Expected type %d, got %d", tt.expectedType, ans.Header().Rrtype)
				}

				switch rr := ans.(type) {
				case *dns.A:
					if rr.A.String() != tt.expectedAnswer {
						t.Errorf("Expected A record %s, got %s", tt.expectedAnswer, rr.A.String())
					}
				case *dns.CNAME:
					if rr.Target != tt.expectedAnswer {
						t.Errorf("Expected CNAME record %s, got %s", tt.expectedAnswer, rr.Target)
					}
				case *dns.MX:
					if rr.Mx != tt.expectedAnswer {
						t.Errorf("Expected MX record %s, got %s", tt.expectedAnswer, rr.Mx)
					}
				}

				if !tt.expectRelay && !msg.Authoritative {
					t.Error("Expected message to be authoritative for local records")
				}
			}
		})
	}
}

// TestHandlerWithoutRelay tests the handler without relay configuration
func TestHandlerWithoutRelay(t *testing.T) {
	records := map[string][]DNSRecord{
		"example.com.": {
			{
				Domain:     "example.com.",
				Value:      "192.168.1.1",
				TTL:        300,
				RecordType: ARecord,
			},
			{
				Domain:     "example.com.",
				Value:      "mail.example.com.",
				TTL:        300,
				RecordType: MXRecord,
				Priority:   10,
			},
		},
	}

	// Create handler without relay
	relayConfig := config.RelayConfig{
		Enabled: false,
	}
	handler := NewHandler(records, relayConfig)

	// Test cases for different record types
	testCases := []struct {
		name     string
		qtype    uint16
		expected bool
	}{
		{"A Record", dns.TypeA, true},
		{"MX Record", dns.TypeMX, true},
		{"TXT Record", dns.TypeTXT, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := &mockResponseWriter{msgs: make([]*dns.Msg, 0)}
			r := new(dns.Msg)
			r.SetQuestion("example.com.", tc.qtype)

			handler.ServeDNS(w, r)

			if len(w.msgs) != 1 {
				t.Fatal("Expected response message")
			}

			msg := w.msgs[0]
			hasAnswer := len(msg.Answer) > 0
			if hasAnswer != tc.expected {
				t.Errorf("Expected answer presence: %v, got: %v", tc.expected, hasAnswer)
			}
			if !msg.Authoritative {
				t.Error("Expected message to be authoritative when relay is disabled")
			}
		})
	}
}
