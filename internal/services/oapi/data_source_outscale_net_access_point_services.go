package oapi

import (
	"context"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleNetAccessPointServices() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleNetAccessPointServicesRead,
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

func DataSourceOutscaleNetAccessPointServicesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")

	var err error
	req := oscgo.ReadNetAccessPointServicesRequest{}
	if filtersOk {
		req.Filters, err = buildOutscaleDataSourcesNAPSFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp oscgo.ReadNetAccessPointServicesResponse

	err = retry.Retry(20*time.Second, func() *retry.RetryError {
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

	d.SetId(id.UniqueId())

	return nil
}

func buildOutscaleDataSourcesNAPSFilters(set *schema.Set) (*oscgo.FiltersService, error) {
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
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
