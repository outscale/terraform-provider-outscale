package outscale

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOutscaleOAPIVpcs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVpcsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"net_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_range": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"dhcp_options_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"net_id": {
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
				},
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPIVpcsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadNetsRequest{}

	filters, filtersOk := d.GetOk("filter")
	netIds, netIdsOk := d.GetOk("net_id")

	if !filtersOk && !netIdsOk {
		return fmt.Errorf("filters or net_id(s) must be provided")
	}

	if filtersOk {
		req.SetFilters(buildOutscaleOAPIDataSourceNetFilters(filters.(*schema.Set)))
	}

	if netIdsOk {
		ids := make([]string, len(netIds.([]interface{})))

		for k, v := range netIds.([]interface{}) {
			ids[k] = v.(string)
		}
		var filters oscgo.FiltersNet
		filters.SetNetIds(ids)
		req.SetFilters(filters)

	}

	var err error
	var resp oscgo.ReadNetsResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.NetApi.ReadNets(context.Background()).ReadNetsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}
	if len(resp.GetNets()) == 0 {
		return fmt.Errorf("no matching VPC found")
	}

	d.SetId(resource.UniqueId())

	nets := make([]map[string]interface{}, len(resp.GetNets()))

	for i, v := range resp.GetNets() {
		net := make(map[string]interface{})

		net["net_id"] = v.GetNetId()
		net["ip_range"] = v.GetIpRange()
		net["dhcp_options_set_id"] = v.GetDhcpOptionsSetId()
		net["tenancy"] = v.GetTenancy()
		net["state"] = v.GetState()
		if v.Tags != nil {
			net["tags"] = tagsOSCAPIToMap(v.GetTags())
		}

		nets[i] = net
	}

	if err := d.Set("nets", nets); err != nil {
		return err
	}

	return nil
}
