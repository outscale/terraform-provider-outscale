package outscale

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/spf13/cast"
)

func DataSourceOutscaleSubregions() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleSubregionsRead,

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

func DataSourceOutscaleSubregionsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")

	var err error
	filtersReq := &oscgo.FiltersSubregion{}
	if filtersOk {
		filtersReq, err = buildOutscaleDataSourceSubregionsFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	req := oscgo.ReadSubregionsRequest{Filters: filtersReq}

	var resp oscgo.ReadSubregionsResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.SubregionApi.ReadSubregions(context.Background()).ReadSubregionsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	subregions := resp.GetSubregions()

	return resourceDataAttrSetter(d, func(set AttributeSetter) error {
		d.SetId(id.UniqueId())

		subs := make([]map[string]interface{}, len(subregions))
		for i, subregion := range subregions {
			subs[i] = map[string]interface{}{
				"location_code":  subregion.GetLocationCode(),
				"subregion_name": subregion.GetSubregionName(),
				"region_name":    subregion.GetRegionName(),
				"state":          subregion.GetState(),
			}
		}

		return set("subregions", subs)
	})
}

func buildOutscaleDataSourceSubregionsFilters(set *schema.Set) (*oscgo.FiltersSubregion, error) {
	filters := &oscgo.FiltersSubregion{}
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string

		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, cast.ToString(e))
		}

		switch name := m["name"].(string); name {
		case "subregion_names":
			filters.SetSubregionNames(filterValues)
		case "states":
			filters.SetStates(filterValues)
		case "region_names":
			filters.SetRegionNames(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return filters, nil
}
