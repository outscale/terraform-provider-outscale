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

func datasourceOutscaleOAPIInternetServices() *schema.Resource {
	return &schema.Resource{
		Read: datasourceOutscaleOAPIInternetServicesRead,
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
						"tags": dataSourceTagsSchema(),
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

func datasourceOutscaleOAPIInternetServicesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	internetID, internetIDOk := d.GetOk("internet_service_ids")

	if !filtersOk && !internetIDOk {
		return fmt.Errorf("One of filters, or instance_id must be assigned")
	}

	// Build up search parameters
	params := oscgo.ReadInternetServicesRequest{
		Filters: &oscgo.FiltersInternetService{},
	}
	filter := oscgo.FiltersInternetService{}
	if internetIDOk {
		i := internetID.([]string)
		in := make([]string, len(i))
		copy(in, i)
		filter.SetInternetServiceIds(in)
		params.SetFilters(filter)
	}

	if filtersOk {
		params.Filters = buildOutscaleOSCAPIDataSourceInternetServiceFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadInternetServicesResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.InternetServiceApi.ReadInternetServices(context.Background()).ReadInternetServicesRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string

	if err != nil {
		errString = err.Error()

		return fmt.Errorf("[DEBUG] Error reading Internet Services (%s)", errString)
	}

	log.Printf("[DEBUG] Setting OAPI LIN Internet Gateways id (%s)", err)

	d.SetId(resource.UniqueId())

	result := resp.GetInternetServices()
	return internetServicesOAPIDescriptionAttributes(d, result)
}

func internetServicesOAPIDescriptionAttributes(d *schema.ResourceData, internetServices []oscgo.InternetService) error {

	i := make([]map[string]interface{}, len(internetServices))
	for k, v := range internetServices {
		im := make(map[string]interface{})
		if v.GetState() != "" {
			im["state"] = v.GetState()
		}

		if v.GetNetId() != "" {
			im["net_id"] = v.GetNetId()
		}
		if v.GetInternetServiceId() != "" {
			im["internet_service_id"] = v.GetInternetServiceId()
		}
		if v.Tags != nil {
			im["tags"] = tagsOSCAPIToMap(v.GetTags())
		}
		i[k] = im
	}

	return d.Set("internet_services", i)
}
