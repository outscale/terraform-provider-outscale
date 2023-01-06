package outscale

import (
	"context"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOutscaleOAPINetAccessPointServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPINetAccessPointServicesRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_ranges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
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

func dataSourceOutscaleOAPINetAccessPointServicesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")

	filtersReq := oscgo.FiltersService{}
	if filtersOk {
		filtersReq = buildOutscaleDataSourcesNAPSFilters(filters.(*schema.Set))
	}
	req := oscgo.ReadNetAccessPointServicesRequest{Filters: &filtersReq}

	var resp oscgo.ReadNetAccessPointServicesResponse
	var err error

	err = resource.Retry(20*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.NetAccessPointApi.ReadNetAccessPointServices(
			context.Background()).
			ReadNetAccessPointServicesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}

	naps := resp.GetServices()[:]
	nap_ret := make([]map[string]interface{}, len(naps))

	for k, v := range naps {
		n := make(map[string]interface{})
		n["ip_ranges"] = v.GetIpRanges()
		n["service_id"] = v.GetServiceId()
		n["service_name"] = v.GetServiceName()
		nap_ret[k] = n
	}

	if err := d.Set("services", nap_ret); err != nil {
		return err
	}

	d.SetId(resource.UniqueId())

	return nil
}

func buildOutscaleDataSourcesNAPSFilters(set *schema.Set) oscgo.FiltersService {
	var filters oscgo.FiltersService

	for _, v := range set.List() {
		m := v.(map[string]interface{})
		filterValues := make([]string, 0)
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "service_ids":
			filters.SetServiceIds(filterValues)
		case "service_names":
			filters.SetServiceNames(filterValues)
		default:
			log.Printf("[Debug] Unknown net access point services Filter Name: %s. default", name)
		}
	}
	return filters
}
