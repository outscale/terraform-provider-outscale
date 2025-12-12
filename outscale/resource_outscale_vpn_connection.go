package outscale

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/terraform-provider-outscale/utils"
	"github.com/spf13/cast"

	oscgo "github.com/outscale/osc-sdk-go/v2"
)

func ResourceOutscaleVPNConnection() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleVPNConnectionCreate,
		Read:   ResourceOutscaleVPNConnectionRead,
		Update: ResourceOutscaleVPNConnectionUpdate,
		Delete: ResourceOutscaleVPNConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
			"connection_type": {
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
			"vpn_connection_id": {
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

func ResourceOutscaleVPNConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.CreateVpnConnectionRequest{
		ClientGatewayId:  d.Get("client_gateway_id").(string),
		VirtualGatewayId: d.Get("virtual_gateway_id").(string),
		ConnectionType:   d.Get("connection_type").(string),
	}

	if staticRoutesOnly, ok := d.GetOkExists("static_routes_only"); ok {
		req.SetStaticRoutesOnly(cast.ToBool(staticRoutesOnly))
	}
	var resp oscgo.CreateVpnConnectionResponse
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.VpnConnectionApi.CreateVpnConnection(context.Background()).CreateVpnConnectionRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating Outscale VPN Conecction: %s", err)
	}

	d.SetId(*resp.GetVpnConnection().VpnConnectionId)

	err = createOAPITagsSDK(conn, d)
	if err != nil {
		return err
	}

	return ResourceOutscaleVPNConnectionRead(d, meta)
}

func ResourceOutscaleVPNConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	vpnConnectionID := d.Id()

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available", "failed"},
		Refresh:    vpnConnectionRefreshFunc(conn, &vpnConnectionID),
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	r, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Outscale VPN Connection(%s) to become ready: %s", vpnConnectionID, err)
	}

	resp := r.(oscgo.ReadVpnConnectionsResponse)
	if utils.IsResponseEmpty(len(resp.GetVpnConnections()), "VpnConnection", d.Id()) {
		d.SetId("")
		return nil
	}
	vpnConnection := resp.GetVpnConnections()[0]

	if err := d.Set("client_gateway_configuration", vpnConnection.GetClientGatewayConfiguration()); err != nil {
		return err
	}
	if err := d.Set("vpn_connection_id", vpnConnection.GetVpnConnectionId()); err != nil {
		return err
	}
	if err := d.Set("state", vpnConnection.GetState()); err != nil {
		return err
	}
	if err := d.Set("static_routes_only", vpnConnection.GetStaticRoutesOnly()); err != nil {
		return err
	}
	if err := d.Set("client_gateway_id", vpnConnection.GetClientGatewayId()); err != nil {
		return err
	}
	if err := d.Set("virtual_gateway_id", vpnConnection.GetVirtualGatewayId()); err != nil {
		return err
	}
	if err := d.Set("connection_type", vpnConnection.GetConnectionType()); err != nil {
		return err
	}
	if err := d.Set("routes", flattenVPNConnection(vpnConnection.GetRoutes())); err != nil {
		return err
	}
	if err := d.Set("tags", flattenOAPITagsSDK(vpnConnection.GetTags())); err != nil {
		return err
	}
	if err := d.Set("vgw_telemetries", flattenVgwTelemetries(vpnConnection.GetVgwTelemetries())); err != nil {
		return err
	}
	return nil
}

func ResourceOutscaleVPNConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if err := updateOAPITagsSDK(conn, d); err != nil {
		return err
	}
	return ResourceOutscaleVPNConnectionRead(d, meta)
}

func ResourceOutscaleVPNConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	vpnConnectionID := d.Id()

	req := oscgo.DeleteVpnConnectionRequest{
		VpnConnectionId: vpnConnectionID,
	}
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, httpResp, err := conn.VpnConnectionApi.DeleteVpnConnection(context.Background()).DeleteVpnConnectionRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{"deleted", "failed"},
		Refresh:    vpnConnectionRefreshFunc(conn, &vpnConnectionID),
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Outscale VPN Connection(%s) to become deleted: %s", vpnConnectionID, err)
	}

	return nil
}

func vpnConnectionRefreshFunc(conn *oscgo.APIClient, vpnConnectionID *string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {

		filter := oscgo.ReadVpnConnectionsRequest{
			Filters: &oscgo.FiltersVpnConnection{
				VpnConnectionIds: &[]string{*vpnConnectionID},
			},
		}
		resp, httpResp, err := conn.VpnConnectionApi.ReadVpnConnections(context.Background()).ReadVpnConnectionsRequest(filter).Execute()
		if err != nil {
			if httpResp != nil {
				switch {
				case httpResp.StatusCode == http.StatusServiceUnavailable:
					return nil, "pending", nil
				case httpResp.StatusCode == http.StatusNotFound:
					return nil, "deleted", nil
				default:
					return nil, "failed", fmt.Errorf("Error on vpnConnectionRefresh: %s", err)
				}
			} else {
				return nil, "failed", fmt.Errorf("Error on vpnConnectionRefresh: %s", err)
			}
		}

		if len(resp.GetVpnConnections()) == 0 {
			return nil, "failed", fmt.Errorf("error on vpnConnectionRefresh: there are not vpn connections(%s)", *vpnConnectionID)
		}

		vpnConnection := resp.GetVpnConnections()[0]

		return resp, vpnConnection.GetState(), nil
	}
}

func flattenVPNConnection(routes []oscgo.RouteLight) []map[string]interface{} {
	routesMap := make([]map[string]interface{}, len(routes))

	for i, route := range routes {
		routesMap[i] = map[string]interface{}{
			"destination_ip_range": route.GetDestinationIpRange(),
			"route_type":           route.GetRouteType(),
			"state":                route.GetState(),
		}
	}
	return routesMap
}

func flattenVgwTelemetries(vgwTelemetries []oscgo.VgwTelemetry) []map[string]interface{} {
	vgwTelemetriesMap := make([]map[string]interface{}, len(vgwTelemetries))

	for i, vgwTelemetry := range vgwTelemetries {
		vgwTelemetriesMap[i] = map[string]interface{}{
			"accepted_route_count":   vgwTelemetry.GetAcceptedRouteCount(),
			"last_state_change_date": vgwTelemetry.GetLastStateChangeDate(),
			"outside_ip_address":     vgwTelemetry.GetOutsideIpAddress(),
			"state":                  vgwTelemetry.GetState(),
			"state_description":      vgwTelemetry.GetStateDescription(),
		}
	}
	return vgwTelemetriesMap
}
