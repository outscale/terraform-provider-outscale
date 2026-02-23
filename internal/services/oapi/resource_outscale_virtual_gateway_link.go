package oapi

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleVirtualGatewayLink() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleVirtualGatewayLinkCreate,
		ReadContext:   ResourceOutscaleVirtualGatewayLinkRead,
		DeleteContext: ResourceOutscaleVirtualGatewayLinkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
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

func ResourceOutscaleVirtualGatewayLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutCreate)

	netID := d.Get("net_id").(string)
	vgwID := d.Get("virtual_gateway_id").(string)

	createOpts := osc.LinkVirtualGatewayRequest{
		NetId:            netID,
		VirtualGatewayId: vgwID,
	}
	log.Printf("[DEBUG] VPN Gateway attachment options: %#v", createOpts)

	_, err := client.LinkVirtualGateway(ctx, createOpts, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error attaching virtual gateway %q to vpc %q: %s",
			vgwID, netID, err)
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"detached", "attaching"},
		Target:  []string{"attached"},
		Timeout: timeout,
		Refresh: vpnGatewayLinkStateRefresh(ctx, client, netID, vgwID, timeout),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for virtual gateway %q to attach to vpc %q: %s",
			vgwID, netID, err)
	}
	log.Printf("[DEBUG] Virtual Gateway %q attached to VPC %q.", vgwID, netID)

	d.SetId(vgwID)

	return ResourceOutscaleVirtualGatewayLinkRead(ctx, d, meta)
}

func ResourceOutscaleVirtualGatewayLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutRead)

	vgwID := d.Id()

	resp, err := client.ReadVirtualGateways(ctx, osc.ReadVirtualGatewaysRequest{
		Filters: &osc.FiltersVirtualGateway{VirtualGatewayIds: &[]string{vgwID}},
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.VirtualGateways == nil || utils.IsResponseEmpty(len(*resp.VirtualGateways), "VirtualGateway", d.Id()) {
		d.SetId("")
		return nil
	}
	vgw := (*resp.VirtualGateways)[0]
	if ptr.From(vgw.State) == "deleted" {
		log.Printf("[INFO] VPN Gateway %q appears to have been deleted.", vgwID)
		d.SetId("")
		return nil
	}

	vga := oapiVpnGatewayGetLink(vgw)
	if len(ptr.From(vgw.NetToVirtualGatewayLinks)) == 0 || *vga.State == "detached" {
		// d.Set("net_id", "")
		return nil
	}

	if err := d.Set("net_id", ptr.From(vga.NetId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("virtual_gateway_id", ptr.From(vgw.VirtualGatewayId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("net_to_virtual_gateway_links", flattenNetToVirtualGatewayLinks(vgw.NetToVirtualGatewayLinks)); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenNetToVirtualGatewayLinks(netToVirtualGatewayLinks *[]osc.NetToVirtualGatewayLink) []map[string]interface{} {
	res := make([]map[string]interface{}, len(*netToVirtualGatewayLinks))

	if len(*netToVirtualGatewayLinks) > 0 {
		for i, n := range *netToVirtualGatewayLinks {
			res[i] = map[string]interface{}{
				"state":  n.State,
				"net_id": n.NetId,
			}
		}
	}
	return res
}

func ResourceOutscaleVirtualGatewayLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutDelete)

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

	_, err := client.UnlinkVirtualGateway(ctx, osc.UnlinkVirtualGatewayRequest{
		VirtualGatewayId: d.Id(),
		NetId:            netID.(string),
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for it to be fully detached before continuing
	log.Printf("[DEBUG] Waiting for VPN gateway (%s) to detach", d.Get("virtual_gateway_id").(string))
	stateConf := &retry.StateChangeConf{
		Pending: []string{"attached", "detaching", "available"},
		Target:  []string{"detached"},
		Timeout: timeout,
		Refresh: vpnGatewayAttachStateRefreshFunc(ctx, client, d.Get("virtual_gateway_id").(string), "detached", timeout),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf(
			"error waiting for vpn gateway (%s) to detach: %s",
			d.Get("virtual_gateway_id").(string), err)
	}

	return nil
}

func vpnGatewayLinkStateRefresh(ctx context.Context, client *osc.Client, vpcID, vgwID string, timeout time.Duration) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.ReadVirtualGateways(ctx, osc.ReadVirtualGatewaysRequest{Filters: &osc.FiltersVirtualGateway{
			VirtualGatewayIds: &[]string{vgwID},
			LinkNetIds:        &[]string{vpcID},
		}}, options.WithRetryTimeout(timeout))
		if err != nil {
			return nil, "", err
		}

		vgw := ptr.From(resp.VirtualGateways)[0]
		if len(ptr.From(vgw.NetToVirtualGatewayLinks)) == 0 {
			return vgw, "detached", nil
		}

		vga := oapiVpnGatewayGetLink(vgw)

		log.Printf("[DEBUG] VPN Gateway %q attachment status: %s", vgwID, *vga.State)
		return vgw, *vga.State, nil
	}
}
