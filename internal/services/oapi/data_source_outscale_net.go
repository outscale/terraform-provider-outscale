package oapi

import (
	"context"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleVpc() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleVpcRead,

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

func DataSourceOutscaleVpcRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	var err error
	req := osc.ReadNetsRequest{}

	if v, ok := d.GetOk("filter"); ok {
		req.Filters, err = buildOutscaleDataSourceNetFilters(v.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if id := d.Get("net_id"); id != "" {
		req.Filters.NetIds = &[]string{id.(string)}
	}

	resp, err := client.ReadNets(ctx, req, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.Nets == nil || len(*resp.Nets) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	if len(*resp.Nets) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	net := (*resp.Nets)[0]

	d.SetId(net.NetId)

	if err := d.Set("net_id", net.NetId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ip_range", net.IpRange); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dhcp_options_set_id", net.DhcpOptionsSetId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenancy", net.Tenancy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", net.State); err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(d.Set("tags", FlattenOAPITagsSDK(net.Tags)))
}

func buildOutscaleDataSourceNetFilters(set *schema.Set) (*osc.FiltersNet, error) {
	var filters osc.FiltersNet
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "dhcp_options_set_ids":
			filters.DhcpOptionsSetIds = &filterValues
		case "ip_ranges":
			filters.IpRanges = &filterValues
		case "net_ids":
			filters.NetIds = &filterValues
		case "states":
			filters.States = new(lo.Map(filterValues, func(s string, _ int) osc.NetState { return osc.NetState(s) }))
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
