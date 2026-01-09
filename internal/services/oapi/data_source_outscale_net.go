package oapi

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleVpc() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleVpcRead,

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

			"tags": TagsSchemaComputedSDK(),
		},
	}
}

func DataSourceOutscaleVpcRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	var err error
	req := oscgo.ReadNetsRequest{}

	if v, ok := d.GetOk("filter"); ok {
		req.Filters, err = buildOutscaleDataSourceNetFilters(v.(*schema.Set))
		if err != nil {
			return err
		}
	}

	if id := d.Get("net_id"); id != "" {
		req.Filters.SetNetIds([]string{id.(string)})
	}

	var resp oscgo.ReadNetsResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
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

	return d.Set("tags", FlattenOAPITagsSDK(net.GetTags()))
}

func buildOutscaleDataSourceNetFilters(set *schema.Set) (*oscgo.FiltersNet, error) {
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
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
