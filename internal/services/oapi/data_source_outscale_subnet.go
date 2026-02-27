package oapi

import (
	"context"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleSubnet() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleSubnetRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"subregion_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": TagsSchemaComputedSDK(),
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"available_ips_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"map_public_ip_on_launch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func DataSourceOutscaleSubnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	req := osc.ReadSubnetsRequest{}

	if id := d.Get("subnet_id"); id != "" {
		req.Filters = &osc.FiltersSubnet{SubnetIds: &[]string{id.(string)}}
	}

	filters, filtersOk := d.GetOk("filter")

	var err error
	if filtersOk {
		req.Filters, err = buildOutscaleSubnetDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadSubnets(ctx, req, options.WithRetryTimeout(120*time.Second))
	if err != nil {
		errString := err.Error()
		return diag.Errorf("error reading subnet (%s)", errString)
	}

	if resp.Subnets == nil || len(*resp.Subnets) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	if len(*resp.Subnets) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	subnet := (*resp.Subnets)[0]

	d.SetId(subnet.SubnetId)
	if err := d.Set("subnet_id", subnet.SubnetId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("net_id", subnet.NetId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("subregion_name", subnet.SubregionName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ip_range", subnet.IpRange); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", subnet.State); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("map_public_ip_on_launch", subnet.MapPublicIpOnLaunch); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(subnet.Tags)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("available_ips_count", subnet.AvailableIpsCount); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildOutscaleSubnetDataSourceFilters(set *schema.Set) (*osc.FiltersSubnet, error) {
	var filters osc.FiltersSubnet
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "available_ips_counts":
			filters.AvailableIpsCounts = new(utils.StringSliceToIntSlice(filterValues))
		case "ip_ranges":
			filters.IpRanges = &filterValues
		case "net_ids":
			filters.NetIds = &filterValues
		case "states":
			filters.States = new(lo.Map(filterValues, func(s string, _ int) osc.SubnetState { return osc.SubnetState(s) }))
		case "subnet_ids":
			filters.SubnetIds = &filterValues
		case "subregion_names":
			filters.SubregionNames = &filterValues
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
