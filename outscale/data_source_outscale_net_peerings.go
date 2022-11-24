package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceLinPeeringsConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLinPeeringsConnectionRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"net_peerings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accepter_net": vpcPeeringConnectionOptionsSchema(),
						"net_peering_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_net": vpcPeeringConnectionOptionsSchema(),
						"state": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
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
						"tags": dataSourceTagsSchema(),
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

func dataSourceLinPeeringsConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*Client).OSCAPI

	log.Printf("[DEBUG] Reading VPC Peering Connections.")

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("One of filters must be assigned")
	}

	params := oscgo.ReadNetPeeringsRequest{}
	params.SetFilters(buildLinPeeringConnectionFilters(filters.(*schema.Set)))

	var resp oscgo.ReadNetPeeringsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.NetPeeringApi.ReadNetPeerings(context.Background()).ReadNetPeeringsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp.StatusCode, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error reading the Net Peerings %s", err)
	}

	if resp.GetNetPeerings() == nil || len(resp.GetNetPeerings()) == 0 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	peerings := resp.GetNetPeerings()

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(resource.UniqueId())

		if err := set("net_peerings", getNetPeerings(peerings)); err != nil {
			log.Printf("[DEBUG] Net Peerings ERR %+v", err)
			return err
		}
		return nil
	})
}

func getNetPeerings(peerings []oscgo.NetPeering) (res []map[string]interface{}) {
	for _, p := range peerings {
		res = append(res, map[string]interface{}{
			"accepter_net":   getNetPeeringAccepterNet(p.GetAccepterNet()),
			"net_peering_id": p.GetNetPeeringId(),
			"source_net":     getNetPeeringSourceNet(p.GetSourceNet()),
			"state":          getNetPeeringState(p.GetState()),
			//"tags":           getTagSet(p.Tags),
		})
	}
	return res
}

func getNetPeeringAccepterNet(a oscgo.AccepterNet) map[string]interface{} {
	return map[string]interface{}{
		"ip_range":   a.GetIpRange(),
		"account_id": a.GetAccountId(),
		"net_id":     a.GetNetId(),
	}
}

func getNetPeeringSourceNet(a oscgo.SourceNet) map[string]interface{} {
	return map[string]interface{}{
		"ip_range":   a.GetIpRange(),
		"account_id": a.GetAccountId(),
		"net_id":     a.GetNetId(),
	}
}

func getNetPeeringState(a oscgo.NetPeeringState) map[string]interface{} {
	return map[string]interface{}{
		"name":    a.Name,
		"message": a.Message,
	}
}
