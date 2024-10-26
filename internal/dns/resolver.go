package dns

import (
	"fmt"
	"net"
)

// ResolveServiceIP attempts to resolve Docker service name to IP
func ResolveServiceIP(serviceName string) (string, error) {
	ips, err := net.LookupIP(serviceName)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String(), nil
		}
	}

	return "", fmt.Errorf("no IPv4 address found for service: %s", serviceName)
}
