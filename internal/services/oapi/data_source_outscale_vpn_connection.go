package oapi

import (
	"context"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOutscaleVPNConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleVPNConnectionRead,
		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpn_connection_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"virtual_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"static_routes_only": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"client_gateway_configuration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_ip_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"route_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"tags": TagsSchemaComputedSDK(),
			"vgw_telemetries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"accepted_route_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_state_change_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"outside_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_description": {
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

func DataSourceOutscaleVPNConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	vpnConnectionID, vpnConnectionOk := d.GetOk("vpn_connection_id")

	if !filtersOk && !vpnConnectionOk {
		return diag.Errorf("one of filters, or vpn_connection_id must be assigned")
	}

	params := osc.ReadVpnConnectionsRequest{}

	if vpnConnectionOk {
		params.Filters = &osc.FiltersVpnConnection{
			VpnConnectionIds: &[]string{vpnConnectionID.(string)},
		}
	}

	var err error
	if filtersOk {
		params.Filters, err = buildOutscaleDataSourceVPNConnectionFilters(filters.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resp, err := client.ReadVpnConnections(ctx, params, options.WithRetryTimeout(5*time.Minute))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.VpnConnections == nil || len(*resp.VpnConnections) == 0 {
		return diag.FromErr(ErrNoResults)
	}
	if len(*resp.VpnConnections) > 1 {
		return diag.FromErr(ErrMultipleResults)
	}

	vpnConnection := (*resp.VpnConnections)[0]

	if err := d.Set("client_gateway_id", ptr.From(vpnConnection.ClientGatewayId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("virtual_gateway_id", ptr.From(vpnConnection.VirtualGatewayId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("connection_type", ptr.From(vpnConnection.ConnectionType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("static_routes_only", ptr.From(vpnConnection.StaticRoutesOnly)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("client_gateway_configuration", ptr.From(vpnConnection.ClientGatewayConfiguration)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpn_connection_id", ptr.From(vpnConnection.VpnConnectionId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", ptr.From(vpnConnection.State)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("routes", flattenVPNConnection(ptr.From(vpnConnection.Routes))); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(ptr.From(vpnConnection.Tags))); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vgw_telemetries", flattenVgwTelemetries(ptr.From(vpnConnection.VgwTelemetries))); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ptr.From(vpnConnection.VpnConnectionId))

	return nil
}

func buildOutscaleDataSourceVPNConnectionFilters(set *schema.Set) (*osc.FiltersVpnConnection, error) {
	var filters osc.FiltersVpnConnection
	for _, v := range set.List() {
		m := v.(map[string]interface{})
		var filterValues []string
		for _, e := range m["values"].([]interface{}) {
			filterValues = append(filterValues, e.(string))
		}

		var filteBgpAsnsValues []int
		for _, e := range m["values"].([]interface{}) {
			filteBgpAsnsValues = append(filteBgpAsnsValues, cast.ToInt(e))
		}

		switch name := m["name"].(string); name {
		case "vpn_connection_ids":
			filters.VpnConnectionIds = &filterValues
		case "virtual_gateway_ids":
			filters.VirtualGatewayIds = &filterValues
		case "client_gateway_ids":
			filters.ClientGatewayIds = &filterValues
		case "connection_types":
			filters.ConnectionTypes = &filterValues
		case "route_destination_ip_ranges":
			filters.RouteDestinationIpRanges = &filterValues
		case "states":
			filters.States = &filterValues
		case "static_routes_only":
			filters.StaticRoutesOnly = new(cast.ToBool(filterValues[0]))
		case "bgp_asns":
			filters.BgpAsns = &filteBgpAsnsValues
		case "tag_keys":
			filters.TagKeys = &filterValues
		case "tag_values":
			filters.TagValues = &filterValues
		case "tags":
			filters.Tags = &filterValues
		default:
			return nil, utils.UnknownDataSourceFilterError(name)
		}
	}
	return &filters, nil
}
