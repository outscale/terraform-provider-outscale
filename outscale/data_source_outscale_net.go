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

func dataSourceOutscaleOAPIVpc() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVpcRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
				Optional: true,
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": dataSourceTagsSchema(),
		},
	}
}

func dataSourceOutscaleOAPIVpcRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.ReadNetsRequest{}

	if v, ok := d.GetOk("filter"); ok {
		req.SetFilters(buildOutscaleOAPIDataSourceNetFilters(v.(*schema.Set)))
	}

	if id := d.Get("net_id"); id != "" {
		req.Filters.SetNetIds([]string{id.(string)})
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
		return fmt.Errorf("No matching Net found")
	}
	if len(resp.GetNets()) > 1 {
		return fmt.Errorf("Multiple Nets matched; use additional constraints to reduce matches to a single Net")
	}

	net := resp.GetNets()[0]

	d.SetId(net.GetNetId())

	if err := d.Set("net_id", net.GetNetId()); err != nil {
		return err
	}

	if err := d.Set("ip_range", net.GetIpRange()); err != nil {
		return err
	}
	if err := d.Set("dhcp_options_set_id", net.GetDhcpOptionsSetId()); err != nil {
		return err
	}
	if err := d.Set("tenancy", net.GetTenancy()); err != nil {
		return err
	}
	if err := d.Set("state", net.GetState()); err != nil {
		return err
	}

	return d.Set("tags", tagsOSCAPIToMap(net.GetTags()))
}

func buildOutscaleOAPIDataSourceNetFilters(set *schema.Set) oscgo.FiltersNet {
	var filters oscgo.FiltersNet
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "dhcp_options_set_ids":
			filters.SetDhcpOptionsSetIds(filterValues)
		case "ip_ranges":
			filters.SetIpRanges(filterValues)
		case "net_ids":
			filters.SetNetIds(filterValues)
		case "states":
			filters.SetStates(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
