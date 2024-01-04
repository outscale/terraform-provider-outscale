package outscale

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOutscaleOAPILinPeeringsConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPILinPeeringsConnectionRead,

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

func dataSourceOutscaleOAPILinPeeringsConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	log.Printf("[DEBUG] Reading VPC Peering Connections.")

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("One of filters must be assigned")
	}

	params := oscgo.ReadNetPeeringsRequest{}
	params.SetFilters(buildOutscaleOAPILinPeeringConnectionFilters(filters.(*schema.Set)))

	var resp oscgo.ReadNetPeeringsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.NetPeeringApi.ReadNetPeerings(context.Background()).ReadNetPeeringsRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error reading the Net Peerings %s", err)
	}
	peerings := resp.GetNetPeerings()

	if peerings == nil || len(peerings) == 0 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}
	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(resource.UniqueId())

		if err := set("net_peerings", setNetPeeringsAttributtes(peerings)); err != nil {
			return err
		}
		return nil
	})
}

func setNetPeeringsAttributtes(peerings []oscgo.NetPeering) (res []map[string]interface{}) {

	for _, p := range peerings {
		netP := map[string]interface{}{
			"net_peering_id": p.GetNetPeeringId(),
		}
		if p.HasAccepterNet() {
			if !reflect.DeepEqual(p.GetAccepterNet(), oscgo.AccepterNet{}) {
				netP["accepter_net"] = getOAPINetPeeringAccepterNet(p.GetAccepterNet())
			}
		}
		if p.HasSourceNet() {
			if !reflect.DeepEqual(p.GetSourceNet(), oscgo.SourceNet{}) {
				netP["source_net"] = getOAPINetPeeringSourceNet(p.GetSourceNet())
			}
		}
		if p.HasState() {
			netP["state"] = getOAPINetPeeringState(p.GetState())
		}
		if p.HasTags() {
			netP["tags"] = getOapiTagSet(p.Tags)
		}
		res = append(res, netP)
	}
	return
}

func getOAPINetPeeringAccepterNet(a oscgo.AccepterNet) []map[string]interface{} {
	return []map[string]interface{}{{
		"ip_range":   a.GetIpRange(),
		"account_id": a.GetAccountId(),
		"net_id":     a.GetNetId(),
	}}
}

func getOAPINetPeeringSourceNet(a oscgo.SourceNet) []map[string]interface{} {
	return []map[string]interface{}{{
		"ip_range":   a.GetIpRange(),
		"account_id": a.GetAccountId(),
		"net_id":     a.GetNetId(),
	}}
}

func getOAPINetPeeringState(a oscgo.NetPeeringState) []map[string]interface{} {
	return []map[string]interface{}{{
		"name":    a.GetName(),
		"message": a.GetMessage(),
	}}
}
