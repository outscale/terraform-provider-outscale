package oapi

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleNetPeerings() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleNetPeeringsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"net_peerings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accepter_net": vpcOAPIPeeringConnectionOptionsSchema(),
						"net_peering_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_net": vpcOAPIPeeringConnectionOptionsSchema(),
						"state": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"message": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"tags": TagsSchemaComputedSDK(),
					},
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func DataSourceOutscaleNetPeeringsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	log.Printf("[DEBUG] Reading VPC Peering connections.")

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return diag.FromErr(ErrFilterRequired)
	}

	var err error
	params := osc.ReadNetPeeringsRequest{}
	params.Filters, err = buildOutscaleLinPeeringConnectionFilters(filters.(*schema.Set))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ReadNetPeerings(ctx, params, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.Errorf("error reading the net peerings %s", err)
	}
	peerings := resp.NetPeerings

	if peerings == nil || len(*peerings) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	return diag.FromErr(resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(id.UniqueId())

		if err := set("net_peerings", setNetPeeringsAttributtes(*peerings)); err != nil {
			return err
		}
		return nil
	}))
}

func setNetPeeringsAttributtes(peerings []osc.NetPeering) (res []map[string]interface{}) {
	for _, p := range peerings {
		netP := map[string]interface{}{
			"net_peering_id": p.NetPeeringId,
		}
		if !reflect.DeepEqual(p.AccepterNet, osc.AccepterNet{}) {
			netP["accepter_net"] = getOAPINetPeeringAccepterNet(p.AccepterNet)
		}
		if !reflect.DeepEqual(p.SourceNet, osc.SourceNet{}) {
			netP["source_net"] = getOAPINetPeeringSourceNet(p.SourceNet)
		}
		netP["state"] = getOAPINetPeeringState(p.State)

		if p.Tags != nil {
			netP["tags"] = FlattenOAPITagsSDK(p.Tags)
		}
		res = append(res, netP)
	}
	return
}

func getOAPINetPeeringAccepterNet(a osc.AccepterNet) []map[string]interface{} {
	return []map[string]interface{}{{
		"ip_range":   a.IpRange,
		"account_id": a.AccountId,
		"net_id":     a.NetId,
	}}
}

func getOAPINetPeeringSourceNet(a osc.SourceNet) []map[string]interface{} {
	return []map[string]interface{}{{
		"ip_range":   a.IpRange,
		"account_id": a.AccountId,
		"net_id":     a.NetId,
	}}
}

func getOAPINetPeeringState(a osc.NetPeeringState) []map[string]interface{} {
	return []map[string]interface{}{{
		"name":    a.Name,
		"message": a.Message,
	}}
}
