package dnsproviders

import (
	"context"
	"os"

	cf "github.com/cloudflare/cloudflare-go"
)

const (
	CFToken = "CF_API_TOKEN"
)

type Cloudflare struct {
	ZoneName string `json:"zoneName"`
	Ttl      int    `json:"ttl"`
}

func (c *Cloudflare) EnsureRecordA(ctx context.Context, host, ipAddr string) error {
	client, err := cf.NewWithAPIToken(os.Getenv(CFToken))
	if err != nil {
		return err
	}

	zoneID, err := client.ZoneIDByName(c.ZoneName)
	if err != nil {
		return err
	}
	records, err := client.DNSRecords(ctx, zoneID, cf.DNSRecord{
		Type:     "A",
		Name:     host,
		Content:  ipAddr,
		ZoneID:   zoneID,
		ZoneName: c.ZoneName,
	})
	if err != nil {
		return err
	}

	if len(records) == 0 {
		_, err := client.CreateDNSRecord(ctx, zoneID, cf.DNSRecord{
			Type:     "A",
			Name:     host,
			Content:  ipAddr,
			TTL:      c.Ttl,
			ZoneID:   zoneID,
			ZoneName: c.ZoneName,
		})
		return err
	}
	return nil
}
