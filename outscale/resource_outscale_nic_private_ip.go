package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleNetworkInterfacePrivateIP() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleNetworkInterfacePrivateIPCreate,
		Read:   resourceOutscaleNetworkInterfacePrivateIPRead,
		Delete: resourceOutscaleNetworkInterfacePrivateIPDelete,

		Schema: map[string]*schema.Schema{
			"allow_reassignment": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"secondary_private_ip_address_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"network_interface_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"private_ip_address": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceOutscaleNetworkInterfacePrivateIPCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	input := &fcu.AssignPrivateIpAddressesInput{
		NetworkInterfaceId: aws.String(d.Get("network_interface_id").(string)),
	}

	if v, ok := d.GetOk("allow_reassignment"); ok {
		input.AllowReassignment = aws.Bool(v.(bool))
	}

	if v, ok := d.GetOk("secondary_private_ip_address_count"); ok {
		input.SecondaryPrivateIpAddressCount = aws.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("private_ip_address"); ok {
		input.PrivateIpAddresses = expandStringList(v.([]interface{}))
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err = conn.VM.AssignPrivateIpAddresses(input)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure to assign Private IPs: %s", err)
	}

	d.SetId(*input.NetworkInterfaceId)

	return resourceOutscaleNetworkInterfacePrivateIPRead(d, meta)
}

func resourceOutscaleNetworkInterfacePrivateIPRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	interfaceID := d.Get("network_interface_id").(string)

	req := &fcu.DescribeNetworkInterfacesInput{
		NetworkInterfaceIds: []*string{aws.String(interfaceID)},
	}

	var resp *fcu.DescribeNetworkInterfacesOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		resp, err = conn.VM.DescribeNetworkInterfaces(req)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidNetworkInterfaceID.NotFound") {
			// The ENI is gone now, so just remove the attachment from the state
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving ENI: %s", err)
	}
	if len(resp.NetworkInterfaces) != 1 {
		return fmt.Errorf("Unable to find ENI (%s): %#v", interfaceID, resp.NetworkInterfaces)
	}

	eni := resp.NetworkInterfaces[0]

	if eni.NetworkInterfaceId == nil {
		// Interface is no longer attached, remove from state
		d.SetId("")
		return nil
	}

	var ips []string
	for _, v := range eni.PrivateIpAddresses {
		ips = append(ips, *v.PrivateIpAddress)
	}

	_, ok := d.GetOk("allow_reassignment")

	d.Set("allow_reassignment", ok)
	d.Set("private_ip_address", ips)
	d.Set("secondary_private_ip_address_count", len(eni.PrivateIpAddresses))
	d.Set("network_interface_id", eni.NetworkInterfaceId)
	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceOutscaleNetworkInterfacePrivateIPDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	input := &fcu.UnassignPrivateIpAddressesInput{
		NetworkInterfaceId: aws.String(d.Id()),
	}

	if v, ok := d.GetOk("private_ip_address"); ok {
		input.PrivateIpAddresses = expandStringList(v.([]interface{}))
	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := conn.VM.UnassignPrivateIpAddresses(input)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failure to unassign Private IPs: %s", err)
	}

	return nil
}
