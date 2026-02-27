package oapi

import (
	"context"
	"fmt"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
	"github.com/spf13/cast"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
)

func ResourceOutscaleClientGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleClientGatewayCreate,
		ReadContext:   ResourceOutscaleClientGatewayRead,
		UpdateContext: ResourceOutscaleClientGatewayUpdate,
		DeleteContext: ResourceOutscaleClientGatewayDelete,
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
			"bgp_asn": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"connection_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_gateway_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": TagsSchemaSDK(),
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceOutscaleClientGatewayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutCreate)

	req := osc.CreateClientGatewayRequest{
		BgpAsn:         cast.ToInt(d.Get("bgp_asn")),
		ConnectionType: d.Get("connection_type").(string),
		PublicIp:       d.Get("public_ip").(string),
	}

	resp, err := client.CreateClientGateway(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(*ptr.From(resp.ClientGateway).ClientGatewayId)
	err = createOAPITagsSDK(ctx, client, timeout, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceOutscaleClientGatewayRead(ctx, d, meta)
}

func ResourceOutscaleClientGatewayRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	clientGatewayID := d.Id()

	stateConf := &retry.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"available", "failed", "deleted"},
		Timeout: timeout,
		Refresh: clientGatewayRefreshFunc(ctx, client, timeout, &clientGatewayID),
	}

	r, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for outscale client gateway (%s) to become ready: %s", clientGatewayID, err)
	}

	resp := r.(*osc.ReadClientGatewaysResponse)
	if resp.ClientGateways == nil || utils.IsResponseEmpty(len(*resp.ClientGateways), "ClientGateway", d.Id()) ||
		ptr.From((*resp.ClientGateways)[0].State) == "deleted" {
		d.SetId("")
		return nil
	}

	clientGateway := (*resp.ClientGateways)[0]

	if err := d.Set("bgp_asn", ptr.From(clientGateway.BgpAsn)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("connection_type", ptr.From(clientGateway.ConnectionType)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("public_ip", ptr.From(clientGateway.PublicIp)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("client_gateway_id", ptr.From(clientGateway.ClientGatewayId)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", ptr.From(clientGateway.State)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", FlattenOAPITagsSDK(ptr.From(clientGateway.Tags))); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceOutscaleClientGatewayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutUpdate)

	if err := updateOAPITagsSDK(ctx, client, timeout, d); err != nil {
		return diag.FromErr(err)
	}

	return ResourceOutscaleClientGatewayRead(ctx, d, meta)
}

func ResourceOutscaleClientGatewayDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutDelete)

	gatewayID := d.Id()
	req := osc.DeleteClientGatewayRequest{
		ClientGatewayId: gatewayID,
	}

	_, err := client.DeleteClientGateway(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted", "failed"},
		Timeout: timeout,
		Refresh: clientGatewayRefreshFunc(ctx, client, timeout, &gatewayID),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for outscale client gateway (%s) to become deleted: %s", gatewayID, err)
	}

	return nil
}

func clientGatewayRefreshFunc(ctx context.Context, client *osc.Client, timeout time.Duration, gatewayID *string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := osc.ReadClientGatewaysRequest{
			Filters: &osc.FiltersClientGateway{
				ClientGatewayIds: &[]string{*gatewayID},
			},
		}
		resp, err := client.ReadClientGateways(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil || len(ptr.From(resp.ClientGateways)) == 0 {
			switch {
			case osc.IsConflict(err):
				return nil, "pending", nil
			default:
				return nil, "failed", fmt.Errorf("error on clientgatewayrefresh: %s", err)
			}
		}

		gateway := (*resp.ClientGateways)[0]

		return resp, *gateway.State, nil
	}
}
