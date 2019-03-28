package outscale

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/oapi"
	"github.com/terraform-providers/terraform-provider-outscale/utils"
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOutscaleOAPINetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	di := d.Get("device_number").(int)
	iID := d.Get("vm_id").(string)
	nicID := d.Get("nic_id").(string)

	opts := &oapi.LinkNicRequest{
		DeviceNumber: int64(di),
		VmId:         iID,
		NicId:        nicID,
	}

	log.Printf("[DEBUG] Attaching network interface (%s) to instance (%s)", nicID, iID)

	var resp *oapi.POST_LinkNicResponses
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		resp, err = conn.POST_LinkNic(*opts)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil || resp.OK == nil {
		if err != nil {
			errString = err.Error()
		} else if resp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(resp.Code401))
		} else if resp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(resp.Code400))
		} else if resp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(resp.Code500))
		}
		return fmt.Errorf("Error creating Outscale LinkNic: %s", errString)

	}

	result := resp.OK

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

	d.SetId(result.LinkNicId)
	return resourceOutscaleOAPINetworkInterfaceAttachmentRead(d, meta)
}

func resourceOutscaleOAPINetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	interfaceID := d.Get("nic_id").(string)

	req := &oapi.ReadNicsRequest{
		Filters: oapi.FiltersNic{NicIds: []string{interfaceID}},
	}

	var describeResp *oapi.POST_ReadNicsResponses
	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		describeResp, err = conn.POST_ReadNics(*req)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	var errString string

	if err != nil || describeResp.OK == nil {
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "InvalidNetworkInterfaceID.NotFound") {
				// The ENI is gone now, so just remove the attachment from the state
				d.SetId("")
				return nil
			}
			errString = err.Error()
		} else if describeResp.Code401 != nil {
			errString = fmt.Sprintf("ErrorCode: 401, %s", utils.ToJSONString(describeResp.Code401))
		} else if describeResp.Code400 != nil {
			errString = fmt.Sprintf("ErrorCode: 400, %s", utils.ToJSONString(describeResp.Code400))
		} else if describeResp.Code500 != nil {
			errString = fmt.Sprintf("ErrorCode: 500, %s", utils.ToJSONString(describeResp.Code500))
		}
		return fmt.Errorf("Could not find network interface: %s", errString)

	}

	result := describeResp.OK

	if len(result.Nics) != 1 {
		return fmt.Errorf("Unable to find ENI (%s): %#v", interfaceID, result.Nics)
	}

	eni := result.Nics[0]

	if reflect.DeepEqual(eni.LinkNic, oapi.LinkNic{}) {
		// Interface is no longer attached, remove from state
		d.SetId("")
		return nil
	}

	link := eni.LinkNic

	d.Set("device_number", link.DeviceNumber)
	d.Set("vm_id", link.VmId)
	d.Set("state", link.State)
	d.Set("delete_on_vm_deletion", link.DeleteOnVmDeletion)
	d.Set("vm_account_id", link.VmAccountId)
	d.Set("request_id", result.ResponseContext.RequestId)

	return nil
}

func resourceOutscaleOAPINetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	interfaceID := d.Get("nic_id").(string)

	dr := &oapi.UnlinkNicRequest{
		LinkNicId: d.Id(),
	}

	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		_, err = conn.POST_UnlinkNic(*dr)
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
