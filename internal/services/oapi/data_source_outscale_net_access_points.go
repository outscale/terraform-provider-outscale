package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/samber/lo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func napSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"filter": dataSourceFiltersSchema(),
		"net_access_points": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"net_access_point_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"net_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"route_table_ids": {
						Type:     schema.TypeList,
						Computed: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"service_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"state": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tags": TagsSchemaComputedSDK(),
				},
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func DataSourceOutscaleNetAccessPoints() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleNetAccessPointsRead,

		Schema: getDataSourceSchemas(napSchema()),
	}
}

func buildOutscaleDataSourcesNAPFilters(set *schema.Set) (*osc.FiltersNetAccessPoint, error) {
	filters := osc.FiltersNetAccessPoint{}

	for _, v := range set.List() {
		m := v.(map[string]interface{})
		filterValues := make([]string, 0)
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "net_ids":
			filters.NetIds = &filterValues
		case "service_names":
			filters.ServiceNames = &filterValues
		case "states":
			filters.States = new(lo.Map(filterValues, func(s string, _ int) osc.NetAccessPointState { return osc.NetAccessPointState(s) }))
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		case "net_access_point_ids":
			filters.NetAccessPointIds = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}

func DataSourceOutscaleNetAccessPointsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	req := osc.ReadNetAccessPointsRequest{}
	var err error

	if filtersOk {
		req.Filters, err = buildOutscaleDataSourcesNAPFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	resp, err := client.ReadNetAccessPoints(ctx, req, options.WithRetryTimeout(30*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	naps := ptr.From(resp.NetAccessPoints)[:]
	if len(naps) < 1 {
		return diag.FromErr(ErrNoResults)
	}

	nap_ret := make([]map[string]interface{}, len(naps))
	for k, v := range naps {
		n := make(map[string]interface{})

		n["net_access_point_id"] = v.NetAccessPointId
		n["route_table_ids"] = utils.StringSlicePtrToInterfaceSlice(v.RouteTableIds)
		n["net_id"] = v.NetId
		n["service_name"] = v.ServiceName
		n["state"] = ptr.From(v.State)
		n["tags"] = FlattenOAPITagsSDK(ptr.From(v.Tags))
		nap_ret[k] = n
	}

	err = d.Set("net_access_points", nap_ret)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.UniqueId())

	return nil
}
