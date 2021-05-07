package dnsproviders

import (
	"context"
	"errors"
)

type Nameserver interface {
	EnsureRecordA(ctx context.Context, host, ipAddr string) error
}

func NewNameserver(nsType string, c Config) (Nameserver, error) {
	var resolver Nameserver
	switch {
	case nsType == "aws":
		resolver = &AWS{
			HostedZoneID: c.Aws.HostedZoneID,
			Ttl:          c.Aws.Ttl,
		}
	case nsType == "cloudflare":
		resolver = &Cloudflare{
			ZoneName: c.Cf.ZoneName,
			Ttl:      c.Cf.Ttl,
		}
	default:
		return nil, errors.New("invalid resolver type")
	}
	return resolver, nil
}

type Config struct {
	Aws AWS        `json:"aws,omitempty"`
	Cf  Cloudflare `json:"cloudflare,omitempty"`
}
