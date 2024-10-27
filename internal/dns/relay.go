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
}

func (e *RelayError) Error() string {
	return fmt.Sprintf("relay to %s failed: %v", e.Server, e.Err)
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
	var lastErr error

	for _, ns := range r.config.Nameservers {
		// Ensure server address has port
		if !strings.Contains(ns, ":") {
			ns = ns + ":" + defaultDNSPort
		}

		// log.Printf("Attempting relay to %s", ns)
		log.Printf("relay_attempt: server=%s, query=%s", ns, req.Question[0].Name)
		response, _, err := r.client.Exchange(req, ns)
		if err != nil {
			log.Printf("Failed to relay to %s: %v", ns, err)
			lastErr = err
			continue
		}

		log.Printf("Got response from %s with code: %v", ns, response.Rcode)
		return response, nil
	}

	// if lastErr != nil {
	// 	return nil, fmt.Errorf("all nameservers failed, last error: %v", lastErr)
	// }

	// return nil, fmt.Errorf("no nameservers configured")
	return nil, fmt.Errorf("all nameservers failed, last error: %v", lastErr)
}
