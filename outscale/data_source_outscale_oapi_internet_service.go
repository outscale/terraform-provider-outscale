package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"tags": tagsOAPIListSchemaComputed(),
		},
	}
}

func datasourceOutscaleOAPIInternetServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	filters, filtersOk := d.GetOk("filter")
	internetID, insternetIDOk := d.GetOk("internet_service_id")

	if filtersOk == false && insternetIDOk == false {
		return fmt.Errorf("One of filters, or instance_id must be assigned")
	}

	// Build up search parameters
	params := &oapi.ReadInternetServicesRequest{}

	if insternetIDOk {
		params.Filters = oapi.FiltersInternetService{
			InternetServiceIds: []string{internetID.(string)},
		}

	}

	if filtersOk {
		params.Filters = buildOutscaleOAPIDataSourceInternetServiceFilters(filters.(*schema.Set))
	}

	var resp *oapi.POST_ReadInternetServicesResponses
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		resp, err = conn.POST_ReadInternetServices(*params)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})

	var errString string

	if err != nil || resp.OK == nil {
		if err != nil {
			errString = err.Error()

		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}

		return fmt.Errorf("[DEBUG] Error reading Internet Service id (%s)", errString)
	}

	if len(resp.OK.InternetServices) <= 0 {
		return fmt.Errorf("Error reading Internet Service: Internet Services is not found with the seatch criteria")
	}

	result := resp.OK.InternetServices[0]

	log.Printf("[DEBUG] Setting OAPI Internet Service id (%s)", err)

	d.Set("request_id", resp.OK.ResponseContext.RequestId)
	d.Set("internet_service_id", result.InternetServiceId)
	d.Set("state", result.State)
	d.Set("net_id", result.NetId)

	return d.Set("tags", tagsOAPIToMap(result.Tags))
}

func buildOutscaleOAPIDataSourceInternetServiceFilters(set *schema.Set) oapi.FiltersInternetService {
	var filters oapi.FiltersInternetService
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "internet_service_ids":
			filters.InternetServiceIds = filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
