package oapi

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleVirtualGateways() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceOutscaleVirtualGatewaysRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"virtual_gateway_id": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"virtual_gateways": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"clientection_type": {
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

func DataSourceOutscaleVirtualGatewaysRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*client.OutscaleClient).OSC

	filter, filtersOk := d.GetOk("filter")
	_, vpnOk := d.GetOk("virtual_gateway_id")

	if !filtersOk && !vpnOk {
		return fmt.Errorf("one of virtual_gateway_id or filter must be assigned")
	}

	var err error
	params := osc.ReadVirtualGatewaysRequest{}
	if filtersOk {
		params.Filters, err = buildOutscaleAPIVirtualGatewayFilters(filter.(*schema.Set))
		if err != nil {
			return err
		}
	}

	var resp osc.ReadVirtualGatewaysResponse
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := client.VirtualGatewayApi.ReadVirtualGateways(ctx).ReadVirtualGatewaysRequest(params).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return err
	}
	if resp.GetVirtualGateways() == nil || len(resp.GetVirtualGateways()) == 0 {
		return ErrNoResults
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
		vpn["clientection_type"] = v.GetclientectionType()
		vpn["virtual_gateway_id"] = v.GetVirtualGatewayId()
		vpn["tags"] = FlattenOAPITagsSDK(v.Tags)

		vpns[k] = vpn
	}
	d.Set("virtual_gateways", vpns)
	d.SetId(id.UniqueId())

	return nil
}
