package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPINatService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPINatServiceRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"nat_service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Attributes
			"public_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceOutscaleOAPINatServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	filters, filtersOk := d.GetOk("filter")
	natGatewayID, natGatewayIDOK := d.GetOk("nat_service_id")

	if filtersOk == false && natGatewayIDOK == false {
		return fmt.Errorf("filters, or owner must be assigned, or nat_service_id must be provided")
	}

	params := &oapi.ReadNatServicesRequest{}

	if filtersOk {
		params.Filters = buildOutscaleOAPINatServiceDataSourceFilters(filters.(*schema.Set))
	}
	if natGatewayIDOK && natGatewayID.(string) != "" {
		params.Filters.NatServiceIds = []string{natGatewayID.(string)}
	}

	var resp *oapi.POST_ReadNatServicesResponses
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error

		resp, err = conn.POST_ReadNatServices(*params)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
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

		return fmt.Errorf("[DEBUG] Error reading Nar Service (%s)", errString)
	}

	response := resp.OK

	if len(response.NatServices) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(response.NatServices) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more " +
			"specific search criteria")
	}

	return ngOAPIDescriptionAttributes(d, response.NatServices[0])
}

// populate the numerous fields that the image description returns.
func ngOAPIDescriptionAttributes(d *schema.ResourceData, ng oapi.NatService) error {

	d.SetId(ng.NatServiceId)
	d.Set("nat_service_id", ng.NatServiceId)

	if ng.State != "" {
		d.Set("state", ng.State)
	}
	if ng.SubnetId != "" {
		d.Set("subnet_id", ng.SubnetId)
	}
	if ng.NetId != "" {
		d.Set("net_id", ng.NetId)
	}

	addresses := make([]map[string]interface{}, len(ng.PublicIps))

	for k, v := range ng.PublicIps {
		address := make(map[string]interface{})
		if v.PublicIpId != "" {
			address["public_ip_id"] = v.PublicIpId
		}
		if v.PublicIp != "" {
			address["public_ip"] = v.PublicIp
		}
		addresses[k] = address
	}
	if err := d.Set("public_ips", addresses); err != nil {
		return err
	}

	return nil
}

func buildOutscaleOAPINatServiceDataSourceFilters(set *schema.Set) oapi.FiltersNatService {
	var filters oapi.FiltersNatService
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "nat_service_ids":
			filters.NatServiceIds = filterValues
		case "net_ids":
			filters.NetIds = filterValues
		case "states":
			filters.States = filterValues
		case "subnet_ids":
			filters.SubnetIds = filterValues
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
