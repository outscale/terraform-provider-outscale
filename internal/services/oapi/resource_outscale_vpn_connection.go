package oapi

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"
)

func ResourceOutscaleVPNConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleVPNclientectionCreate,
		ReadContext:   ResourceOutscaleVPNclientectionRead,
		UpdateContext: ResourceOutscaleVPNclientectionUpdate,
		DeleteContext: ResourceOutscaleVPNclientectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Update: schema.DefaultTimeout(UpdateDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
		},

		Schema: map[string]*schema.Schema{
			"client_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"virtual_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_connection_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"static_routes_only": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
				ForceNew: true,
			},
			"client_gateway_configuration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpn_clientection_id": {
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
			"tags": TagsSchemaSDK(),
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

func ResourceOutscaleVPNclientectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutCreate)

	req := osc.CreateVpnConnectionRequest{
		ClientGatewayId:  d.Get("client_gateway_id").(string),
		VirtualGatewayId: d.Get("virtual_gateway_id").(string),
		ConnectionType:   d.Get("client_connection_type").(string),
	}

	if staticRoutesOnly, ok := d.GetOkExists("static_routes_only"); ok {
		req.StaticRoutesOnly = ptr.To(cast.ToBool(staticRoutesOnly))
	}
	resp, err := client.CreateVpnConnection(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating outscale vpn conecction: %s", err))
	}

	d.SetId(*resp.VpnConnection.VpnConnectionId)

	err = createOAPITagsSDK(ctx, client, timeout, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceOutscaleVPNclientectionRead(ctx, d, meta)
}

func ResourceOutscaleVPNclientectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	vpnclientectionID := d.Id()

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available", "failed"},
		Refresh: vpnclientectionRefreshFunc(ctx, client, timeout, &vpnclientectionID),
	}

	r, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for outscale vpn clientection(%s) to become ready: %s", vpnclientectionID, err))
	}

	resp := r.(osc.ReadVpnConnectionsResponse)
	if resp.VpnConnections == nil || utils.IsResponseEmpty(len(*resp.VpnConnections), "Vpnclientection", d.Id()) {
		d.SetId("")
		return nil
	}
	vpnConnection := (*resp.VpnConnections)[0]

	if err := d.Set("client_gateway_configuration", vpnConnection.ClientGatewayConfiguration); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vpn_clientection_id", vpnConnection.VpnConnectionId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", vpnConnection.State); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("static_routes_only", vpnConnection.StaticRoutesOnly); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("client_gateway_id", vpnConnection.ClientGatewayId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("virtual_gateway_id", vpnConnection.VirtualGatewayId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("clientection_type", vpnConnection.ConnectionType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("routes", flattenVPNclientection(ptr.From(vpnConnection.Routes))); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(ptr.From(vpnConnection.Tags))); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vgw_telemetries", flattenVgwTelemetries(ptr.From(vpnConnection.VgwTelemetries))); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceOutscaleVPNclientectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutUpdate)

	if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
		return diag.FromErr(err)
	}

	return ResourceOutscaleVPNclientectionRead(ctx, d, meta)
}

func ResourceOutscaleVPNclientectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutDelete)

	id := d.Id()

	req := osc.DeleteVpnConnectionRequest{
		VpnConnectionId: id,
	}
	_, err := client.DeleteVpnConnection(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted", "failed"},
		Refresh: vpnclientectionRefreshFunc(ctx, client, timeout, &id),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for outscale vpn clientection(%s) to become deleted: %s", id, err))
	}

	return nil
}

func vpnclientectionRefreshFunc(ctx context.Context, client *osc.Client, timeout time.Duration, id *string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		filter := osc.ReadVpnConnectionsRequest{
			Filters: &osc.FiltersVpnConnection{
				VpnConnectionIds: &[]string{*id},
			},
		}
		resp, err := client.ReadVpnConnections(ctx, filter, options.WithRetryTimeout(timeout))
		if err != nil {
			return nil, "failed", fmt.Errorf("error on vpnconnectionrefresh: %s", err)
		}
		if resp.VpnConnections == nil || len(*resp.VpnConnections) == 0 {
			return nil, "failed", fmt.Errorf("error on vpnconnectionrefresh: there are not vpn connections(%s)", *id)
		}

		vpnConnection := (*resp.VpnConnections)[0]

		return resp, *vpnConnection.State, nil
	}
}

func flattenVPNclientection(routes []osc.RouteLight) []map[string]interface{} {
	routesMap := make([]map[string]interface{}, len(routes))

	for i, route := range routes {
		routesMap[i] = map[string]interface{}{
			"destination_ip_range": route.DestinationIpRange,
			"route_type":           route.RouteType,
			"state":                route.State,
		}
	}

	return routesMap
}

func flattenVgwTelemetries(vgwTelemetries []osc.VgwTelemetry) []map[string]interface{} {
	vgwTelemetriesMap := make([]map[string]interface{}, len(vgwTelemetries))

	for i, vgwTelemetry := range vgwTelemetries {
		vgwTelemetriesMap[i] = map[string]interface{}{
			"accepted_route_count":   vgwTelemetry.AcceptedRouteCount,
			"last_state_change_date": vgwTelemetry.LastStateChangeDate,
			"outside_ip_address":     vgwTelemetry.OutsideIpAddress,
			"state":                  vgwTelemetry.State,
			"state_description":      vgwTelemetry.StateDescription,
		}
	}

	return vgwTelemetriesMap
}
