package outscale

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOAPINetworkInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceOutscaleOAPINetworkInterfaceAttachmentCreate,
		Read:   resourceOutscaleOAPINetworkInterfaceAttachmentRead,
		Delete: resourceOutscaleOAPINetworkInterfaceAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"device_number": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"vm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"nic_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delete_on_vm_deletion": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vm_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"link_nic_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPINetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	di := d.Get("device_number").(int)
	iID := d.Get("vm_id").(string)
	nicID := d.Get("nic_id").(string)

	opts := oscgo.LinkNicRequest{
		DeviceNumber: int32(di),
		VmId:         iID,
		NicId:        nicID,
	}

	log.Printf("[DEBUG] Attaching network interface (%s) to instance (%s)", nicID, iID)

	var resp oscgo.LinkNicResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.NicApi.LinkNic(context.Background(), &oscgo.LinkNicOpts{LinkNicRequest: optional.NewInterface(opts)})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		errString := err.Error()
		return fmt.Errorf("Error creating Outscale LinkNic: %s", errString)

	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"false"},
		Target:     []string{"true"},
		Refresh:    nicLinkRefreshFunc(conn, nicID),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Volume (%s) to attach to Instance: %s, error: %s", nicID, iID, err)
	}

	d.SetId(resp.GetLinkNicId())
	return resourceOutscaleOAPINetworkInterfaceAttachmentRead(d, meta)
}

func resourceOutscaleOAPINetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	interfaceID := d.Get("nic_id").(string)

	req := oscgo.ReadNicsRequest{
		Filters: &oscgo.FiltersNic{NicIds: &[]string{interfaceID}},
	}

	var resp oscgo.ReadNicsResponse
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.NicApi.ReadNics(context.Background(), &oscgo.ReadNicsOpts{ReadNicsRequest: optional.NewInterface(req)})
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
		errString := err.Error()
		return fmt.Errorf("Could not find network interface: %s", errString)

	}
	if len(resp.GetNics()) != 1 {
		return fmt.Errorf("Unable to find ENI (%s): %#v", interfaceID, resp.GetNics())
	}

	eni := resp.GetNics()[0]

	if reflect.DeepEqual(eni.GetLinkNic(), oscgo.LinkNic{}) {
		// Interface is no longer attached, remove from state
		d.SetId("")
		return nil
	}

	link := eni.GetLinkNic()

	if link.GetVmAccountId() != "" {
		d.Set("vm_account_id", link.GetVmAccountId())
	}
	if link.GetState() != "" {
		d.Set("state", link.GetState())
	}

	d.Set("device_number", fmt.Sprintf("%d", link.GetDeviceNumber()))
	d.Set("vm_id", link.GetVmId())
	d.Set("delete_on_vm_deletion", link.GetDeleteOnVmDeletion())
	d.Set("link_nic_id", link.GetLinkNicId())
	d.Set("request_id", resp.ResponseContext.GetRequestId())

	return nil
}

func resourceOutscaleOAPINetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	interfaceID := d.Get("nic_id").(string)

	dr := oscgo.UnlinkNicRequest{
		LinkNicId: d.Id(),
	}

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.NicApi.UnlinkNic(context.Background(), &oscgo.UnlinkNicOpts{UnlinkNicRequest: optional.NewInterface(dr)})
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
		Refresh: nicLinkRefreshFunc(conn, interfaceID),
		Timeout: 10 * time.Minute,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for ENI (%s) to become dettached: %s", interfaceID, err)
	}

	return nil
}
