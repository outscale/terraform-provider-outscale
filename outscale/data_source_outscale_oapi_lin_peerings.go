package outscale

import (
	"fmt"
	"log"
	"time"

	"github.com/outscale/osc-go/oapi"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
						"tags": tagsOAPIListSchemaComputed(),
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
	conn := meta.(*OutscaleClient).OAPI

	log.Printf("[DEBUG] Reading VPC Peering Connections.")

	filters, filtersOk := d.GetOk("filter")
	if !filtersOk {
		return fmt.Errorf("One of filters must be assigned")
	}

	params := oapi.ReadNetPeeringsRequest{
		Filters: buildOutscaleOAPILinPeeringConnectionFilters(filters.(*schema.Set)),
	}

	var resp *oapi.POST_ReadNetPeeringsResponses
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.POST_ReadNetPeerings(params)
		return resource.RetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error reading the Net Peerings %s", err)
	}

	if resp.OK.NetPeerings == nil || len(resp.OK.NetPeerings) == 0 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	peerings := resp.OK.NetPeerings

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(resource.UniqueId())

		if err := set("net_peerings", getOAPINetPeerings(peerings)); err != nil {
			log.Printf("[DEBUG] Net Peerings ERR %+v", err)
			return err
		}
		return d.Set("request_id", resp.OK.ResponseContext.RequestId)
	})
}

func getOAPINetPeerings(peerings []oapi.NetPeering) (res []map[string]interface{}) {
	for _, p := range peerings {
		res = append(res, map[string]interface{}{
			"accepter_net":   getOAPINetPeeringAccepterNet(p.AccepterNet),
			"net_peering_id": p.NetPeeringId,
			"source_net":     getOAPINetPeeringSourceNet(p.SourceNet),
			"state":          getOAPINetPeeringState(p.State),
			"tags":           getOapiTagSet(p.Tags),
		})
	}
	return res
}

func getOAPINetPeeringAccepterNet(a oapi.AccepterNet) map[string]interface{} {
	return map[string]interface{}{
		"ip_range":   a.IpRange,
		"account_id": a.AccountId,
		"net_id":     a.NetId,
	}
}

func getOAPINetPeeringSourceNet(a oapi.SourceNet) map[string]interface{} {
	return map[string]interface{}{
		"ip_range":   a.IpRange,
		"account_id": a.AccountId,
		"net_id":     a.NetId,
	}
}

func getOAPINetPeeringState(a oapi.NetPeeringState) map[string]interface{} {
	return map[string]interface{}{
		"name":    a.Name,
		"message": a.Message,
	}
}
