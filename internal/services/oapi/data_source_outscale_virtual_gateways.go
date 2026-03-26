package oapi

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
)

func DataSourceOutscaleVirtualGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleVirtualGatewaysRead,

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

func DataSourceOutscaleVirtualGatewaysRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filter, filtersOk := d.GetOk("filter")
	_, vpnOk := d.GetOk("virtual_gateway_id")

	if !filtersOk && !vpnOk {
		return diag.Errorf("one of virtual_gateway_id or filter must be assigned")
	}

	var err error
	params := osc.ReadVirtualGatewaysRequest{}
	if filtersOk {
		params.Filters, err = buildOutscaleAPIVirtualGatewayFilters(filter.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadVirtualGateways(ctx, params, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.VirtualGateways == nil || len(*resp.VirtualGateways) == 0 {
		return diag.FromErr(ErrNoResults)
	}

	vpns := make([]map[string]any, len(*resp.VirtualGateways))

	for k, v := range *resp.VirtualGateways {
		vpn := make(map[string]any)
		vs := make([]map[string]any, len(v.NetToVirtualGatewayLinks))

		for k, v1 := range v.NetToVirtualGatewayLinks {
			vp := make(map[string]any)
			vp["state"] = v1.State
			vp["net_id"] = ptr.From(v1.NetId)

			vs[k] = vp
		}
		vpn["net_to_virtual_gateway_links"] = vs
		vpn["state"] = v.State
		vpn["connection_type"] = v.ConnectionType
		vpn["virtual_gateway_id"] = v.VirtualGatewayId
		vpn["tags"] = FlattenOAPITagsSDK(v.Tags)

		vpns[k] = vpn
	}
	if err := d.Set("virtual_gateways", vpns); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id.UniqueId())

	return nil
}
