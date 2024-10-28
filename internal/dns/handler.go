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
	// Normalize all record names to lowercase and ensure they're fully qualified
	normalizedRecords := make(map[string][]DNSRecord)
	for k, v := range records {
		normalizedKey := dns.CanonicalName(k)
		normalizedRecords[strings.ToLower(normalizedKey)] = make([]DNSRecord, len(v))
		for i, rec := range v {
			// Create a copy of the record
			newRec := rec
			// Ensure the value is fully qualified for CNAME records
			if rec.RecordType == CNAMERecord {
				newRec.Value = dns.CanonicalName(rec.Value)
			}
			normalizedRecords[strings.ToLower(normalizedKey)][i] = newRec
		}
	}

	var relay *RelayClient
	if relayConfig.Enabled {
		var err error
		relay, err = NewRelayClient(relayConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize relay client: %w", err)
		}
	}

	return &Handler{
		records: normalizedRecords,
		relay:   relay,
	}, nil
}

func (h *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Compress = true

	for _, q := range r.Question {
		log.Printf("Query for %s (type: %v)", q.Name, dns.TypeToString[q.Qtype])

		// Try to find matching records
		matchingRecords := h.findMatchingRecords(q.Name)
		log.Printf("Found %d matching records for %s", len(matchingRecords), q.Name)

		// Domain exists (found matching records)
		if len(matchingRecords) > 0 {
			answers := h.processRecords(q, matchingRecords)
			if len(answers) > 0 {
				m.Answer = append(m.Answer, answers...)
				log.Printf("Added %d answers for %s", len(answers), q.Name)
				continue // Skip relay if we have local answers
			}
			// Domain exists but no matching record type - return NOERROR with no answers
			m.Rcode = dns.RcodeSuccess
			continue
		}

		// Domain doesn't exist locally - try relay if enabled
		if h.relay != nil {
			log.Printf("No local records found for %s, attempting relay", q.Name)

			relayReq := new(dns.Msg)
			relayReq.SetQuestion(q.Name, q.Qtype)
			relayReq.RecursionDesired = true

			relayResp, err := h.relay.Relay(relayReq)
			if err != nil {
				log.Printf("Relay failed: %v", err)
				m.Rcode = dns.RcodeNameError // Return NXDOMAIN on relay failure
				continue
			}

			if relayResp.Rcode != dns.RcodeSuccess {
				log.Printf("Relay returned non-success code: %v", dns.RcodeToString[relayResp.Rcode])
				// Convert SERVFAIL to NXDOMAIN when appropriate
				if relayResp.Rcode == dns.RcodeServerFailure {
					m.Rcode = dns.RcodeNameError
				} else {
					m.Rcode = relayResp.Rcode
				}
				continue
			}

			m.Answer = append(m.Answer, relayResp.Answer...)
			m.Ns = append(m.Ns, relayResp.Ns...)
			m.Extra = append(m.Extra, relayResp.Extra...)

			if len(relayResp.Answer) > 0 {
				m.Authoritative = false
			}
		} else {
			// No relay and domain doesn't exist - return NXDOMAIN
			m.Rcode = dns.RcodeNameError
		}
	}

	if err := w.WriteMsg(m); err != nil {
		log.Printf("Error writing DNS response: %v", err)
	} else {
		log.Printf("Successfully wrote DNS response with %d answers", len(m.Answer))
	}
}

