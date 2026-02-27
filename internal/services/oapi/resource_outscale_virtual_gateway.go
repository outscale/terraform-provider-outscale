package oapi

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func ResourceOutscaleVirtualGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleVirtualGatewayCreate,
		ReadContext:   ResourceOutscaleVirtualGatewayRead,
		UpdateContext: ResourceOutscaleVirtualGatewayUpdate,
		DeleteContext: ResourceOutscaleVirtualGatewayDelete,
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
			"connection_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"net_to_virtual_gateway_links": {
				Type:     schema.TypeList,
				Optional: true,
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
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"virtual_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": TagsSchemaSDK(),
		},
	}
}

func ResourceOutscaleVirtualGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)
	clientType, clientecTypeOk := d.GetOk("connection_type")
	createOpts := osc.CreateVirtualGatewayRequest{}
	if clientecTypeOk {
		createOpts.ConnectionType = clientType.(string)
	}

	resp, err := client.CreateVirtualGateway(ctx, createOpts, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error creating vpn gateway: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available"},
		Timeout: timeout,
		Refresh: virtualGatewayStateRefreshFunc(ctx, client, *resp.VirtualGateway.VirtualGatewayId, "deleted", timeout),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for instance (%s) to become created: %s", d.Id(), err)
	}

	virtualGateway := resp.VirtualGateway
	d.SetId(ptr.From(virtualGateway.VirtualGatewayId))

	if d.IsNewResource() {
		if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
			return diag.FromErr(err)
		}
	}
	return ResourceOutscaleVirtualGatewayRead(ctx, d, meta)
}

func ResourceOutscaleVirtualGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	resp, err := client.ReadVirtualGateways(ctx, osc.ReadVirtualGatewaysRequest{
		Filters: &osc.FiltersVirtualGateway{VirtualGatewayIds: &[]string{d.Id()}},
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		fmt.Printf("\n\n[ERROR] Error finding VpnGateway: %s", err)
		return diag.FromErr(err)
	}

	if resp.VirtualGateways == nil || utils.IsResponseEmpty(len(*resp.VirtualGateways), "VirtualGateway", d.Id()) {
		d.SetId("")
		return nil
	}
	virtualGateway := (*resp.VirtualGateways)[0]
	if ptr.From(virtualGateway.State) == "deleted" {
		d.SetId("")
		return nil
	}
	if virtualGateway.NetToVirtualGatewayLinks != nil {
		vs := make([]map[string]interface{}, len(*virtualGateway.NetToVirtualGatewayLinks))
		for k, v := range *virtualGateway.NetToVirtualGatewayLinks {
			vp := make(map[string]interface{})
			vp["state"] = v.State
			vp["net_id"] = v.NetId
			vs[k] = vp
		}
		d.Set("net_to_virtual_gateway_links", vs)
	}

	d.Set("connection_type", ptr.From(virtualGateway.ConnectionType))
	d.Set("virtual_gateway_id", ptr.From(virtualGateway.VirtualGatewayId))

	d.Set("state", ptr.From(virtualGateway.State))
	d.Set("tags", FlattenOAPITagsSDK(ptr.From(virtualGateway.Tags)))

	return nil
}

func ResourceOutscaleVirtualGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutUpdate)

	if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func ResourceOutscaleVirtualGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutDelete)

	_, err := client.DeleteVirtualGateway(ctx, osc.DeleteVirtualGatewayRequest{VirtualGatewayId: d.Id()}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")

	return nil
}

// vpnGatewayAttachStateRefreshFunc returns a retry.StateRefreshFunc that is used to watch
// the state of a VPN gateway's attachment
func vpnGatewayAttachStateRefreshFunc(ctx context.Context, client *osc.Client, id string, expected string, timeout time.Duration) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.ReadVirtualGateways(ctx, osc.ReadVirtualGatewaysRequest{
			Filters: &osc.FiltersVirtualGateway{VirtualGatewayIds: &[]string{id}},
		}, options.WithRetryTimeout(timeout))
		if err != nil {
			fmt.Printf("[ERROR] Error on VpnGatewayStateRefresh: %s", err)
			return nil, "", err
		}

		if resp.VirtualGateways == nil {
			return nil, "", nil
		}

		virtualGateway := (*resp.VirtualGateways)[0]
		if virtualGateway.NetToVirtualGatewayLinks == nil || len(*virtualGateway.NetToVirtualGatewayLinks) == 0 {
			return virtualGateway, "detached", nil
		}
		vpnAttachment := oapiVpnGatewayGetLink(virtualGateway)

		return virtualGateway, *vpnAttachment.State, nil
	}
}

func oapiVpnGatewayGetLink(vgw osc.VirtualGateway) *osc.NetToVirtualGatewayLink {
	for _, v := range ptr.From(vgw.NetToVirtualGatewayLinks) {
		if ptr.From(v.State) == "attached" {
			return &v
		}
	}
	return &osc.NetToVirtualGatewayLink{State: new("detached")}
}

func virtualGatewayStateRefreshFunc(ctx context.Context, client *osc.Client, instanceID, failState string, timeout time.Duration) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.ReadVirtualGateways(ctx, osc.ReadVirtualGatewaysRequest{
			Filters: &osc.FiltersVirtualGateway{
				VirtualGatewayIds: &[]string{instanceID},
			},
		}, options.WithRetryTimeout(timeout))
		if err != nil {
			log.Printf("[ERROR] error on InstanceStateRefresh: %s", err)
			return nil, "", err
		}

		if resp.VirtualGateways == nil {
			return nil, "", nil
		}

		virtualGateway := (*resp.VirtualGateways)[0]
		state := *virtualGateway.State

		if state == failState {
			return virtualGateway, state, fmt.Errorf("failed to reach target state:: %v", *virtualGateway.State)
		}

		return virtualGateway, state, nil
	}
}
