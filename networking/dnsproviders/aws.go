package dnsproviders

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

var comment = "Updated via DDNS Controller"

type AWS struct {
	HostedZoneID string `json:"hostedZoneID"`
	Ttl          int64  `json:"ttl"`
}

func (a *AWS) EnsureRecordA(ctx context.Context, host, ipAddr string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	client := route53.NewFromConfig(cfg)
	_, err = client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: types.ChangeActionUpsert,
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: &host,
						Type: types.RRTypeA,
						ResourceRecords: []types.ResourceRecord{
							{
								Value: &ipAddr,
							},
						},
						TTL: &a.Ttl,
					},
				},
			},
			Comment: &comment,
		},
		HostedZoneId: &a.HostedZoneID,
	})
	if err != nil {
		return err
	}
	return nil
}
