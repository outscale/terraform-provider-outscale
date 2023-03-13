package outscale

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func dataSourceOutscaleOAPIVirtualGateways() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceOutscaleOAPIVirtualGatewaysRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(false),
			"virtual_gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_gateway_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_to_virtual_gateway_links": {
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
								},
							},
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

func dataSourceOutscaleOAPIVirtualGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	req := oscgo.ReadVirtualGatewaysRequest{}

	if filter, filtersOk := d.GetOk("filter"); filtersOk {
		req.SetFilters(buildOutscaleAPIVirtualGatewayFilters(filter.(*schema.Set)))
	}

	var resp oscgo.ReadVirtualGatewaysResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.VirtualGatewayApi.ReadVirtualGateways(context.Background()).ReadVirtualGatewaysRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}
	if !resp.HasVirtualGateways() || len(resp.GetVirtualGateways()) < 1 {
		return fmt.Errorf("your query returned no results, please change your search criteria and try again")
	}

	vpns := make([]map[string]interface{}, len(resp.GetVirtualGateways()))

	for k, v := range resp.GetVirtualGateways() {
		vpn := make(map[string]interface{})
		vs := make([]map[string]interface{}, len(v.GetNetToVirtualGatewayLinks()))

		for k, v1 := range v.GetNetToVirtualGatewayLinks() {
			vp := make(map[string]interface{})
			vp["state"] = v1.GetState()
			vp["net_id"] = v1.GetNetId()

			vs[k] = vp
		}
		vpn["net_to_virtual_gateway_links"] = vs
		vpn["state"] = v.GetState()
		vpn["connection_type"] = v.GetConnectionType()
		vpn["virtual_gateway_id"] = v.GetVirtualGatewayId()
		vpn["tags"] = tagsOSCAPIToMap(v.GetTags())

		vpns[k] = vpn
	}
	d.Set("virtual_gateways", vpns)
	d.SetId(resource.UniqueId())

	return nil
}
