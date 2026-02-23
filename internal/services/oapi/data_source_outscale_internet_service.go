package oapi

import (
	"context"
	"log"
	"time"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleInternetService() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleInternetServiceRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": TagsSchemaComputedSDK(),
		},
	}
}

func DataSourceOutscaleInternetServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	internetID, insternetIDOk := d.GetOk("internet_service_id")

	if !filtersOk && !insternetIDOk {
		return diag.Errorf("one of filters, or instance_id must be assigned")
	}

	// Build up search parameters
	var err error
	params := osc.ReadInternetServicesRequest{}

	if insternetIDOk {
		params.Filters = &osc.FiltersInternetService{
			InternetServiceIds: &[]string{internetID.(string)},
		}
	}

	if filtersOk {
		params.Filters, err = buildOutscaleOSCAPIDataSourceInternetServiceFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadInternetServices(ctx, params, options.WithRetryTimeout(120*time.Second))
	if err != nil {
		return diag.Errorf("error reading internet service id (%s)", err)
	}

	if resp.InternetServices == nil || len(*resp.InternetServices) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	result := (*resp.InternetServices)[0]

	log.Printf("[DEBUG] Setting OAPI Internet Service id (%s)", err)

	if err := d.Set("internet_service_id", result.InternetServiceId); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("state", result.State); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("net_id", result.NetId); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.InternetServiceId)

	return diag.FromErr(d.Set("tags", FlattenOAPITagsSDK(result.Tags)))
}

func buildOutscaleOSCAPIDataSourceInternetServiceFilters(set *schema.Set) (*osc.FiltersInternetService, error) {
	var filters osc.FiltersInternetService
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "internet_service_ids":
			filters.InternetServiceIds = &filterValues
		case "link_net_ids":
			filters.LinkNetIds = &filterValues
		case "link_states":
			filters.LinkStates = &filterValues
		case "tags":
			filters.Tags = &filterValues
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
