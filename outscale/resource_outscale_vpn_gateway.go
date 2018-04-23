package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleVpnGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleVpnGatewayCreate,
		Read:   resourceOutscaleVpnGatewayRead,
		Delete: resourceOutscaleVpnGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"attachments": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpn_gateway_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_set": tagsSchemaComputed(),
			"tag":     tagsSchema(),
		},
	}
}

func resourceOutscaleVpnGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	createOpts := &fcu.CreateVpnGatewayInput{
		Type: aws.String("ipsec.1"),
	}

	var resp *fcu.CreateVpnGatewayOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.CreateVpnGateway(createOpts)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
	if err != nil {
		return fmt.Errorf("Error creating VPN gateway: %s", err)
	}

	vpnGateway := resp.VpnGateway
	d.SetId(*vpnGateway.VpnGatewayId)

	if d.IsNewResource() {
		if err := setTags(conn, d); err != nil {
			return err
		}
		d.SetPartial("tag_set")
	}

	return resourceOutscaleVpnGatewayRead(d, meta)
}

func resourceOutscaleVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	var resp *fcu.DescribeVpnGatewaysOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, err = conn.VM.DescribeVpnGateways(&fcu.DescribeVpnGatewaysInput{
			VpnGatewayIds: []*string{aws.String(d.Id())},
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
		if strings.Contains(fmt.Sprint(err), "InvalidVpnGatewayID.NotFound") {
			d.SetId("")
			return nil
		} else {
			fmt.Printf("\n\n[ERROR] Error finding VpnGateway: %s", err)
			return err
		}
	}

	vpnGateway := resp.VpnGateways[0]
	if vpnGateway == nil || *vpnGateway.State == "deleted" {
		d.SetId("")
		return nil
	}

	vpnAttachment := vpnGatewayGetAttachment(vpnGateway)
	if len(vpnGateway.VpcAttachments) == 0 || *vpnAttachment.State == "detached" {
		d.Set("vpc_id", "")
	} else {
		d.Set("vpc_id", *vpnAttachment.VpcId)
	}

	vs := make([]map[string]interface{}, len(vpnGateway.VpcAttachments))

	for k, v := range vpnGateway.VpcAttachments {
		vp := make(map[string]interface{})

		vp["state"] = *v.State
		vp["vpc_id"] = *v.VpcId

		vs[k] = vp
	}

	d.Set("vpn_gateway_id", vpnGateway.VpnGatewayId)
	d.Set("attachments", vs)
	d.Set("state", vpnGateway.State)
	d.Set("tag_set", tagsToMap(vpnGateway.Tags))
	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceOutscaleVpnGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.VM.DeleteVpnGateway(&fcu.DeleteVpnGatewayInput{
			VpnGatewayId: aws.String(d.Id()),
		})
		if err == nil {
			return nil
		}

		ec2err, ok := err.(awserr.Error)
		if !ok {
			return resource.RetryableError(err)
		}

		switch ec2err.Code() {
		case "InvalidVpnGatewayID.NotFound":
			return nil
		case "IncorrectState":
			return resource.RetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
}

func resourceOutscaleVpnGatewayAttach(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	if d.Get("vpc_id").(string) == "" {
		fmt.Printf(
			"[DEBUG] Not attaching VPN Gateway '%s' as no VPC ID is set",
			d.Id())
		return nil
	}

	fmt.Printf(
		"[INFO] Attaching VPN Gateway '%s' to VPC '%s'",
		d.Id(),
		d.Get("vpc_id").(string))

	req := &fcu.AttachVpnGatewayInput{
		VpnGatewayId: aws.String(d.Id()),
		VpcId:        aws.String(d.Get("vpc_id").(string)),
	}

	err := resource.Retry(30*time.Second, func() *resource.RetryError {
		_, err := conn.VM.AttachVpnGateway(req)
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
		return err
	}

	// Wait for it to be fully attached before continuing
	fmt.Printf("[DEBUG] Waiting for VPN gateway (%s) to attach", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"detached", "attaching"},
		Target:  []string{"attached"},
		Refresh: vpnGatewayAttachStateRefreshFunc(conn, d.Id(), "available"),
		Timeout: 5 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for VPN gateway (%s) to attach: %s",
			d.Id(), err)
	}

	return nil
}

func resourceOutscaleVpnGatewayDetach(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	// Get the old VPC ID to detach from
	vpcID, _ := d.GetChange("vpc_id")

	if vpcID.(string) == "" {
		fmt.Printf(
			"[DEBUG] Not detaching VPN Gateway '%s' as no VPC ID is set",
			d.Id())
		return nil
	}

	fmt.Printf(
		"[INFO] Detaching VPN Gateway '%s' from VPC '%s'",
		d.Id(),
		vpcID.(string))

	wait := true

	var err error
	err = resource.Retry(30*time.Second, func() *resource.RetryError {
		_, err = conn.VM.DetachVpnGateway(&fcu.DetachVpnGatewayInput{
			VpnGatewayId: aws.String(d.Id()),
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
	fmt.Printf("[DEBUG] Waiting for VPN gateway (%s) to detach", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"attached", "detaching", "available"},
		Target:  []string{"detached"},
		Refresh: vpnGatewayAttachStateRefreshFunc(conn, d.Id(), "detached"),
		Timeout: 5 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for vpn gateway (%s) to detach: %s",
			d.Id(), err)
	}

	return nil
}

// vpnGatewayAttachStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// the state of a VPN gateway's attachment
func vpnGatewayAttachStateRefreshFunc(conn *fcu.Client, id string, expected string) resource.StateRefreshFunc {
	var start time.Time
	return func() (interface{}, string, error) {
		if start.IsZero() {
			start = time.Now()
		}

		var resp *fcu.DescribeVpnGatewaysOutput
		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.VM.DescribeVpnGateways(&fcu.DescribeVpnGatewaysInput{
				VpnGatewayIds: []*string{aws.String(id)},
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
			if strings.Contains(fmt.Sprint(err), "InvalidVpnGatewayID.NotFound") {
				resp = nil
			} else {
				fmt.Printf("[ERROR] Error on VpnGatewayStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			return nil, "", nil
		}

		vpnGateway := resp.VpnGateways[0]
		if len(vpnGateway.VpcAttachments) == 0 {
			return vpnGateway, "detached", nil
		}

		vpnAttachment := vpnGatewayGetAttachment(vpnGateway)
		return vpnGateway, *vpnAttachment.State, nil
	}
}

func vpnGatewayGetAttachment(vgw *fcu.VpnGateway) *fcu.VpcAttachment {
	for _, v := range vgw.VpcAttachments {
		if *v.State == "attached" {
			return v
		}
	}
	return &fcu.VpcAttachment{State: aws.String("detached")}
}
