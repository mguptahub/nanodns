package dns

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/mguptahub/nanodns/pkg/config"
	"github.com/miekg/dns"
)

type Handler struct {
	records map[string][]DNSRecord
	relay   *RelayClient
}

func NewHandler(records map[string][]DNSRecord, relayConfig config.RelayConfig) (*Handler, error) {
	var relay *RelayClient
	if relayConfig.Enabled {
		// relay, _ = NewRelayClient(relayConfig)
		var err error
		relay, err = NewRelayClient(relayConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize relay client: %w", err)
		}
	}

	return &Handler{
		records: records,
		relay:   relay,
	}, nil
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	for _, q := range r.Question {
		log.Printf("Query for %s (type: %v)", q.Name, dns.TypeToString[q.Qtype])

		// Try local records first
		if recs, exists := h.records[q.Name]; exists {
			for _, rec := range recs {
				if answer := h.createAnswer(q, rec); answer != nil {
					m.Answer = append(m.Answer, answer)
				}
			}
		}

		// If no local records found and relay is enabled, try relay
		if len(m.Answer) == 0 && h.relay != nil {
			log.Printf("No local records found for %s, attempting relay", q.Name)

			// Create a new message for just this question
			relayReq := new(dns.Msg)
			relayReq.SetQuestion(q.Name, q.Qtype)
			relayReq.RecursionDesired = true

			relayResp, err := h.relay.Relay(relayReq)
			if err != nil {
				log.Printf("Relay failed: %v", err)
				continue
			}

			// Validate relay response
			if relayResp.Rcode != dns.RcodeSuccess {
				log.Printf("Relay returned non-success code: %v", dns.RcodeToString[relayResp.Rcode])
				continue
			}

			// Add answers from relay response
			m.Answer = append(m.Answer, relayResp.Answer...)
			m.Ns = append(m.Ns, relayResp.Ns...)
			m.Extra = append(m.Extra, relayResp.Extra...)

			// If we got answers from relay, we're not authoritative
			if len(relayResp.Answer) > 0 {
				m.Authoritative = false
			}
		}
	}

	w.WriteMsg(m)
}

func (h *Handler) createAnswer(q dns.Question, rec DNSRecord) dns.RR {
	switch {
	case q.Qtype == dns.TypeA && rec.RecordType == ARecord:
		return h.createARecord(q, rec)
	case q.Qtype == dns.TypeCNAME && rec.RecordType == CNAMERecord:
		return h.createCNAMERecord(q, rec)
	case q.Qtype == dns.TypeMX && rec.RecordType == MXRecord:
		return h.createMXRecord(q, rec)
	case q.Qtype == dns.TypeTXT && rec.RecordType == TXTRecord:
		return h.createTXTRecord(q, rec)
	default:
		return nil
	}
}

func (h *Handler) createARecord(q dns.Question, rec DNSRecord) dns.RR {
	var ip net.IP
	if rec.IsService {
		resolvedIP, err := ResolveServiceIP(rec.Value)
		if err != nil {
			log.Printf("Failed to resolve service %s: %v", rec.Value, err)
			return nil
		}
		ip = net.ParseIP(resolvedIP)
	} else {
		ip = net.ParseIP(rec.Value)
	}

	if ip == nil {
		log.Printf("Invalid IP address for %s", rec.Value)
		return nil
	}

	return &dns.A{
		Hdr: dns.RR_Header{
			Name:   q.Name,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    rec.TTL,
		},
		A: ip,
	}
}

func (h *Handler) createCNAMERecord(q dns.Question, rec DNSRecord) dns.RR {
	return &dns.CNAME{
		Hdr: dns.RR_Header{
			Name:   q.Name,
			Rrtype: dns.TypeCNAME,
			Class:  dns.ClassINET,
			Ttl:    rec.TTL,
		},
		Target: rec.Value,
	}
}

func (h *Handler) createMXRecord(q dns.Question, rec DNSRecord) dns.RR {
	return &dns.MX{
		Hdr: dns.RR_Header{
			Name:   q.Name,
			Rrtype: dns.TypeMX,
			Class:  dns.ClassINET,
			Ttl:    rec.TTL,
		},
		Preference: rec.Priority,
		Mx:         rec.Value,
	}
}

func (h *Handler) createTXTRecord(q dns.Question, rec DNSRecord) dns.RR {
	// Split TXT record by spaces if it contains multiple strings
	txtParts := strings.Split(rec.Value, " ")
	return &dns.TXT{
		Hdr: dns.RR_Header{
			Name:   q.Name,
			Rrtype: dns.TypeTXT,
			Class:  dns.ClassINET,
			Ttl:    rec.TTL,
		},
		Txt: txtParts,
	}
}
