package outscale

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceOutscaleOAPIVolumeLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIVolumeLinkCreate,
		Read:   resourceOAPIVolumeLinkRead,
		Delete: resourceOAPIVolumeLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getOAPIVolumeLinkSchema(),
	}
}

func getOAPIVolumeLinkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Arguments
		"device_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"vm_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"volume_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"force_unlink": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
			Computed: true,
		},
		// Attributes
		"delete_on_vm_termination": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"request_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceOAPIVolumeLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI
	name := d.Get("device_name").(string)
	iID := d.Get("vm_id").(string)
	vID := d.Get("volume_id").(string)

	// Find out if the volume is already attached to the instance, in which case
	// we have nothing to do
	request := oscgo.ReadVolumesRequest{
		Filters: &oscgo.FiltersVolume{
			VolumeIds: &[]string{vID},
		},
	}
	var err error
	var resp oscgo.ReadVolumesResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.VolumeApi.ReadVolumes(context.Background()).ReadVolumesRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})

	if (err != nil) || isElegibleToLink(resp.GetVolumes(), iID) {
		// This handles the situation where the instance is created by
		// a spot request and whilst the request has been fulfilled the
		// instance is not running yet

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending"},
			Target:     []string{"running"},
			Refresh:    vmStateRefreshFunc(conn, iID, ""),
			Timeout:    10 * time.Minute,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"error waiting for the VM(%s) to become ready to be linked with the volume: %s",
				iID, err)
		}

		// not attached
		opts := oscgo.LinkVolumeRequest{
			DeviceName: name,
			VmId:       iID,
			VolumeId:   vID,
		}

		log.Printf("[DEBUG] Attaching Volume (%s) to Instance (%s)", vID, iID)

		err := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			_, httpResp, err := conn.VolumeApi.LinkVolume(context.Background()).LinkVolumeRequest(opts).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("[WARN] Error attaching volume (%s) to instance (%s), message:'%s'", vID, iID, err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"attaching"},
		Target:     []string{"attached"},
		Refresh:    volumeOAPIAttachmentStateRefreshFunc(conn, vID, iID),
		Timeout:    5 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Volume (%s) to attach to Instance: %s, error: %s",
			vID, iID, err)
	}

	d.SetId(vID)
	return resourceOAPIVolumeLinkRead(d, meta)
}

func isElegibleToLink(volumes []oscgo.Volume, instanceID string) bool {
	elegible := true

	if len(volumes) > 0 {
		for _, link := range volumes[0].GetLinkedVolumes() {
			if instanceID == link.GetVmId() {
				elegible = false
				break
			}
		}
	}

	return elegible
}

func volumeOAPIAttachmentStateRefreshFunc(conn *oscgo.APIClient, volumeID, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		request := oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{
				VolumeIds: &[]string{volumeID},
			},
		}

		var err error
		var resp oscgo.ReadVolumesResponse

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			rp, httpResp, err := conn.VolumeApi.ReadVolumes(context.Background()).ReadVolumesRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})

		if err != nil {
			return nil, "failed", err
		}

		if len(resp.GetVolumes()) > 0 {
			v := resp.GetVolumes()[0]
			for _, a := range v.GetLinkedVolumes() {
				if a.GetVmId() == instanceID {
					return a, a.GetState(), nil
				}
			}
		}
		// assume detached if volume count is 0
		return 42, "detached", nil
	}
}

func resourceOAPIVolumeLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	request := oscgo.ReadVolumesRequest{
		Filters: &oscgo.FiltersVolume{
			VolumeIds: &[]string{d.Id()},
		},
	}

	var err error
	var resp oscgo.ReadVolumesResponse
	var statusCode int
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		rp, httpResp, err := conn.VolumeApi.ReadVolumes(context.Background()).ReadVolumesRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		statusCode = httpResp.StatusCode
		return nil
	})

	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Outscale volume %s for instance: %s: %#v", d.Get("volume_id").(string), d.Get("vm_id").(string), err)
	}
	if utils.IsResponseEmpty(len(resp.GetVolumes()), "VolumeLink", d.Id()) {
		d.SetId("")
		return nil
	}
	var linkedVolume oscgo.LinkedVolume

	for _, vol := range resp.GetVolumes()[0].GetLinkedVolumes() {
		linkedVolume = vol
	}

	if err := d.Set("device_name", linkedVolume.GetDeviceName()); err != nil {
		return fmt.Errorf("error sertting %s in Volume Link(%s): %s", `device_name`, linkedVolume.GetVolumeId(), err)
	}
	if err := d.Set("vm_id", linkedVolume.GetVmId()); err != nil {
		return fmt.Errorf("error sertting %s in Volume Link(%s): %s", `vm_id`, linkedVolume.GetVolumeId(), err)
	}
	if err := d.Set("volume_id", linkedVolume.GetVolumeId()); err != nil {
		return fmt.Errorf("error sertting %s in Volume Link(%s): %s", `volume_id`, linkedVolume.GetVolumeId(), err)
	}
	if err := d.Set("delete_on_vm_termination", linkedVolume.GetDeleteOnVmDeletion()); err != nil {
		return fmt.Errorf("error sertting %s in Volume Link(%s): %s", `delete_on_vm_termination`, linkedVolume.GetVolumeId(), err)
	}
	if err := d.Set("state", linkedVolume.GetState()); err != nil {
		return fmt.Errorf("error sertting %s in Volume Link(%s): %s", `state`, linkedVolume.GetVolumeId(), err)
	}
	if len(resp.GetVolumes()) == 0 || resp.GetVolumes()[0].GetState() == "available" || isElegibleToLink(resp.GetVolumes(), d.Get("vm_id").(string)) {
		log.Printf("[DEBUG] Volume Attachment (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return nil
}

func resourceOAPIVolumeLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	if _, ok := d.GetOk("skip_destroy"); ok {
		log.Printf("[INFO] Found skip_destroy to be true, removing attachment %q from state", d.Id())
		d.SetId("")
		return nil
	}

	vID := d.Id()
	iID := d.Get("vm_id").(string)

	opts := oscgo.UnlinkVolumeRequest{
		VolumeId: vID,
	}

	force, forceOk := d.GetOk("force_unlink")
	if forceOk {
		opts.SetForceUnlink(force.(bool))

	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		_, httpResp, err := conn.VolumeApi.UnlinkVolume(context.Background()).UnlinkVolumeRequest(opts).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Failed to detach Volume (%s) from Instance (%s): %s",
			vID, iID, err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"detaching"},
		Target:     []string{"detached"},
		Refresh:    volumeOAPIAttachmentStateRefreshFunc(conn, vID, iID),
		Timeout:    5 * time.Minute,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	log.Printf("[DEBUG] Detaching Volume (%s) from Instance (%s)", vID, iID)
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Volume (%s) to detach from Instance: %s",
			vID, iID)
	}
	d.SetId("")
	return nil
}
