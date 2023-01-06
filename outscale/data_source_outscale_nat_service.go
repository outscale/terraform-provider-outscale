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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": dataSourceTagsSchema(),
		},
	}
}

func dataSourceOutscaleOAPINatServiceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	natGatewayID, natGatewayIDOK := d.GetOk("nat_service_id")

	if !filtersOk && !natGatewayIDOK {
		return fmt.Errorf("filters, or owner must be assigned, or nat_service_id must be provided")
	}

	params := oscgo.ReadNatServicesRequest{}

	if filtersOk {
		params.SetFilters(buildOutscaleOAPINatServiceDataSourceFilters(filters.(*schema.Set)))
	}
	if natGatewayIDOK && natGatewayID.(string) != "" {
		filter := oscgo.FiltersNatService{}
		filter.SetNatServiceIds([]string{natGatewayID.(string)})
		params.SetFilters(filter)
	}

	var resp oscgo.ReadNatServicesResponse
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		rp, httpResp, err := conn.NatServiceApi.ReadNatServices(context.Background()).ReadNatServicesRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	var errString string

	if err != nil {
		errString = err.Error()

		return fmt.Errorf("[DEBUG] Error reading Nar Service (%s)", errString)
	}

	if len(resp.GetNatServices()) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	if len(resp.GetNatServices()) > 1 {
		return fmt.Errorf("your query returned more than one result, please try a more " +
			"specific search criteria")
	}

	return ngOAPIDescriptionAttributes(d, resp.GetNatServices()[0])
}

// populate the numerous fields that the image description returns.
func ngOAPIDescriptionAttributes(d *schema.ResourceData, ng oscgo.NatService) error {

	d.SetId(ng.GetNatServiceId())

	if err := d.Set("nat_service_id", ng.GetNatServiceId()); err != nil {
		return err
	}

	if ng.GetState() != "" {
		if err := d.Set("state", ng.State); err != nil {
			return err
		}

	}
	if ng.GetSubnetId() != "" {
		if err := d.Set("subnet_id", ng.GetSubnetId()); err != nil {
			return err
		}

	}
	if ng.GetNetId() != "" {
		if err := d.Set("net_id", ng.GetNetId()); err != nil {
			return err
		}

	}

	addresses := make([]map[string]interface{}, len(ng.GetPublicIps()))

	for k, v := range ng.GetPublicIps() {
		address := make(map[string]interface{})
		if v.GetPublicIpId() != "" {
			address["public_ip_id"] = v.PublicIpId
		}
		if v.GetPublicIp() != "" {
			address["public_ip"] = v.PublicIp
		}
		addresses[k] = address
	}
	if err := d.Set("public_ips", addresses); err != nil {
		return err
	}
	if err := d.Set("tags", tagsOSCAPIToMap(ng.GetTags())); err != nil {
		return err
	}

	return nil
}

func buildOutscaleOAPINatServiceDataSourceFilters(set *schema.Set) oscgo.FiltersNatService {
	var filters oscgo.FiltersNatService
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		switch name := m["name"].(string); name {
		case "nat_service_ids":
			filters.SetNatServiceIds(filterValues)
		case "net_ids":
			filters.SetNetIds(filterValues)
		case "states":
			filters.SetStates(filterValues)
		case "subnet_ids":
			filters.SetSubnetIds(filterValues)
		case "tag_keys":
			filters.SetTagKeys(filterValues)
		case "tag_values":
			filters.SetTagValues(filterValues)
		case "tags":
			filters.SetTags(filterValues)
		default:
			log.Printf("[Debug] Unknown Filter Name: %s.", name)
		}
	}
	return filters
}
