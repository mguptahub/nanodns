package dns

import (
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

type Handler struct {
	records map[string][]DNSRecord
}

func NewHandler(records map[string][]DNSRecord) *Handler {
	return &Handler{
		records: records,
	}
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true

	for _, q := range r.Question {
		log.Printf("Query for %s (type: %v)", q.Name, dns.TypeToString[q.Qtype])

		if recs, exists := h.records[q.Name]; exists {
			for _, rec := range recs {
				if answer := h.createAnswer(q, rec); answer != nil {
					m.Answer = append(m.Answer, answer)
				}
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
