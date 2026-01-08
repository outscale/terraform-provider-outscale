package oapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleVirtualGatewayLink() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleVirtualGatewayLinkCreate,
		Read:   ResourceOutscaleVirtualGatewayLinkRead,
		Delete: ResourceOutscaleVirtualGatewayLinkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"net_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"virtual_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dry_run": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
		},
	}
}

func ResourceOutscaleVirtualGatewayLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	netID := d.Get("net_id").(string)
	vgwID := d.Get("virtual_gateway_id").(string)

	createOpts := oscgo.LinkVirtualGatewayRequest{
		NetId:            netID,
		VirtualGatewayId: vgwID,
	}
	log.Printf("[DEBUG] VPN Gateway attachment options: %#v", createOpts)

	var err error

	err = retry.Retry(30*time.Second, func() *retry.RetryError {
		_, httpResp, err := conn.VirtualGatewayApi.LinkVirtualGateway(context.Background()).LinkVirtualGatewayRequest(createOpts).Execute()
		if err != nil {
			if httpResp.StatusCode == http.StatusNotFound {
				return retry.RetryableError(
					fmt.Errorf("gateway not found, retry for eventual consistancy"))
			}
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error attaching virtual gateway %q to vpc %q: %s",
			vgwID, netID, err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"detached", "attaching"},
		Target:     []string{"attached"},
		Refresh:    vpnGatewayLinkStateRefresh(conn, netID, vgwID),
		Timeout:    15 * time.Minute,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for virtual gateway %q to attach to vpc %q: %s",
			vgwID, netID, err)
	}
	log.Printf("[DEBUG] Virtual Gateway %q attached to VPC %q.", vgwID, netID)

	d.SetId(vgwID)

	return ResourceOutscaleVirtualGatewayLinkRead(d, meta)
}

func ResourceOutscaleVirtualGatewayLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	vgwID := d.Id()

	var resp oscgo.ReadVirtualGatewaysResponse
	var err error
	var statusCode int
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.VirtualGatewayApi.ReadVirtualGateways(context.Background()).ReadVirtualGatewaysRequest(oscgo.ReadVirtualGatewaysRequest{
			Filters: &oscgo.FiltersVirtualGateway{VirtualGatewayIds: &[]string{vgwID}},
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		statusCode = httpResp.StatusCode
		return nil
	})
	if err != nil {
		if statusCode == http.StatusNotFound {
			log.Printf("[WARN] VPN Gateway %q not found.", vgwID)
			d.SetId("")
			return nil
		}
		return err
	}
	if utils.IsResponseEmpty(len(resp.GetVirtualGateways()), "VirtualGateway", d.Id()) {
		d.SetId("")
		return nil
	}
	vgw := resp.GetVirtualGateways()[0]
	if vgw.GetState() == "deleted" {
		log.Printf("[INFO] VPN Gateway %q appears to have been deleted.", vgwID)
		d.SetId("")
		return nil
	}

	vga := oapiVpnGatewayGetLink(vgw)
	if len(vgw.GetNetToVirtualGatewayLinks()) == 0 || vga.GetState() == "detached" {
		// d.Set("net_id", "")
		return nil
	}

	if err := d.Set("net_id", vga.GetNetId()); err != nil {
		return err
	}
	if err := d.Set("virtual_gateway_id", vgw.GetVirtualGatewayId()); err != nil {
		return err
	}
	if err := d.Set("net_to_virtual_gateway_links", flattenNetToVirtualGatewayLinks(vgw.NetToVirtualGatewayLinks)); err != nil {
		return err
	}
	return nil
}

func flattenNetToVirtualGatewayLinks(netToVirtualGatewayLinks *[]oscgo.NetToVirtualGatewayLink) []map[string]interface{} {
	res := make([]map[string]interface{}, len(*netToVirtualGatewayLinks))

	if len(*netToVirtualGatewayLinks) > 0 {
		for i, n := range *netToVirtualGatewayLinks {
			res[i] = map[string]interface{}{
				"state":  n.GetState(),
				"net_id": n.GetNetId(),
			}
		}
	}
	return res
}

func ResourceOutscaleVirtualGatewayLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	// Get the old VPC ID to detach from
	netID, _ := d.GetChange("net_id")

	if netID.(string) == "" {
		fmt.Printf(
			"[DEBUG] Not detaching Virtual Gateway '%s' as no VPC ID is set",
			d.Get("virtual_gateway_id").(string))
		return nil
	}

	fmt.Printf(
		"[INFO] Detaching Virtual Gateway '%s' from VPC '%s'",
		d.Get("virtual_gateway_id").(string),
		netID.(string))

	wait := true

	var err error
	var statusCode int
	err = retry.Retry(30*time.Second, func() *retry.RetryError {
		_, httpResp, err := conn.VirtualGatewayApi.UnlinkVirtualGateway(context.Background()).UnlinkVirtualGatewayRequest(oscgo.UnlinkVirtualGatewayRequest{
			VirtualGatewayId: d.Id(),
			NetId:            netID.(string),
		}).Execute()
		if err != nil {
			if httpResp.StatusCode == http.StatusNotFound {
				return retry.RetryableError(
					fmt.Errorf("gateway not found, retry for eventual consistancy"))
			}
			return utils.CheckThrottling(httpResp, err)
		}
		statusCode = httpResp.StatusCode
		return nil
	})
	if err != nil {
		if statusCode == http.StatusNotFound {
			err = nil
			wait = false
		}
		if err != nil {
			return err
		}
	}

	if !wait {
		return nil
	}

	// Wait for it to be fully detached before continuing
	log.Printf("[DEBUG] Waiting for VPN gateway (%s) to detach", d.Get("virtual_gateway_id").(string))
	stateConf := &retry.StateChangeConf{
		Pending: []string{"attached", "detaching", "available"},
		Target:  []string{"detached"},
		Refresh: vpnGatewayAttachStateRefreshFunc(conn, d.Get("virtual_gateway_id").(string), "detached"),
		Timeout: 5 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"error waiting for vpn gateway (%s) to detach: %s",
			d.Get("virtual_gateway_id").(string), err)
	}

	return nil
}

func vpnGatewayLinkStateRefresh(conn *oscgo.APIClient, vpcID, vgwID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var err error
		var resp oscgo.ReadVirtualGatewaysResponse
		var statusCode int
		err = retry.Retry(30*time.Second, func() *retry.RetryError {
			rp, httpResp, err := conn.VirtualGatewayApi.ReadVirtualGateways(context.Background()).ReadVirtualGatewaysRequest(oscgo.ReadVirtualGatewaysRequest{Filters: &oscgo.FiltersVirtualGateway{
				VirtualGatewayIds: &[]string{vgwID},
				LinkNetIds:        &[]string{vpcID},
			}}).Execute()
			if err != nil {
				if httpResp.StatusCode == http.StatusNotFound {
					return retry.RetryableError(
						fmt.Errorf("gateway not found, retry for eventual consistancy"))
				}
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			statusCode = httpResp.StatusCode
			return nil
		})
		if err != nil {
			if statusCode == http.StatusNotFound {
				log.Printf("[WARN] VPN Gateway %q not found.", vgwID)
				return nil, "", nil
			}
			return nil, "", err
		}

		vgw := resp.GetVirtualGateways()[0]
		if len(vgw.GetNetToVirtualGatewayLinks()) == 0 {
			return vgw, "detached", nil
		}

		vga := oapiVpnGatewayGetLink(vgw)

		log.Printf("[DEBUG] VPN Gateway %q attachment status: %s", vgwID, *vga.State)
		return vgw, *vga.State, nil
	}
}
