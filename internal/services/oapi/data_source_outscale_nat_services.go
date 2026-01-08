package oapi

import (
	"context"
	"fmt"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleNatServices() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleNatServicesRead,

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
						"tags": TagsSchemaComputedSDK(),
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

func DataSourceOutscaleNatServicesRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	filters, filtersOk := d.GetOk("filter")
	natGatewayID, natGatewayIDOK := d.GetOk("nat_service_ids")

	if !filtersOk && !natGatewayIDOK {
		return fmt.Errorf("filters, or owner must be assigned, or nat_service_id must be provided")
	}

	var err error
	params := oscgo.ReadNatServicesRequest{}
	if filtersOk {
		params.Filters, err = buildOutscaleNatServiceDataSourceFilters(filters.(*schema.Set))
		if err != nil {
			return err
		}
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
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
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

		return fmt.Errorf("error reading nat service (%s)", errString)
	}

	if len(resp.GetNatServices()) < 1 {
		return ErrNoResults
	}

	return ngsOAPIDescriptionAttributes(d, resp.GetNatServices())
}

// populate the numerous fields that the image description returns.
func ngsOAPIDescriptionAttributes(d *schema.ResourceData, ngs []oscgo.NatService) error {
	d.SetId(id.UniqueId())

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
			addng["tags"] = FlattenOAPITagsSDK(v.GetTags())
		}

		addngs[k] = addng
	}

	return d.Set("nat_services", addngs)
}
