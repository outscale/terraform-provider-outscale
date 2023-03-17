package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOutscaleOAPIVpcAttr() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVpcAttrRead,

		Schema: map[string]*schema.Schema{
			//"filter": dataSourceFiltersSchema(),
			"dhcp_options_set_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tenancy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": dataSourceTagsSchema(),
		},
	}
}

func dataSourceOutscaleOAPIVpcAttrRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters := oscgo.FiltersNet{
		NetIds: &[]string{d.Get("net_id").(string)},
	}

	req := oscgo.ReadNetsRequest{
		Filters: &filters,
	}

	var resp oscgo.ReadNetsResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.NetApi.ReadNets(context.Background()).ReadNetsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		log.Printf("[DEBUG] Error reading lin (%s)", utils.GetErrorResponse(err))
	}

	if len(resp.GetNets()) == 0 {
		d.SetId("")
		return fmt.Errorf("oAPI Net not found")
	}

	d.SetId(resp.GetNets()[0].GetNetId())

	d.Set("ip_range", resp.GetNets()[0].GetIpRange())
	d.Set("tenancy", resp.GetNets()[0].Tenancy)
	d.Set("dhcp_options_set_id", resp.GetNets()[0].GetDhcpOptionsSetId())
	d.Set("net_id", resp.GetNets()[0].GetNetId())
	d.Set("state", resp.GetNets()[0].GetState())

	return d.Set("tags", tagsOSCAPIToMap(resp.GetNets()[0].GetTags()))
}
