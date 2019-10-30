package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/outscale/osc-go/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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
	conn := meta.(*OutscaleClient).OAPI

	filters, filtersOk := d.GetOk("filter")
	natGatewayID, natGatewayIDOK := d.GetOk("nat_service_ids")

	if filtersOk == false && natGatewayIDOK == false {
		return fmt.Errorf("filters, or owner must be assigned, or nat_service_id must be provided")
	}

	params := &oapi.ReadNatServicesRequest{}
	if filtersOk {
		params.Filters = buildOutscaleOAPINatServiceDataSourceFilters(filters.(*schema.Set))
	}
	if natGatewayIDOK {
		ids := make([]string, len(natGatewayID.([]interface{})))

		for k, v := range natGatewayID.([]interface{}) {
			ids[k] = v.(string)
		}

		params.Filters.NatServiceIds = ids
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

	d.Set("request_id", response.ResponseContext.RequestId)

	return ngsOAPIDescriptionAttributes(d, response.NatServices)
}

// populate the numerous fields that the image description returns.
func ngsOAPIDescriptionAttributes(d *schema.ResourceData, ngs []oapi.NatService) error {

	d.SetId(resource.UniqueId())

	addngs := make([]map[string]interface{}, len(ngs))

	for k, v := range ngs {
		addng := make(map[string]interface{})

		ngas := make([]interface{}, len(v.PublicIps))

		for i, w := range v.PublicIps {
			nga := make(map[string]interface{})
			if w.PublicIpId != "" {
				nga["public_ip_id"] = w.PublicIpId
			}
			if w.PublicIp != "" {
				nga["public_ip"] = w.PublicIp
			}
			ngas[i] = nga
		}
		addng["public_ips"] = ngas

		if v.NatServiceId != "" {
			addng["nat_service_id"] = v.NatServiceId
		}
		if v.State != "" {
			addng["state"] = v.State
		}
		if v.SubnetId != "" {
			addng["subnet_id"] = v.SubnetId
		}
		if v.NetId != "" {
			addng["net_id"] = v.NetId
		}

		addngs[k] = addng
	}

	return d.Set("nat_services", addngs)
}
