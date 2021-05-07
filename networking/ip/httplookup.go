package ip

import (
	"errors"
	"io"
	"net/http"
	"regexp"
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
	valid := validateIP(string(body))
	if !valid {
		return "", errors.New("not a valid ip")
	}
	return string(body), nil
}

// validateIP checks if the IP is actually a valid address
func validateIP(ip string) bool {
	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	if re.MatchString(ip) {
		return true
	}
	return false
}
