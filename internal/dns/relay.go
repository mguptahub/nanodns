package dns

import (
	"fmt"
	"log"
	"strings"

	"github.com/mguptahub/nanodns/pkg/config"
	"github.com/miekg/dns"
)

type RelayClient struct {
	config config.RelayConfig
	client *dns.Client
}

type RelayError struct {
	Server string
	Err    error
	Query  string // DNS query that failed
	Rcode  int    // DNS response code if available
}

func (e *RelayError) Error() string {
	if e.Rcode != 0 {
		return fmt.Sprintf("relay to %s failed for query %s: %v (rcode: %d)",
			e.Server, e.Query, e.Err, e.Rcode)
	}
	return fmt.Sprintf("relay to %s failed for query %s: %v",
		e.Server, e.Query, e.Err)
}

// NewRelayClient creates a new RelayClient with the provided configuration.
// It returns an error if the configuration is invalid.
func NewRelayClient(config config.RelayConfig) (*RelayClient, error) {
	if config.Timeout <= 0 {
		return nil, fmt.Errorf("timeout must be positive")
	}
	if len(config.Nameservers) == 0 {
		return nil, fmt.Errorf("at least one nameserver must be configured")
	}
	return &RelayClient{
		config: config,
		client: &dns.Client{
			Timeout: config.Timeout,
		},
	}, nil
}

const defaultDNSPort = "53"

// Relay forwards the DNS request to configured upstream nameservers.
// It attempts each nameserver in sequence until a successful response is received.
// Returns the first successful response or an error if all nameservers fail.
func (r *RelayClient) Relay(req *dns.Msg) (*dns.Msg, error) {
	if len(req.Question) == 0 {
		return nil, fmt.Errorf("empty question in DNS request")
	}
	var lastErr error

	for _, ns := range r.config.Nameservers {
		// Ensure server address has port
		if !strings.Contains(ns, ":") {
			ns = ns + ":" + defaultDNSPort
		}

		log.Printf("relay_attempt: server=%s, query=%s", ns, req.Question[0].Name)
		response, rtt, err := r.client.Exchange(req, ns)
		if err != nil {
			log.Printf("relay_failed: server=%s, query=%s, error=%v", ns, req.Question[0].Name, err)
			lastErr = &RelayError{
				Server: ns,
				Err:    err,
				Query:  req.Question[0].Name,
			}
			continue
		}

		log.Printf("relay_success: server=%s, query=%s, rcode=%v, rtt=%v", ns, req.Question[0].Name, response.Rcode, rtt)
		return response, nil
	}

	return nil, fmt.Errorf("all nameservers failed, last error: %v", lastErr)
}
