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

func resourceOutscaleOAPIVpnGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPIVpnGatewayCreate,
		Read:   resourceOutscaleOAPIVpnGatewayRead,
		Delete: resourceOutscaleOAPIVpnGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"lin_to_vpn_gateway_link": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"lin_id": &schema.Schema{
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
			"lin_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag":  tagsSchemaComputed(),
			"tags": tagsSchema(),
		},
	}
}

func resourceOutscaleOAPIVpnGatewayCreate(d *schema.ResourceData, meta interface{}) error {
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
		d.SetPartial("tag")
	}

	return resourceOutscaleOAPIVpnGatewayRead(d, meta)
}

func resourceOutscaleOAPIVpnGatewayRead(d *schema.ResourceData, meta interface{}) error {
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
		}
		fmt.Printf("\n\n[ERROR] Error finding VpnGateway: %s", err)
		return err
	}

	vpnGateway := resp.VpnGateways[0]
	if vpnGateway == nil || *vpnGateway.State == "deleted" {
		d.SetId("")
		return nil
	}

	vpnAttachment := oapiVpnGatewayGetAttachment(vpnGateway)
	if len(vpnGateway.VpcAttachments) == 0 || *vpnAttachment.State == "detached" {
		d.Set("lin_id", "")
	} else {
		d.Set("lin_id", *vpnAttachment.VpcId)
	}

	vs := make([]map[string]interface{}, len(vpnGateway.VpcAttachments))

	for k, v := range vpnGateway.VpcAttachments {
		vp := make(map[string]interface{})

		vp["state"] = *v.State
		vp["lin_id"] = *v.VpcId

		vs[k] = vp
	}

	d.Set("vpn_gateway_id", vpnGateway.VpnGatewayId)
	d.Set("lin_to_vpn_gateway_link", vs)
	d.Set("state", vpnGateway.State)
	d.Set("tag", tagsToMap(vpnGateway.Tags))
	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceOutscaleOAPIVpnGatewayDelete(d *schema.ResourceData, meta interface{}) error {
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

func oapiVpnGatewayGetAttachment(vgw *fcu.VpnGateway) *fcu.VpcAttachment {
	for _, v := range vgw.VpcAttachments {
		if *v.State == "attached" {
			return v
		}
	}
	return &fcu.VpcAttachment{State: aws.String("detached")}
}
