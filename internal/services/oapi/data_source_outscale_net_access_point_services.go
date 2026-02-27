package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleNetAccessPointServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleNetAccessPointServicesRead,
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

func DataSourceOutscaleNetAccessPointServicesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")

	var err error
	req := osc.ReadNetAccessPointServicesRequest{}
	if filtersOk {
		req.Filters, err = buildOutscaleDataSourcesNAPSFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadNetAccessPointServices(ctx, req, options.WithRetryTimeout(20*time.Second))
	if err != nil {
		return diag.FromErr(err)
	}

	naps := ptr.From(resp.Services)[:]
	nap_ret := make([]map[string]interface{}, len(naps))

	for k, v := range naps {
		n := make(map[string]interface{})
		n["ip_ranges"] = v.IpRanges
		n["service_id"] = v.ServiceId
		n["service_name"] = v.ServiceName
		nap_ret[k] = n
	}

	if err := d.Set("services", nap_ret); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.UniqueId())

	return nil
}

func buildOutscaleDataSourcesNAPSFilters(set *schema.Set) (*osc.FiltersService, error) {
	var filters osc.FiltersService

	for _, v := range set.List() {
		m := v.(map[string]interface{})
		filterValues := make([]string, 0)
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "service_ids":
			filters.ServiceIds = &filterValues
		case "service_names":
			filters.ServiceNames = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
