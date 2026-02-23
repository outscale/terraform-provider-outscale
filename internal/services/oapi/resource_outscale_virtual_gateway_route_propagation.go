package oapi

import (
	"context"
	"fmt"
	"log"

	"github.com/outscale/goutils/sdk/ptr"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"
)

func ResourceOutscaleVirtualGatewayRoutePropagation() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleVpnGatewayRoutePropagationEnable,
		ReadContext:   ResourceOutscaleVpnGatewayRoutePropagationRead,
		DeleteContext: ResourceOutscaleVpnGatewayRoutePropagationDisable,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
		},
		Schema: map[string]*schema.Schema{
			"virtual_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enable": {
				Type:     schema.TypeBool,
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

func ResourceOutscaleVpnGatewayRoutePropagationEnable(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutCreate)

	gwID := d.Get("virtual_gateway_id").(string)
	rtID := d.Get("route_table_id").(string)
	enable := d.Get("enable").(bool)

	log.Printf("\n\n[INFO] Enabling virtual gateway route propagation from %s to %s", gwID, rtID)

	_, err := client.UpdateRoutePropagation(ctx, osc.UpdateRoutePropagationRequest{
		VirtualGatewayId: gwID,
		RouteTableId:     rtID,
		Enable:           enable,
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error enabling vgw propagation: %s", err)
	}

	d.SetId(fmt.Sprintf("%s_%s", gwID, rtID))

	return nil
}

func ResourceOutscaleVpnGatewayRoutePropagationDisable(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutDelete)

	gwID := d.Get("virtual_gateway_id").(string)
	rtID := d.Get("route_table_id").(string)
	enable := d.Get("enable").(bool)

	log.Printf("\n\n[INFO] Disabling VGW propagation from %s to %s", gwID, rtID)

	_, err := client.UpdateRoutePropagation(ctx, osc.UpdateRoutePropagationRequest{
		VirtualGatewayId: gwID,
		RouteTableId:     rtID,
		Enable:           enable,
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error disabling vgw propagation: %s", err)
	}
	d.SetId("")

	return nil
}

func ResourceOutscaleVpnGatewayRoutePropagationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutRead)

	gwID := d.Get("virtual_gateway_id").(string)
	rtID := d.Get("route_table_id").(string)

	resp, err := client.ReadRouteTables(ctx, osc.ReadRouteTablesRequest{
		Filters: &osc.FiltersRouteTable{RouteTableIds: &[]string{rtID}},
	}, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.RouteTables == nil || utils.IsResponseEmpty(len(*resp.RouteTables), "VirtualGatewayRoutePropagation", d.Id()) {
		d.SetId("")
		return nil
	}
	rt := (*resp.RouteTables)[0]

	exists := false
	for _, vgw := range rt.RoutePropagatingVirtualGateways {
		if ptr.From(vgw.VirtualGatewayId) == gwID {
			exists = true
		}
	}
	if !exists {
		log.Printf("\n\n[INFO] %s is no longer propagating to %s, so dropping route propagation from state", rtID, gwID)
		d.SetId("")
		return nil
	}

	return nil
}
