package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func napSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"net_access_point_ids": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"net_ids": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"service_names": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"states": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"tag_keys": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"tag_values": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"tags": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
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
					"tags": dataSourceTagsSchema(),
				},
			},
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func dataSourceOutscaleNetAccessPoints() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleNetAccessPointsRead,

		Schema: getDataSourceSchemas(napSchema()),
	}
}

func buildOutscaleDataSourcesNAPFilters(set *schema.Set) *oscgo.FiltersNetAccessPoint {
	filters := new(oscgo.FiltersNetAccessPoint)

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
		case "net_access_point_id":
			filters.NetAccessPointIds = &filterValues
		default:
			filters.NetAccessPointIds = &filterValues
			log.Printf("[Debug] Unknown Filter Name: %s. default to 'net_access_point_id'", name)
		}
	}
	return filters
}

func dataSourceOutscaleNetAccessPointsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	napid, napidOk := d.GetOk("net_access_point_ids")
	filters, filtersOk := d.GetOk("filter")
	filter := new(oscgo.FiltersNetAccessPoint)

	if !napidOk && !filtersOk {
		return fmt.Errorf("One of filters, or net_access_point_ids must be assigned")
	}

	if filtersOk {
		filter = buildOutscaleDataSourcesNAPFilters(filters.(*schema.Set))
	} else {
		filter = &oscgo.FiltersNetAccessPoint{
			NetAccessPointIds: &[]string{napid.(string)},
		}
	}

	req := &oscgo.ReadNetAccessPointsRequest{
		Filters: filter,
	}

	var resp oscgo.ReadNetAccessPointsResponse
	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		resp, _, err = conn.NetAccessPointApi.ReadNetAccessPoints(
			context.Background()).
			ReadNetAccessPointsRequest(*req).Execute()
		if err != nil {
			return utils.CheckThrottling(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	naps := *resp.NetAccessPoints
	nap_len := len(naps)
	nap_ret := make([]map[string]interface{}, nap_len)

	for k, v := range naps {
		n := make(map[string]interface{})

		n["net_access_point_id"] = v.NetAccessPointId
		n["route_table_ids"] = flattenStringList(v.RouteTableIds)
		n["net_id"] = v.NetId
		n["service_name"] = v.ServiceName
		n["state"] = v.State
		n["tags"] = tagsOSCAPIToMap(v.GetTags())
		nap_ret[k] = n
	}

	err = d.Set("net_access_point", nap_ret)
	if err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return nil
}
