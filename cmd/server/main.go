package main

import (
	"log"

	"github.com/mguptahub/nanodns/internal/dns"
	"github.com/mguptahub/nanodns/pkg/config"
	externaldns "github.com/miekg/dns"
)

func main() {
	// Load records from environment variables
	records := dns.LoadRecords()

	// Create DNS handler
	handler := dns.NewHandler(records)
	externaldns.HandleFunc(".", handler.ServeDNS)

	// Configure server
	port := config.GetDNSPort()
	server := &externaldns.Server{
		Addr: ":" + port,
		Net:  "udp",
	}

	log.Printf("Starting DNS server on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	defer server.Shutdown()
}