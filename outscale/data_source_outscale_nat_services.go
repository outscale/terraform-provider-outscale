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

func dataSourceOutscaleOAPINatServices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPINatServicesRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"nat_service_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Attributes
			"nat_services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"nat_service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_id": {
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

func dataSourceOutscaleOAPINatServicesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	natGatewayID, natGatewayIDOK := d.GetOk("nat_service_ids")

	if !filtersOk && !natGatewayIDOK {
		return fmt.Errorf("filters, or owner must be assigned, or nat_service_id must be provided")
	}

	params := oscgo.ReadNatServicesRequest{}
	if filtersOk {
		params.SetFilters(buildOutscaleOAPINatServiceDataSourceFilters(filters.(*schema.Set)))
	}
	if natGatewayIDOK {
		ids := make([]string, len(natGatewayID.([]interface{})))

		for k, v := range natGatewayID.([]interface{}) {
			ids[k] = v.(string)
		}
		filter := oscgo.FiltersNatService{}
		filter.SetNatServiceIds(ids)
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

	return ngsOAPIDescriptionAttributes(d, resp.GetNatServices())
}

// populate the numerous fields that the image description returns.
func ngsOAPIDescriptionAttributes(d *schema.ResourceData, ngs []oscgo.NatService) error {
	d.SetId(resource.UniqueId())

	addngs := make([]map[string]interface{}, len(ngs))

	for k, v := range ngs {
		addng := make(map[string]interface{})

		ngas := make([]interface{}, len(v.GetPublicIps()))

		for i, w := range v.GetPublicIps() {
			nga := make(map[string]interface{})
			if w.GetPublicIpId() != "" {
				nga["public_ip_id"] = w.GetPublicIpId()
			}
			if w.GetPublicIp() != "" {
				nga["public_ip"] = w.GetPublicIp()
			}
			ngas[i] = nga
		}
		addng["public_ips"] = ngas

		if v.GetNatServiceId() != "" {
			addng["nat_service_id"] = v.GetNatServiceId()
		}
		if v.GetState() != "" {
			addng["state"] = v.GetState()
		}
		if v.GetSubnetId() != "" {
			addng["subnet_id"] = v.GetSubnetId()
		}
		if v.GetNetId() != "" {
			addng["net_id"] = v.GetNetId()
		}
		if v.GetTags() != nil {
			addng["tags"] = tagsOSCAPIToMap(v.GetTags())
		}

		addngs[k] = addng
	}

	return d.Set("nat_services", addngs)
}
