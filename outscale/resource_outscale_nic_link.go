package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleNetworkInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleNetworkInterfaceAttachmentCreate,
		Read:   resourceOutscaleNetworkInterfaceAttachmentRead,
		Delete: resourceOutscaleNetworkInterfaceAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"device_index": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network_interface_id": {
				Type:     schema.TypeString,
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

func resourceOutscaleNetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	di := d.Get("device_index").(int)
	iID := d.Get("instance_id").(string)
	nicID := d.Get("network_interface_id").(string)

	opts := &fcu.AttachNetworkInterfaceInput{
		DeviceIndex:        aws.Int64(int64(di)),
		InstanceId:         aws.String(iID),
		NetworkInterfaceId: aws.String(nicID),
	}

	log.Printf("[DEBUG] Attaching network interface (%s) to instance (%s)", nicID, iID)

	var resp *fcu.AttachNetworkInterfaceOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		resp, err = conn.VM.AttachNetworkInterface(opts)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"false"},
		Target:     []string{"true"},
		Refresh:    networkInterfaceAttachmentRefreshFunc(conn, nicID),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Volume (%s) to attach to Instance: %s, error: %s", nicID, iID, err)
	}

	d.SetId(*resp.AttachmentId)
	return resourceOutscaleNetworkInterfaceAttachmentRead(d, meta)
}

func resourceOutscaleNetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
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

	if eni.Attachment == nil {
		// Interface is no longer attached, remove from state
		d.SetId("")
		return nil
	}

	d.Set("device_index", eni.Attachment.DeviceIndex)
	d.Set("instance_id", eni.Attachment.InstanceId)
	d.Set("network_interface_id", eni.NetworkInterfaceId)
	d.Set("request_id", resp.RequestId)

	return nil
}

func resourceOutscaleNetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	interfaceID := d.Get("network_interface_id").(string)

	dr := &fcu.DetachNetworkInterfaceInput{
		AttachmentId: aws.String(d.Id()),
		Force:        aws.Bool(true),
	}

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		_, err = conn.VM.DetachNetworkInterface(dr)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidAttachmentID.NotFound") {
			return fmt.Errorf("Error detaching ENI: %s", err)
		}
	}

	log.Printf("[DEBUG] Waiting for ENI (%s) to become dettached", interfaceID)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"true"},
		Target:  []string{"false"},
		Refresh: networkInterfaceAttachmentRefreshFunc(conn, interfaceID),
		Timeout: 10 * time.Minute,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for ENI (%s) to become dettached: %s", interfaceID, err)
	}

	return nil
}
