package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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
						"internet_service_ids": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": tagsOAPIListSchemaComputed(),
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
	conn := meta.(*OutscaleClient).OAPI

	filters, filtersOk := d.GetOk("filter")
	internetID, internetIDOk := d.GetOk("internet_service_ids")

	if filtersOk == false && internetIDOk == false {
		return fmt.Errorf("One of filters, or instance_id must be assigned")
	}

	// Build up search parameters
	params := &oapi.ReadInternetServicesRequest{
		Filters: oapi.FiltersInternetService{},
	}

	if internetIDOk {
		i := internetID.([]string)
		in := make([]string, len(i))
		for k, v := range i {
			in[k] = v
		}
		params.Filters.InternetServiceIds = in
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

		return fmt.Errorf("[DEBUG] Error reading Internet Services (%s)", errString)
	}

	log.Printf("[DEBUG] Setting OAPI LIN Internet Gateways id (%s)", err)

	d.Set("request_id", resp.OK.ResponseContext.RequestId)
	d.SetId(resource.UniqueId())

	result := resp.OK.InternetServices
	return internetServicesOAPIDescriptionAttributes(d, result)
}

func flattenOAPIInternetGwsAttachements(attachements []*fcu.InternetGatewayAttachment) []map[string]interface{} {
	res := make([]map[string]interface{}, len(attachements))

	for i, a := range attachements {
		res[i]["state"] = a.State
		res[i]["net_id"] = a.VpcId
	}

	return res
}

func internetServicesOAPIDescriptionAttributes(d *schema.ResourceData, internetServices []oapi.InternetService) error {

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
			im["tags"] = tagsOAPIToMap(v.Tags)
		}
		i[k] = im
	}

	return d.Set("internet_services", i)
}
