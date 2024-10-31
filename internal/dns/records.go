package dns

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/mguptahub/nanodns/pkg/config"
)

type RecordType string

const (
	ARecord     RecordType = "A"
	CNAMERecord RecordType = "CNAME"
	MXRecord    RecordType = "MX"
	TXTRecord   RecordType = "TXT"

	// Record separator
	RecordSeparator = "|"
)

type DNSRecord struct {
	Domain     string
	Value      string
	TTL        uint32
	RecordType RecordType
	IsService  bool
	Priority   uint16 // For MX records
}

var records = make(map[string][]DNSRecord)

// LoadRecords loads DNS records from environment variables
func LoadRecords() map[string][]DNSRecord {
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		key := pair[0]
		value := pair[1]

		if strings.HasPrefix(key, "A_") ||
			strings.HasPrefix(key, "CNAME_") ||
			strings.HasPrefix(key, "MX_") ||
			strings.HasPrefix(key, "TXT_") {

			record, err := parseRecord(key, value)
			if err != nil {
				log.Printf("Error parsing record %s: %v", key, err)
				continue
			}
			domain := record.Domain
			records[domain] = append(records[domain], record)
		}
	}

	logLoadedRecords()
	return records
}

func parseRecord(key, value string) (DNSRecord, error) {
	parts := strings.Split(value, RecordSeparator)

	if len(parts) < 2 {
		return DNSRecord{}, fmt.Errorf("invalid format: expected parts separated by %s", RecordSeparator)
	}

	domain := parts[0]
	if !strings.HasSuffix(domain, ".") {
		domain = domain + "."
	}

	ttl := uint32(config.DefaultTTL)
	record := DNSRecord{
		Domain: domain,
		TTL:    ttl,
	}

	// Set record type and parse value based on prefix
	switch {
	case strings.HasPrefix(key, "A_"):
		record.RecordType = ARecord
		record.Value = parts[1]
		if config.IsServiceRecord(record.Value) {
			record.IsService = true
			record.Value = config.GetServiceName(record.Value)
		}
		if len(parts) > 2 {
			if parsedTTL, err := strconv.ParseUint(parts[2], 10, 32); err == nil {
				record.TTL = uint32(parsedTTL)
			}
		}

	case strings.HasPrefix(key, "CNAME_"):
		record.RecordType = CNAMERecord
		record.Value = parts[1]
		if len(parts) > 2 {
			if parsedTTL, err := strconv.ParseUint(parts[2], 10, 32); err == nil {
				record.TTL = uint32(parsedTTL)
			}
		}

	case strings.HasPrefix(key, "MX_"):
		record.RecordType = MXRecord
		if len(parts) < 3 {
			return DNSRecord{}, fmt.Errorf("MX record requires priority: domain|priority|value[|ttl]")
		}
		priority, err := strconv.ParseUint(parts[1], 10, 16)
		if err != nil {
			return DNSRecord{}, fmt.Errorf("invalid MX priority: %v", err)
		}
		record.Priority = uint16(priority)
		record.Value = parts[2]
		if len(parts) > 3 {
			if parsedTTL, err := strconv.ParseUint(parts[3], 10, 32); err == nil {
				record.TTL = uint32(parsedTTL)
			}
		}

	case strings.HasPrefix(key, "TXT_"):
		record.RecordType = TXTRecord
		record.Value = parts[1]
		if len(parts) > 2 {
			if parsedTTL, err := strconv.ParseUint(parts[2], 10, 32); err == nil {
				record.TTL = uint32(parsedTTL)
			}
		}
	}

	return record, nil
}

func logLoadedRecords() {
	log.Println("Loaded DNS Records")
	for domain, recs := range records {
		for _, rec := range recs {
			var extraInfo string
			switch rec.RecordType {
			case MXRecord:
				extraInfo = fmt.Sprintf(" Priority: %d", rec.Priority)
			case ARecord:
				if rec.IsService {
					extraInfo = " (Docker Service)"
				}
			}
			log.Printf("%s -> %s (TTL: %d, Type: %s%s)",
				domain, rec.Value, rec.TTL, rec.RecordType, extraInfo)
		}
	}
}
