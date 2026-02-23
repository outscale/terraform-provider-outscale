package oapi

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/id"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func DataSourceOutscaleVPNConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceOutscaleVPNConnectionsRead,

		Schema: map[string]*schema.Schema{
			"filter": dataSourceFiltersSchema(),
			"vpn_connection_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpn_connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpn_connection_id": {
							Type:     schema.TypeString,
							Computed: true,
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

func DataSourceOutscaleVPNConnectionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	filters, filtersOk := d.GetOk("filter")
	vpnConnectionIDs, vpnConnectionOk := d.GetOk("vpn_connection_ids")

	if !filtersOk && !vpnConnectionOk {
		return diag.Errorf("one of filters, or vpn_connection_ids must be assigned")
	}

	log.Printf("vpnconnectionIDs: %#+v\n", vpnConnectionIDs)
	params := osc.ReadVpnConnectionsRequest{}

	if vpnConnectionOk {
		params.Filters = &osc.FiltersVpnConnection{
			VpnConnectionIds: utils.InterfaceSliceToStringSlicePtr(vpnConnectionIDs.([]interface{})),
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
	if err := d.Set("vpn_connections", flattenVPNConnections(*resp.VpnConnections)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id.UniqueId())
	return nil
}

func flattenVPNConnections(vpnConnections []osc.VpnConnection) []map[string]interface{} {
	vpnConnectionsMap := make([]map[string]interface{}, len(vpnConnections))

	for i, vpnConnection := range vpnConnections {
		vpnConnectionsMap[i] = map[string]interface{}{
			"vpn_connection_id":            vpnConnection.VpnConnectionId,
			"client_gateway_id":            vpnConnection.ClientGatewayId,
			"virtual_gateway_id":           vpnConnection.VirtualGatewayId,
			"connection_type":              vpnConnection.ConnectionType,
			"static_routes_only":           vpnConnection.StaticRoutesOnly,
			"client_gateway_configuration": vpnConnection.ClientGatewayConfiguration,
			"state":                        vpnConnection.State,
			"routes":                       flattenVPNConnection(ptr.From(vpnConnection.Routes)),
			"tags":                         FlattenOAPITagsSDK(ptr.From(vpnConnection.Tags)),
			"vgw_telemetries":              flattenVgwTelemetries(ptr.From(vpnConnection.VgwTelemetries)),
		}
	}
	return vpnConnectionsMap
}
