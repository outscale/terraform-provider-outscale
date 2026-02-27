package oapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/outscale/goutils/sdk/ptr"
	"github.com/outscale/osc-sdk-go/v3/pkg/options"
	"github.com/outscale/osc-sdk-go/v3/pkg/osc"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleNetworkInterfaceAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceOutscaleNetworkInterfaceAttachmentCreate,
		ReadContext:   ResourceOutscaleNetworkInterfaceAttachmentRead,
		DeleteContext: ResourceOutscaleNetworkInterfaceAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: ResourceOutscaleNetworkInterfaceAttachmentImportStateContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(CreateDefaultTimeout),
			Read:   schema.DefaultTimeout(ReadDefaultTimeout),
			Delete: schema.DefaultTimeout(DeleteDefaultTimeout),
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

func ResourceOutscaleNetworkInterfaceAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutCreate)

	di := d.Get("device_number").(int)
	vmID := d.Get("vm_id").(string)
	nicID := d.Get("nic_id").(string)

	opts := osc.LinkNicRequest{
		DeviceNumber: di,
		VmId:         vmID,
		NicId:        nicID,
	}

	log.Printf("[DEBUG] Attaching network interface (%s) to instance (%s)", nicID, vmID)

	resp, err := client.LinkNic(ctx, opts, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.Errorf("error creating outscale linknic: %s", err)
	}

	d.SetId(ptr.From(resp.LinkNicId))
	return ResourceOutscaleNetworkInterfaceAttachmentRead(ctx, d, meta)
}

func ResourceOutscaleNetworkInterfaceAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC
	timeout := d.Timeout(schema.TimeoutRead)

	nicID := d.Get("nic_id").(string)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"attaching", "detaching"},
		Target:  []string{"attached", "detached", "failed"},
		Timeout: timeout,
		Refresh: nicLinkRefreshFunc(ctx, client, nicID, timeout),
	}

	resp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for nic to attach to Instance: %s, error: %s", nicID, err)
	}

	r := resp.(*osc.ReadNicsResponse)
	if r == nil || r.Nics == nil || utils.IsResponseEmpty(len(*r.Nics), "NicLink", d.Id()) {
		d.SetId("")
		return nil
	}
	linkNic := (*r.Nics)[0].LinkNic
	if err := d.Set("device_number", linkNic.DeviceNumber); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("vm_id", linkNic.VmId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("delete_on_vm_deletion", linkNic.DeleteOnVmDeletion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("link_nic_id", linkNic.LinkNicId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func ResourceOutscaleNetworkInterfaceAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.OutscaleClient).OSC

	timeout := d.Timeout(schema.TimeoutDelete)

	interfaceID := d.Id()

	req := osc.UnlinkNicRequest{
		LinkNicId: interfaceID,
	}

	_, err := client.UnlinkNic(ctx, req, options.WithRetryTimeout(timeout))
	if err != nil {
		return diag.FromErr(err)
	}

	nicID := d.Get("nic_id").(string)

	// log.Printf("[DEBUG] Waiting for ENI (%s) to become dettached", interfaceID)
	stateConf := &retry.StateChangeConf{
		Pending: []string{"detaching"},
		Target:  []string{"detached", "failed"},
		Timeout: timeout,
		Refresh: nicLinkRefreshFunc(ctx, client, nicID, timeout),
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"error waiting for volume to dettached from instance: %s, error: %s", nicID, err)
	}

	return nil
}

func ResourceOutscaleNetworkInterfaceAttachmentImportStateContext(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if d.Id() == "" {
		return nil, errors.New("import error: to import a Nic Link, use the format {nic_id} it must not be empty")
	}

	timeout := d.Timeout(schema.TimeoutRead)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"attaching", "detaching"},
		Target:  []string{"attached", "detached", "failed"},
		Timeout: timeout,
		Refresh: nicLinkRefreshFunc(ctx, meta.(*client.OutscaleClient).OSC, d.Id(), timeout),
	}

	resp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for nic to attach to instance: %s, error: %s", d.Id(), err)
	}
	r := resp.(*osc.ReadNicsResponse)
	if r == nil || r.Nics == nil || len(*r.Nics) == 0 {
		return nil, fmt.Errorf("nic not found: %v", d.Id())
	}

	linkNic := (*r.Nics)[0].LinkNic
	if err := d.Set("device_number", linkNic.DeviceNumber); err != nil {
		return nil, err
	}
	if err := d.Set("vm_id", linkNic.VmId); err != nil {
		return nil, err
	}
	if err := d.Set("nic_id", (*r.Nics)[0].NicId); err != nil {
		return nil, err
	}

	d.SetId(linkNic.LinkNicId)

	return []*schema.ResourceData{d}, nil
}

func nicLinkRefreshFunc(ctx context.Context, client *osc.Client, nicID string, timeout time.Duration) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := osc.ReadNicsRequest{
			Filters: &osc.FiltersNic{
				NicIds: &[]string{nicID},
			},
		}

		resp, err := client.ReadNics(ctx, req, options.WithRetryTimeout(timeout))
		if err != nil {
			return nil, "failed", err
		}
		if resp.Nics == nil || len(*resp.Nics) < 1 {
			return nil, "failed", fmt.Errorf("error to find the nic(%s): %#v", nicID, resp.Nics)
		}

		linkNic := ptr.From((*resp.Nics)[0].LinkNic)
		if reflect.DeepEqual(linkNic, osc.LinkNic{}) {
			return resp, "detached", nil
		}

		return resp, string(linkNic.State), nil
	}
}
