package oapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleNetworkInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		Create: ResourceOutscaleNetworkInterfaceAttachmentCreate,
		Read:   ResourceOutscaleNetworkInterfaceAttachmentRead,
		Delete: ResourceOutscaleNetworkInterfaceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: ResourceOutscaleNetworkInterfaceAttachmentImportState,
		},
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
				Type:     schema.TypeBool,
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

func ResourceOutscaleNetworkInterfaceAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	di := d.Get("device_number").(int)
	vmID := d.Get("vm_id").(string)
	nicID := d.Get("nic_id").(string)

	opts := oscgo.LinkNicRequest{
		DeviceNumber: int32(di),
		VmId:         vmID,
		NicId:        nicID,
	}

	log.Printf("[DEBUG] Attaching network interface (%s) to instance (%s)", nicID, vmID)

	var resp oscgo.LinkNicResponse
	var err error
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.NicApi.LinkNic(context.Background()).LinkNicRequest(opts).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating Outscale LinkNic: %s", err)
	}

	d.SetId(resp.GetLinkNicId())
	return ResourceOutscaleNetworkInterfaceAttachmentRead(d, meta)
}

func ResourceOutscaleNetworkInterfaceAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	nicID := d.Get("nic_id").(string)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"attaching", "detaching"},
		Target:     []string{"attached", "detached", "failed"},
		Refresh:    nicLinkRefreshFunc(conn, nicID),
		Timeout:    5 * time.Minute,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	resp, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for NIC to attach to Instance: %s, error: %s", nicID, err)
	}

	r := resp.(oscgo.ReadNicsResponse)
	if utils.IsResponseEmpty(len(r.GetNics()), "NicLink", d.Id()) {
		d.SetId("")
		return nil
	}
	linkNic := r.GetNics()[0].GetLinkNic()

	if err := d.Set("device_number", linkNic.GetDeviceNumber()); err != nil {
		return err
	}
	if err := d.Set("vm_id", linkNic.GetVmId()); err != nil {
		return err
	}
	if err := d.Set("delete_on_vm_deletion", linkNic.GetDeleteOnVmDeletion()); err != nil {
		return err
	}
	if err := d.Set("link_nic_id", linkNic.GetLinkNicId()); err != nil {
		return err
	}

	return nil
}

func ResourceOutscaleNetworkInterfaceAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	interfaceID := d.Id()

	req := oscgo.UnlinkNicRequest{
		LinkNicId: interfaceID,
	}

	var err error
	var statusCode int
	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		_, httpResp, err := conn.NicApi.UnlinkNic(context.Background()).UnlinkNicRequest(req).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		statusCode = httpResp.StatusCode
		return nil
	})

	if err != nil {
		if statusCode == http.StatusNotFound {
			return fmt.Errorf("Error detaching ENI: %s", err)
		}
	}

	nicID := d.Get("nic_id").(string)

	// log.Printf("[DEBUG] Waiting for ENI (%s) to become dettached", interfaceID)
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"detaching"},
		Target:     []string{"detached", "failed"},
		Refresh:    nicLinkRefreshFunc(conn, nicID),
		Timeout:    5 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Volume to dettache to Instance: %s, error: %s", nicID, err)
	}

	return nil
}

func ResourceOutscaleNetworkInterfaceAttachmentImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if d.Id() == "" {
		return nil, errors.New("import error: to import a Nic Link, use the format {nic_id} it must not be empty")
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"attaching", "detaching"},
		Target:     []string{"attached", "detached", "failed"},
		Refresh:    nicLinkRefreshFunc(meta.(*client.OutscaleClient).OSCAPI, d.Id()),
		Timeout:    5 * time.Minute,
		MinTimeout: 3 * time.Second,
	}

	resp, err := stateConf.WaitForState()
	if err != nil {
		return nil, fmt.Errorf(
			"Error waiting for NIC to attach to Instance: %s, error: %s", d.Id(), err)
	}
	r := resp.(oscgo.ReadNicsResponse)
	linkNic := r.GetNics()[0].GetLinkNic()

	if err := d.Set("device_number", linkNic.GetDeviceNumber()); err != nil {
		return nil, err
	}
	if err := d.Set("vm_id", linkNic.GetVmId()); err != nil {
		return nil, err
	}
	if err := d.Set("nic_id", r.GetNics()[0].GetNicId()); err != nil {
		return nil, err
	}

	d.SetId(linkNic.GetLinkNicId())

	return []*schema.ResourceData{d}, nil
}

func nicLinkRefreshFunc(conn *oscgo.APIClient, nicID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := oscgo.ReadNicsRequest{
			Filters: &oscgo.FiltersNic{
				NicIds: &[]string{nicID},
			},
		}

		var resp oscgo.ReadNicsResponse
		var err error
		err = retry.Retry(5*time.Minute, func() *retry.RetryError {
			rp, httpResp, err := conn.NicApi.ReadNics(context.Background()).ReadNicsRequest(req).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil {
			return nil, "failed", err
		}
		if len(resp.GetNics()) < 1 {
			return nil, "failed", fmt.Errorf("error to find the Outscale Nic(%s): %#v", nicID, resp.GetNics())
		}

		linkNic := resp.GetNics()[0].GetLinkNic()
		if reflect.DeepEqual(linkNic, oscgo.LinkNic{}) {
			return resp, "detached", nil
		}

		return resp, linkNic.GetState(), nil
	}
}