func (h *Handler) processRecords(q dns.Question, records []DNSRecord) []dns.RR {
	var answers []dns.RR

	for _, rec := range records {
		switch rec.RecordType {
		case CNAMERecord:
			// Always add CNAME record first
			if cname := h.createCNAMERecord(q, rec); cname != nil {
				answers = append(answers, cname)
				log.Printf("Added CNAME record: %v", cname)

				// If query was for A record and we have a CNAME, try to resolve the target
				if q.Qtype == dns.TypeA {
					// Ensure CNAME target is fully qualified
					target := dns.CanonicalName(rec.Value)
					// Look for A record matching CNAME target
					targetRecords := h.findMatchingRecords(target)
					for _, targetRec := range targetRecords {
						if targetRec.RecordType == ARecord {
							if a := h.createARecord(dns.Question{
								Name:   q.Name,
								Qtype:  q.Qtype,
								Qclass: q.Qclass,
							}, targetRec); a != nil {
								answers = append(answers, a)
								log.Printf("Added A record for CNAME target: %v", a)
							}
						}
					}
				}
			}
		case ARecord:
			// Only add A record if specifically queried for it
			if q.Qtype == dns.TypeA {
				if a := h.createARecord(q, rec); a != nil {
					answers = append(answers, a)
					log.Printf("Added A record: %v", a)
				}
			}
		case MXRecord:
			// Only add MX record if specifically queried for it
			if q.Qtype == dns.TypeMX {
				if mx := h.createMXRecord(q, rec); mx != nil {
					answers = append(answers, mx)
					log.Printf("Added MX record: %v", mx)

					// Optionally resolve the MX target's A record
					targetRecords := h.findMatchingRecords(rec.Value)
					for _, targetRec := range targetRecords {
						if targetRec.RecordType == ARecord {
							// Create an additional A record for the MX server
							if a := h.createARecord(dns.Question{
								Name:   rec.Value,
								Qtype:  dns.TypeA,
								Qclass: q.Qclass,
							}, targetRec); a != nil {
								answers = append(answers, a)
								log.Printf("Added A record for MX target: %v", a)
							}
						}
					}
				}
			}
		case TXTRecord:
			// Only add TXT record if specifically queried for it
			if q.Qtype == dns.TypeTXT {
				if txt := h.createTXTRecord(q, rec); txt != nil {
					answers = append(answers, txt)
					log.Printf("Added TXT record: %v", txt)
				}
			}
		}
	}

	return answers
}

// findMatchingRecords finds all records that match the query name, including wildcard matches
func (h *Handler) findMatchingRecords(queryName string) []DNSRecord {
	// Normalize query name to lowercase and ensure it's fully qualified
	queryName = dns.CanonicalName(queryName)
	queryName = strings.ToLower(queryName)
	log.Printf("Looking for matches for normalized query: %s", queryName)

	// First try exact match
	if recs, exists := h.records[queryName]; exists {
		log.Printf("Found exact match for %s", queryName)
		return recs
	}

	// Try wildcard matching
	labels := dns.SplitDomainName(queryName)
	if len(labels) <= 1 {
		return nil
	}

	// Try the wildcard pattern
	wildcardName := dns.CanonicalName("*." + strings.Join(labels[1:], "."))
	wildcardName = strings.ToLower(wildcardName)
	log.Printf("Trying wildcard pattern: %s", wildcardName)

	if recs, exists := h.records[wildcardName]; exists {
		log.Printf("Found wildcard match: %s", wildcardName)
		// For each matching wildcard record, create a concrete version
		concreteRecords := make([]DNSRecord, 0, len(recs))
		for _, rec := range recs {
			newRec := rec
			if rec.RecordType == CNAMERecord {
				// For CNAME records, replace the wildcard in the target if it exists
				if strings.HasPrefix(rec.Value, "*.") {
					// Replace the * with the actual subdomain and ensure it's fully qualified
					newValue := strings.Replace(rec.Value, "*", labels[0], 1)
					newRec.Value = dns.CanonicalName(newValue)
				} else {
					newRec.Value = dns.CanonicalName(rec.Value)
				}
			}
			concreteRecords = append(concreteRecords, newRec)
		}
		return concreteRecords
	}

	log.Printf("No wildcard match found for %s", wildcardName)
	return nil
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
	// Ensure target is fully qualified
	target := dns.CanonicalName(rec.Value)
	return &dns.CNAME{
		Hdr: dns.RR_Header{
			Name:   q.Name,
			Rrtype: dns.TypeCNAME,
			Class:  dns.ClassINET,
			Ttl:    rec.TTL,
		},
		Target: target,
	}
}

func (h *Handler) createMXRecord(q dns.Question, rec DNSRecord) dns.RR {
	target := dns.CanonicalName(rec.Value)
	return &dns.MX{
		Hdr: dns.RR_Header{
			Name:   q.Name,
			Rrtype: dns.TypeMX,
			Class:  dns.ClassINET,
			Ttl:    rec.TTL,
		},
		Preference: rec.Priority,
		Mx:         target,
	}
}

func (h *Handler) createTXTRecord(q dns.Question, rec DNSRecord) dns.RR {
	// Split TXT record by spaces if it contains multiple strings
	txtParts := strings.Split(rec.Value, " ")
	// Remove empty strings and trim spaces
	var cleanParts []string
	for _, part := range txtParts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			cleanParts = append(cleanParts, trimmed)
		}
	}

	return &dns.TXT{
		Hdr: dns.RR_Header{
			Name:   q.Name,
			Rrtype: dns.TypeTXT,
			Class:  dns.ClassINET,
			Ttl:    rec.TTL,
		},
		Txt: cleanParts,
	}
}
