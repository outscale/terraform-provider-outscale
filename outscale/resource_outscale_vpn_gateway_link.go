package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleOAPIVpnGatewayLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIVpnGatewayLinkCreate,
		Read:   resourceOutscaleOAPIVpnGatewayLinkRead,
		Delete: resourceOutscaleOAPIVpnGatewayLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"lin_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpn_gateway_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPIVpnGatewayLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	vgwID := d.Get("vpn_gateway_id").(string)

	var resp *fcu.DescribeVpnGatewaysOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeVpnGateways(&fcu.DescribeVpnGatewaysInput{
			VpnGatewayIds: []*string{aws.String(vgwID)},
		})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if ok && awsErr.Code() == "InvalidVPNGatewayID.NotFound" {
			log.Printf("[WARN] VPN Gateway %q not found.", vgwID)
			d.SetId("")
			return nil
		}
		return err
	}

	vgw := resp.VpnGateways[0]
	if *vgw.State == "deleted" {
		log.Printf("[INFO] VPN Gateway %q appears to have been deleted.", vgwID)
		d.SetId("")
		return nil
	}

	vga := oapiVpnGatewayGetAttachment(vgw)
	if len(vgw.VpcAttachments) == 0 || *vga.State == "detached" {
		d.Set("lin_id", "")
		return nil
	}

	d.Set("lin_id", *vga.VpcId)
	d.Set("state", vga.State)
	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceOutscaleOAPIVpnGatewayLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	vpcID := d.Get("lin_id").(string)
	vgwID := d.Get("vpn_gateway_id").(string)

	createOpts := &fcu.AttachVpnGatewayInput{
		VpcId:        aws.String(vpcID),
		VpnGatewayId: aws.String(vgwID),
	}
	log.Printf("[DEBUG] VPN Gateway attachment options: %#v", *createOpts)

	var err error

	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		_, err = conn.VM.AttachVpnGateway(createOpts)
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidVpnGatewayID.NotFound") {
				return resource.RetryableError(
					fmt.Errorf("Gateway not found, retry for eventual consistancy"))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error attaching VPN Gateway %q to VPC %q: %s",
			vgwID, vpcID, err)
	}

	d.SetId(vpnGatewayAttachmentID(vpcID, vgwID))

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"detached", "attaching"},
		Target:     []string{"attached"},
		Refresh:    vpnGatewayAttachmentStateRefresh(conn, vpcID, vgwID),
		Timeout:    15 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for VPN Gateway %q to attach to VPC %q: %s",
			vgwID, vpcID, err)
	}
	log.Printf("[DEBUG] VPN Gateway %q attached to VPC %q.", vgwID, vpcID)

	return resourceOutscaleOAPIVpnGatewayLinkRead(d, meta)
}

func resourceOutscaleOAPIVpnGatewayLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	// Get the old VPC ID to detach from
	vpcID, _ := d.GetChange("lin_id")

	if vpcID.(string) == "" {
		fmt.Printf(
			"[DEBUG] Not detaching VPN Gateway '%s' as no VPC ID is set",
			d.Get("vpn_gateway_id").(string))
		return nil
	}

	fmt.Printf(
		"[INFO] Detaching VPN Gateway '%s' from VPC '%s'",
		d.Get("vpn_gateway_id").(string),
		vpcID.(string))

	wait := true

	var err error
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		_, err = conn.VM.DetachVpnGateway(&fcu.DetachVpnGatewayInput{
			VpnGatewayId: aws.String(d.Get("vpn_gateway_id").(string)),
			VpcId:        aws.String(vpcID.(string)),
		})
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidVpnGatewayID.NotFound") {
				return resource.RetryableError(
					fmt.Errorf("Gateway not found, retry for eventual consistancy"))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidVpnGatewayID.NotFound") {
			err = nil
			wait = false
		} else if strings.Contains(fmt.Sprint(err), "InvalidVpnGatewayAttachment.NotFound") {
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
	log.Printf("[DEBUG] Waiting for VPN gateway (%s) to detach", d.Get("vpn_gateway_id").(string))
	stateConf := &resource.StateChangeConf{
		Pending: []string{"attached", "detaching", "available"},
		Target:  []string{"detached"},
		Refresh: vpnGatewayAttachStateRefreshFunc(conn, d.Get("vpn_gateway_id").(string), "detached"),
		Timeout: 5 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for vpn gateway (%s) to detach: %s",
			d.Get("vpn_gateway_id").(string), err)
	}

	return nil
}

func vpnGatewayAttachmentStateRefresh(conn *fcu.Client, vpcID, vgwID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		var err error
		var resp *fcu.DescribeVpnGatewaysOutput
		err = resource.Retry(30*time.Second, func() *resource.RetryError {
			resp, err = conn.VM.DescribeVpnGateways(&fcu.DescribeVpnGatewaysInput{
				Filters: []*fcu.Filter{
					&fcu.Filter{
						Name:   aws.String("attachment.vpc-id"),
						Values: []*string{aws.String(vpcID)},
					},
				},
				VpnGatewayIds: []*string{aws.String(vgwID)},
			})
			if err != nil {
				if strings.Contains(fmt.Sprint(err), "InvalidVpnGatewayID.NotFound") {
					return resource.RetryableError(
						fmt.Errorf("Gateway not found, retry for eventual consistancy"))
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})

		if err != nil {
			awsErr, ok := err.(awserr.Error)
			if ok {
				switch awsErr.Code() {
				case "InvalidVPNGatewayID.NotFound":
					fallthrough
				case "InvalidVpnGatewayAttachment.NotFound":
					return nil, "", nil
				}
			}

			return nil, "", err
		}

		vgw := resp.VpnGateways[0]
		if len(vgw.VpcAttachments) == 0 {
			return vgw, "detached", nil
		}

		vga := oapiVpnGatewayGetAttachment(vgw)

		log.Printf("[DEBUG] VPN Gateway %q attachment status: %s", vgwID, *vga.State)
		return vgw, *vga.State, nil
	}
}

func vpnGatewayAttachmentID(vpcID, vgwID string) string {
	return fmt.Sprintf("vpn-attachment-%x", hashcode.String(fmt.Sprintf("%s-%s", vpcID, vgwID)))
}
