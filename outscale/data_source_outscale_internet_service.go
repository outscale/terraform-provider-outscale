package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceOutscaleOAPIInternetService() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOutscaleOAPIInternetServiceRead,
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
			"tags": dataSourceTagsSchema(),
		},
	}
}

func datasourceOutscaleOAPIInternetServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	params := oscgo.ReadInternetServicesRequest{}
	if filters, filtersOk := d.GetOk("filter"); filtersOk {
		params.Filters = buildOutscaleOSCAPIDataSourceInternetServiceFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadInternetServicesResponse

	err := resource.Retry(120*time.Second, func() *resource.RetryError {
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
	if err := utils.IsResponseEmptyOrMutiple(len(resp.GetInternetServices()), "Image Export Task"); err != nil {
		return err
	}

	result := resp.GetInternetServices()[0]

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

	return d.Set("tags", tagsOSCAPIToMap(result.GetTags()))
}

func buildOutscaleOSCAPIDataSourceInternetServiceFilters(set *schema.Set) *oscgo.FiltersInternetService {
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
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return &filters
}
