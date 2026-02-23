package oapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/services/oapi/oapihelpers"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func ResourceOutscaleVPNConnectionRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleVPNConnectionRouteCreate,
		ReadContext:   ResourceOutscaleVPNConnectionRouteRead,
		DeleteContext: ResourceOutscaleVPNConnectionRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: ResourceOutscaleVPNConnectionRouteImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
		},

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

func ResourceOutscaleVPNConnectionRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutCreate)

	destinationIPRange := d.Get("destination_ip_range").(string)
	id := d.Get("vpn_connection_id").(string)

	req := osc.CreateVpnConnectionRouteRequest{
		DestinationIpRange: destinationIPRange,
		VpnConnectionId:    id,
	}
	_, err := client.CreateVpnConnectionRoute(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error creating outscale vpn conecction route: %s", err)
	}

	d.SetId(fmt.Sprintf("%s:%s", destinationIPRange, id))

	return ResourceOutscaleVPNConnectionRouteRead(ctx, d, meta)
}

func ResourceOutscaleVPNConnectionRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	destinationIPRange, vpnconnectionID := oapihelpers.ParseVPNConnectionRouteID(d.Id())

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available", "failed"},
		Timeout: timeout,
		Refresh: vpnconnectionRouteRefreshFunc(ctx, client, timeout, &destinationIPRange, &vpnconnectionID),
	}

	val, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for outscale vpn connection route(%s) to become ready: %s", d.Id(), err)
	}
	if val == nil {
		utils.LogManuallyDeleted("VpnconnectionRoute", d.Id())
		d.SetId("")
		return nil
	}

	return nil
}

func ResourceOutscaleVPNConnectionRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutDelete)

	destinationIPRange, id := oapihelpers.ParseVPNConnectionRouteID(d.Id())

	req := osc.DeleteVpnConnectionRouteRequest{
		DestinationIpRange: destinationIPRange,
		VpnConnectionId:    id,
	}
	_, err := client.DeleteVpnConnectionRoute(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted", "failed"},
		Timeout: timeout,
		Refresh: vpnconnectionRouteRefreshFunc(ctx, client, timeout, &destinationIPRange, &id),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for outscale vpn connection route(%s) to become deleted: %s", id, err)
	}

	return nil
}

func vpnconnectionRouteRefreshFunc(ctx context.Context, client *osc.Client, timeout time.Duration, destinationIPRange, id *string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		filter := osc.ReadVpnConnectionsRequest{
			Filters: &osc.FiltersVpnConnection{
				RouteDestinationIpRanges: &[]string{*destinationIPRange},
				VpnConnectionIds:         &[]string{*id},
			},
		}

		resp, err := client.ReadVpnConnections(ctx, filter, options.WithRetryTimeout(timeout))
		if err != nil {
			return nil, "failed", fmt.Errorf("error on vpnconnectionrouterefresh: %s", err)
		}

		if resp.VpnConnections == nil || len(*resp.VpnConnections) == 0 {
			return nil, "failed", nil
		}
		vpnConnection := (*resp.VpnConnections)[0]

		if vpnConnection.Routes != nil {
			for _, route := range *vpnConnection.Routes {
				if route.DestinationIpRange == *destinationIPRange {
					return resp, route.State, nil
				}
			}
		}

		return resp, "pending", nil
	}
}

func ResourceOutscaleVPNConnectionRouteImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	parts := strings.SplitN(d.Id(), "_", 2)
	if len(parts) != 2 {
		return nil, errors.New("import format error: to import a Outscale VPN connection Route, use the format {vpn_connection_id}_{destination_ip_range}")
	}

	vpnconnectionID := parts[0]
	destinationIPRange := parts[1]

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available", "failed"},
		Timeout: timeout,
		Refresh: vpnconnectionRouteRefreshFunc(ctx, client, timeout, &destinationIPRange, &vpnconnectionID),
	}

	val, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for outscale vpn connection route import(%s) to become ready: %s", d.Id(), err)
	}
	if val == nil {
		log.Printf("[WARN] VPN connection route %q could not be found. Removing Route from state.", vpnconnectionID)
		return nil, err
	}

	if err := d.Set("vpn_connection_id", vpnconnectionID); err != nil {
		return nil, fmt.Errorf("error setting `%s` for outscale vpn connection route(%s): %s", "vpn_connection_id", vpnconnectionID, err)
	}
	if err := d.Set("destination_ip_range", destinationIPRange); err != nil {
		return nil, fmt.Errorf("error setting `%s` for outscale vpn connection route(%s): %s", "destination_ip_range", destinationIPRange, err)
	}

	d.SetId(fmt.Sprintf("%s:%s", destinationIPRange, vpnconnectionID))

	return []*schema.ResourceData{d}, nil
}
