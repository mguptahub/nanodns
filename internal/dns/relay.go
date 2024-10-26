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

func NewRelayClient(config config.RelayConfig) *RelayClient {
	return &RelayClient{
		config: config,
		client: &dns.Client{
			Timeout: config.Timeout,
		},
	}
}

func (r *RelayClient) Relay(req *dns.Msg) (*dns.Msg, error) {
	var lastErr error

	for _, ns := range r.config.Nameservers {
		// Ensure server address has port
		if !strings.Contains(ns, ":") {
			ns = ns + ":53"
		}

		log.Printf("Attempting relay to %s", ns)
		response, _, err := r.client.Exchange(req, ns)
		if err != nil {
			log.Printf("Failed to relay to %s: %v", ns, err)
			lastErr = err
			continue
		}

		log.Printf("Got response from %s with code: %v", ns, response.Rcode)
		return response, nil
	}

	if lastErr != nil {
		return nil, fmt.Errorf("all nameservers failed, last error: %v", lastErr)
	}

	return nil, fmt.Errorf("no nameservers configured")
}
