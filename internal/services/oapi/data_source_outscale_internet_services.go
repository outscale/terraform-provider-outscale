package oapi

import (
	"context"
	"log"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleInternetServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleInternetServicesRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"internet_service_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"internet_services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_service_id": {
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
		},
	}
}

func DataSourceOutscaleInternetServicesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	internetID, internetIDOk := d.GetOk("internet_service_ids")

	if !filtersOk && !internetIDOk {
		return diag.Errorf("one of filters, or instance_id must be assigned")
	}

	// Build up search parameters
	params := osc.ReadInternetServicesRequest{
		Filters: &osc.FiltersInternetService{},
	}
	filter := osc.FiltersInternetService{}
	if internetIDOk {
		i := internetID.([]string)
		in := make([]string, len(i))
		copy(in, i)
		filter.InternetServiceIds = &in
		params.Filters = &filter
	}

	var err error
	if filtersOk {
		params.Filters, err = buildOutscaleOSCAPIDataSourceInternetServiceFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadInternetServices(ctx, params, options.WithRetryTimeout(120*time.Second))
	if err != nil {
		return diag.Errorf("error reading internet services (%s)", err)
	}

	log.Printf("[DEBUG] Setting OAPI LIN Internet Gateways id (%s)", err)

	d.SetId(id.UniqueId())

	result := ptr.From(resp.InternetServices)
	return diag.FromErr(internetServicesOAPIDescriptionAttributes(d, result))
}

func internetServicesOAPIDescriptionAttributes(d *schema.ResourceData, internetServices []osc.InternetService) error {
	i := make([]map[string]interface{}, len(internetServices))
	for k, v := range internetServices {
		im := make(map[string]interface{})
		if v.State != "" {
			im["state"] = v.State
		}

		if v.NetId != "" {
			im["net_id"] = v.NetId
		}
		if v.InternetServiceId != "" {
			im["internet_service_id"] = v.InternetServiceId
		}
		if v.Tags != nil {
			im["tags"] = FlattenOAPITagsSDK(v.Tags)
		}
		i[k] = im
	}

	return d.Set("internet_services", i)
}
