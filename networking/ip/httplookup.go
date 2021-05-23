package ip

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// GetIP retrieves IP via web request
func (h *HttpLookup) GetIP() (string, error) {
	resp, err := http.Get(h.Url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	ip := strings.TrimRight(string(body), "\n")
	valid := validateIP(ip)
	if !valid {
		return "", errors.New("not a valid ip")
	}
	return ip, nil
}

// validateIP checks if the IP is actually a valid address
func validateIP(ip string) bool {
	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	if re.MatchString(ip) {
		return true
	}
	return false
}
