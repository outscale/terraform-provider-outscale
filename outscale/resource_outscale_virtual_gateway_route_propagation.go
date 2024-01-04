package outscale

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
)

func resourceOutscaleOAPIVirtualGatewayRoutePropagation() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIVpnGatewayRoutePropagationEnable,
		Read:   resourceOutscaleOAPIVpnGatewayRoutePropagationRead,
		Delete: resourceOutscaleOAPIVpnGatewayRoutePropagationDisable,

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

func resourceOutscaleOAPIVpnGatewayRoutePropagationEnable(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	gwID := d.Get("virtual_gateway_id").(string)
	rtID := d.Get("route_table_id").(string)
	enable := d.Get("enable").(bool)

	log.Printf("\n\n[INFO] Enabling virtual gateway route propagation from %s to %s", gwID, rtID)

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.VirtualGatewayApi.UpdateRoutePropagation(context.Background()).UpdateRoutePropagationRequest(oscgo.UpdateRoutePropagationRequest{
			VirtualGatewayId: gwID,
			RouteTableId:     rtID,
			Enable:           enable,
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error enabling VGW propagation: %s", err)
	}

	d.SetId(fmt.Sprintf("%s_%s", gwID, rtID))

	return nil
}

func resourceOutscaleOAPIVpnGatewayRoutePropagationDisable(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	gwID := d.Get("virtual_gateway_id").(string)
	rtID := d.Get("route_table_id").(string)
	enable := d.Get("enable").(bool)

	log.Printf("\n\n[INFO] Disabling VGW propagation from %s to %s", gwID, rtID)

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, httpResp, err := conn.VirtualGatewayApi.UpdateRoutePropagation(context.Background()).UpdateRoutePropagationRequest(oscgo.UpdateRoutePropagationRequest{
			VirtualGatewayId: gwID,
			RouteTableId:     rtID,
			Enable:           enable,
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error disabling VGW propagation: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceOutscaleOAPIVpnGatewayRoutePropagationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	gwID := d.Get("virtual_gateway_id").(string)
	rtID := d.Get("route_table_id").(string)

	var resp oscgo.ReadRouteTablesResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.RouteTableApi.ReadRouteTables(context.Background()).ReadRouteTablesRequest(oscgo.ReadRouteTablesRequest{
			Filters: &oscgo.FiltersRouteTable{RouteTableIds: &[]string{rtID}},
		}).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return err
	}
	if utils.IsResponseEmpty(len(resp.GetRouteTables()), "VirtualGatewayRoutePropagation", d.Id()) {
		d.SetId("")
		return nil
	}
	rt := resp.GetRouteTables()[0]

	exists := false
	for _, vgw := range rt.GetRoutePropagatingVirtualGateways() {
		if vgw.GetVirtualGatewayId() == gwID {
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
