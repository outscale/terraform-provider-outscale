package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/spf13/cast"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
)

func resourceOutscaleVPNConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleVPNConnectionCreate,
		Read:   resourceOutscaleVPNConnectionRead,
		Update: resourceOutscaleVPNConnectionUpdate,
		Delete: resourceOutscaleVPNConnectionDelete,
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
			"tags": tagsListOAPISchema(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleVPNConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	req := oscgo.CreateVpnConnectionRequest{
		ClientGatewayId:  d.Get("client_gateway_id").(string),
		VirtualGatewayId: d.Get("virtual_gateway_id").(string),
		ConnectionType:   d.Get("connection_type").(string),
	}

	if staticRoutesOnly, ok := d.GetOkExists("static_routes_only"); ok {
		req.SetStaticRoutesOnly(cast.ToBool(staticRoutesOnly))
	}

	vpn, _, err := conn.VpnConnectionApi.CreateVpnConnection(context.Background(),
		&oscgo.CreateVpnConnectionOpts{
			CreateVpnConnectionRequest: optional.NewInterface(req),
		})
	if err != nil {
		return fmt.Errorf("Error creating Outscale VPN Conecction: %s", err)
	}

	if tags, ok := d.GetOk("tags"); ok {
		err := assignTags(tags.([]interface{}), *vpn.GetVpnConnection().VpnConnectionId, conn)
		if err != nil {
			return err
		}
	}

	d.SetId(*vpn.GetVpnConnection().VpnConnectionId)

	return resourceOutscaleVPNConnectionRead(d, meta)
}

func resourceOutscaleVPNConnectionRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	vpnConnectionID := d.Id()

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available", "failed"},
		Refresh:    vpnConnectionRefreshFunc(conn, &vpnConnectionID),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	r, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Outscale VPN Connection(%s) to become ready: %s", vpnConnectionID, err)
	}

	resp := r.(oscgo.ReadVpnConnectionsResponse)
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
	if err := d.Set("tags", tagsOSCAPIToMap(vpnConnection.GetTags())); err != nil {
		return err
	}
	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return err
	}
	return nil
}

func resourceOutscaleVPNConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	d.Partial(true)

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	d.Partial(false)

	return resourceOutscaleVPNConnectionRead(d, meta)
}

func resourceOutscaleVPNConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	vpnConnectionID := d.Id()

	req := oscgo.DeleteVpnConnectionRequest{
		VpnConnectionId: vpnConnectionID,
	}

	_, _, err := conn.VpnConnectionApi.DeleteVpnConnection(context.Background(), &oscgo.DeleteVpnConnectionOpts{
		DeleteVpnConnectionRequest: optional.NewInterface(req),
	})
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{"deleted", "failed"},
		Refresh:    vpnConnectionRefreshFunc(conn, &vpnConnectionID),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Outscale VPN Connection(%s) to become deleted: %s", vpnConnectionID, err)
	}

	return nil
}

func vpnConnectionRefreshFunc(conn *oscgo.APIClient, vpnConnectionID *string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		filter := oscgo.ReadVpnConnectionsRequest{
			Filters: &oscgo.FiltersVpnConnection{
				VpnConnectionIds: &[]string{*vpnConnectionID},
			},
		}

		resp, _, err := conn.VpnConnectionApi.ReadVpnConnections(context.Background(), &oscgo.ReadVpnConnectionsOpts{
			ReadVpnConnectionsRequest: optional.NewInterface(filter),
		})
		if err != nil {
			switch {
			case strings.Contains(fmt.Sprint(err), "RequestLimitExceeded:"):
				return nil, "pending", nil
			case strings.Contains(fmt.Sprint(err), "404"):
				return nil, "deleted", nil
			default:
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
