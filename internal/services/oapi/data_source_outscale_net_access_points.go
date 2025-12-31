package oapi

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
		Read: DataSourceOutscaleNetAccessPointsRead,

		Schema: getDataSourceSchemas(napSchema()),
	}
}

func buildOutscaleDataSourcesNAPFilters(set *schema.Set) (*oscgo.FiltersNetAccessPoint, error) {
	filters := oscgo.FiltersNetAccessPoint{}

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
			filters.States = &filterValues
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		case "net_access_point_ids":
			filters.NetAccessPointIds = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}

func DataSourceOutscaleNetAccessPointsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	req := oscgo.ReadNetAccessPointsRequest{}
	var resp oscgo.ReadNetAccessPointsResponse
	var err error

	if filtersOk {
		req.Filters, err = buildOutscaleDataSourcesNAPFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}
	err = retry.Retry(30*time.Second, func() *retry.RetryError {
		rp, httpResp, err := conn.NetAccessPointApi.ReadNetAccessPoints(
			context.Background()).
			ReadNetAccessPointsRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}

	naps := resp.GetNetAccessPoints()[:]
	if len(naps) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	nap_ret := make([]map[string]interface{}, len(naps))
	for k, v := range naps {
		n := make(map[string]interface{})

		n["net_access_point_id"] = v.NetAccessPointId
		n["route_table_ids"] = utils.StringSlicePtrToInterfaceSlice(v.RouteTableIds)
		n["net_id"] = v.NetId
		n["service_name"] = v.ServiceName
		n["state"] = v.State
		n["tags"] = FlattenOAPITagsSDK(v.GetTags())
		nap_ret[k] = n
	}

	err = d.Set("net_access_points", nap_ret)
	if err != nil {
		return err
	}

	d.SetId(id.UniqueId())

	return nil
}
