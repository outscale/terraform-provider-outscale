package oapi

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleInternetService() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleInternetServiceRead,
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

func DataSourceOutscaleInternetServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	internetID, insternetIDOk := d.GetOk("internet_service_id")

	if !filtersOk && !insternetIDOk {
		return fmt.Errorf("One of filters, or instance_id must be assigned")
	}

	// Build up search parameters
	var err error
	params := oscgo.ReadInternetServicesRequest{}

	if insternetIDOk {
		params.Filters = &oscgo.FiltersInternetService{
			InternetServiceIds: &[]string{internetID.(string)},
		}

	}

	if filtersOk {
		params.Filters, err = buildOutscaleOSCAPIDataSourceInternetServiceFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp oscgo.ReadInternetServicesResponse

	err = retry.Retry(120*time.Second, func() *retry.RetryError {
		var err error
		rp, httpResp, err := conn.InternetServiceApi.ReadInternetServices(context.Background()).ReadInternetServicesRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading Internet Service id (%s)", utils.GetErrorResponse(err))
	}

	if !resp.HasInternetServices() || len(resp.GetInternetServices()) == 0 {
		return fmt.Errorf("Error reading Internet Service: Internet Services is not found with the seatch criteria")
	}

	result := resp.GetInternetServices()[0]

	log.Printf("[DEBUG] Setting OAPI Internet Service id (%s)", err)

	if err := d.Set("internet_service_id", result.GetInternetServiceId()); err != nil {
		return err
	}

	if err := d.Set("state", result.GetState()); err != nil {
		return err
	}

	if err := d.Set("net_id", result.GetNetId()); err != nil {
		return err
	}

	d.SetId(result.GetInternetServiceId())

	return d.Set("tags", FlattenOAPITagsSDK(result.GetTags()))
}

func buildOutscaleOSCAPIDataSourceInternetServiceFilters(set *schema.Set) (*oscgo.FiltersInternetService, error) {
	var filters oscgo.FiltersInternetService
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "internet_service_ids":
			filters.SetInternetServiceIds(filterValues)
		case "link_net_ids":
			filters.SetLinkNetIds(filterValues)
		case "link_states":
			filters.SetLinkStates(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		default:
			return nil, utils.UnknownDataSourceFilterError(context.Background(), name)
		}
	}
	return &filters, nil
}
