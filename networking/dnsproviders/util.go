package dnsproviders

import (
	"net"
)

// IsCurrentEndpoint checks if the host's DNS record(s) match an IP provided. Do not want 'Not Found' errors returned
func IsCurrentEndpoint(host, address string) (bool, error) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		if err.(*net.DNSError).IsNotFound {
			return false, nil
		}
		return false, err
	}
	for _, addr := range addrs {
		if addr == address {
			return true, nil
		}
	}
	return false, nil
}
