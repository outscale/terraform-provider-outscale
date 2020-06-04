package outscale

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/antihax/optional"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	oscgo "github.com/marinsalinas/osc-sdk-go"
)

func resourceOutscaleVPNConnectionRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleVPNConnectionRouteCreate,
		Read:   resourceOutscaleVPNConnectionRouteRead,
		Delete: resourceOutscaleVPNConnectionRouteDelete,

		Schema: map[string]*schema.Schema{
			"destination_ip_range": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpn_connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleVPNConnectionRouteCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	destinationIPRange := d.Get("destination_ip_range").(string)
	vpnConnectionID := d.Get("vpn_connection_id").(string)

	req := oscgo.CreateVpnConnectionRouteRequest{
		DestinationIpRange: destinationIPRange,
		VpnConnectionId:    vpnConnectionID,
	}

	_, _, err := conn.VpnConnectionApi.CreateVpnConnectionRoute(context.Background(),
		&oscgo.CreateVpnConnectionRouteOpts{
			CreateVpnConnectionRouteRequest: optional.NewInterface(req),
		})
	if err != nil {
		return fmt.Errorf("Error creating Outscale VPN Conecction Route: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", destinationIPRange, vpnConnectionID))

	return resourceOutscaleVPNConnectionRouteRead(d, meta)
}

func resourceOutscaleVPNConnectionRouteRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	destinationIPRange, vpnConnectionID := resourceOutscaleVPNConnectionRouteParseID(d.Id())

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"pending"},
		Target:     []string{"available", "failed"},
		Refresh:    vpnConnectionRouteRefreshFunc(conn, &destinationIPRange, &vpnConnectionID),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	r, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Outscale VPN Connection Route(%s) to become ready: %s", d.Id(), err)
	}

	resp := r.(oscgo.ReadVpnConnectionsResponse)

	if err := d.Set("request_id", resp.ResponseContext.GetRequestId()); err != nil {
		return err
	}

	return nil
}

func resourceOutscaleVPNConnectionRouteDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	destinationIPRange, vpnConnectionID := resourceOutscaleVPNConnectionRouteParseID(d.Id())

	req := oscgo.DeleteVpnConnectionRouteRequest{
		DestinationIpRange: destinationIPRange,
		VpnConnectionId:    vpnConnectionID,
	}

	_, _, err := conn.VpnConnectionApi.DeleteVpnConnectionRoute(context.Background(), &oscgo.DeleteVpnConnectionRouteOpts{
		DeleteVpnConnectionRouteRequest: optional.NewInterface(req),
	})
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"deleting"},
		Target:     []string{"deleted", "failed"},
		Refresh:    vpnConnectionRouteRefreshFunc(conn, &destinationIPRange, &vpnConnectionID),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Outscale VPN Connection Route(%s) to become deleted: %s", vpnConnectionID, err)
	}

	return nil
}

func vpnConnectionRouteRefreshFunc(conn *oscgo.APIClient, destinationIPRange, vpnConnectionID *string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		filter := oscgo.ReadVpnConnectionsRequest{
			Filters: &oscgo.FiltersVpnConnection{
				RouteDestinationIpRanges: &[]string{*destinationIPRange},
				VpnConnectionIds:         &[]string{*vpnConnectionID},
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
				return nil, "failed", fmt.Errorf("Error on vpnConnectionRouteRefresh: %s", err)
			}
		}

		vpnConnection := resp.GetVpnConnections()[0]

		routes, ok := vpnConnection.GetRoutesOk()
		if ok {
			for _, route := range routes {
				if route.GetDestinationIpRange() == *destinationIPRange {
					return resp, route.GetState(), nil
				}
			}
		}

		return resp, "pending", nil
	}
}

func resourceOutscaleVPNConnectionRouteParseID(ID string) (string, string) {
	parts := strings.SplitN(ID, ":", 2)
	return parts[0], parts[1]
}
