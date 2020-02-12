package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOAPIVpnGatewayRoutePropagation() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIVpnGatewayRoutePropagationEnable,
		Read:   resourceOutscaleOAPIVpnGatewayRoutePropagationRead,
		Delete: resourceOutscaleOAPIVpnGatewayRoutePropagationDisable,

		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"route_table_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIVpnGatewayRoutePropagationEnable(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	gwID := d.Get("gateway_id").(string)
	rtID := d.Get("route_table_id").(string)

	log.Printf("\n\n[INFO] Enabling VGW propagation from %s to %s", gwID, rtID)

	var err error
	var resp *fcu.EnableVgwRoutePropagationOutput

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.EnableVgwRoutePropagation(&fcu.EnableVgwRoutePropagationInput{
			GatewayId:    aws.String(gwID),
			RouteTableId: aws.String(rtID),
		})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("error enabling VGW propagation: %s", err)
	}

	d.SetId(fmt.Sprintf("%s_%s", gwID, rtID))
	d.Set("request_id", *resp.RequestId)

	return nil
}

func resourceOutscaleOAPIVpnGatewayRoutePropagationDisable(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	gwID := d.Get("gateway_id").(string)
	rtID := d.Get("route_table_id").(string)

	log.Printf("\n\n[INFO] Disabling VGW propagation from %s to %s", gwID, rtID)

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.DisableVgwRoutePropagation(&fcu.DisableVgwRoutePropagationInput{
			GatewayId:    aws.String(gwID),
			RouteTableId: aws.String(rtID),
		})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
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
	conn := meta.(*OutscaleClient).FCU

	gwID := d.Get("gateway_id").(string)
	rtID := d.Get("route_table_id").(string)

	var resp *fcu.DescribeRouteTablesOutput
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeRouteTables(&fcu.DescribeRouteTablesInput{
			RouteTableIds: []*string{aws.String(rtID)},
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "RequestLimitExceeded") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	rt := resp.RouteTables[0]

	exists := false
	for _, vgw := range rt.PropagatingVgws {
		if *vgw.GatewayId == gwID {
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
