package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleVirtualGateway() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleVirtualGatewayRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"virtual_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connection_type": {
				Type:     schema.TypeString,
				Optional: true,
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": TagsSchemaComputedSDK(),
		},
	}
}

func DataSourceOutscaleVirtualGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	virtualId, vpnOk := d.GetOk("virtual_gateway_id")

	if !filtersOk && !vpnOk {
		return diag.Errorf("one of virtual_gateway_id or filter must be assigned")
	}

	params := osc.ReadVirtualGatewaysRequest{}

	if vpnOk {
		params.Filters = &osc.FiltersVirtualGateway{VirtualGatewayIds: &[]string{virtualId.(string)}}
	}

	var err error
	if filtersOk {
		params.Filters, err = buildOutscaleAPIVirtualGatewayFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadVirtualGateways(ctx, params, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.VirtualGateways == nil || len(*resp.VirtualGateways) == 0 {
		return diag.Errorf("no matching virtual gateway found: %#v", params)
	}
	if len(*resp.VirtualGateways) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	vgw := (*resp.VirtualGateways)[0]

	d.SetId(ptr.From(vgw.VirtualGatewayId))
	vs := make([]map[string]interface{}, len(ptr.From(vgw.NetToVirtualGatewayLinks)))

	for k, v := range ptr.From(vgw.NetToVirtualGatewayLinks) {
		vp := make(map[string]interface{})

		vp["state"] = v.State
		vp["net_id"] = v.NetId

		vs[k] = vp
	}

	d.Set("net_to_virtual_gateway_links", vs)
	d.Set("state", ptr.From(vgw.State))
	d.Set("connection_type", ptr.From(vgw.ConnectionType))
	d.Set("tags", FlattenOAPITagsSDK(ptr.From(vgw.Tags)))

	return nil
}

func buildOutscaleAPIVirtualGatewayFilters(set *schema.Set) (*osc.FiltersVirtualGateway, error) {
	var filters osc.FiltersVirtualGateway
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}
		switch name := m["name"].(string); name {
		// case "available_ips_counts":
		// 	filters.AvailableIpsCounts = filterValues
		case "tags":
			filters.Tags = &filterValues
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "states":
			filters.States = &filterValues
		case "connection_types":
			filters.ConnectionTypes = &filterValues
		case "link_net_ids":
			filters.LinkNetIds = &filterValues
		case "link_states":
			filters.LinkStates = &filterValues
		case "virtual_gateway_ids":
			filters.VirtualGatewayIds = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
