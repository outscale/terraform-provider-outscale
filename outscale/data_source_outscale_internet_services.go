package outscale

import (
	"context"
	"fmt"
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
			"filter": dataSourceFiltersSchema(false),
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
	req := oscgo.ReadInternetServicesRequest{}
	if filters, filtersOk := d.GetOk("filter"); filtersOk {
		req.Filters = buildOutscaleOSCAPIDataSourceInternetServiceFilters(filters.(*schema.Set))
	}

	var resp oscgo.ReadInternetServicesResponse
	var err error
	err = resource.Retry(120*time.Second, func() *resource.RetryError {
		rp, httpResp, err := conn.InternetServiceApi.ReadInternetServices(context.Background()).ReadInternetServicesRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("[DEBUG] Error reading Internet Services (%s)", err.Error())
	}

	result := resp.GetInternetServices()
	if len(result) == 0 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	d.SetId(resource.UniqueId())

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
