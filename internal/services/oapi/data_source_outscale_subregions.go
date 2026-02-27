package oapi

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"
)

func DataSourceOutscaleSubregions() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleSubregionsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			// Computed values.
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subregions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subregion_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceOutscaleSubregionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")

	var err error
	filtersReq := &osc.FiltersSubregion{}
	if filtersOk {
		filtersReq, err = buildOutscaleDataSourceSubregionsFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	req := osc.ReadSubregionsRequest{Filters: filtersReq}

	resp, err := client.ReadSubregions(ctx, req, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	subregions := ptr.From(resp.Subregions)

	return diag.FromErr(resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(id.UniqueId())

		subs := make([]map[string]interface{}, len(subregions))
		for i, subregion := range subregions {
			subs[i] = map[string]interface{}{
				"location_code":  subregion.LocationCode,
				"subregion_name": subregion.SubregionName,
				"region_name":    subregion.RegionName,
				"state":          subregion.State,
			}
		}

		return set("subregions", subs)
	}))
}

func buildOutscaleDataSourceSubregionsFilters(set *schema.Set) (*osc.FiltersSubregion, error) {
	filters := &osc.FiltersSubregion{}
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string

		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, cast.ToString(e))
		}

		switch name := m["name"].(string); name {
		case "subregion_names":
			filters.SubregionNames = &filterValues
		case "states":
			filters.States = &filterValues
		case "region_names":
			filters.RegionNames = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return filters, nil
}
