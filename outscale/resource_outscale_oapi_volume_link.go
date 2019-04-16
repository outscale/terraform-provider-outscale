package outscale

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/outscale/osc-go/oapi"
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
			Optional: true,
			ForceNew: true,
			Computed: true,
		},
		"force_unlink": {
			Type:     schema.TypeBool,
			Optional: true,
			ForceNew: true,
			Computed: true,
		},
		"volume_id": {
			Type:     schema.TypeString,
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
	conn := meta.(*OutscaleClient).OAPI
	name := d.Get("device_name").(string)
	iID := d.Get("vm_id").(string)
	vID := d.Get("volume_id").(string)

	// Find out if the volume is already attached to the instance, in which case
	// we have nothing to do
	request := oapi.ReadVolumesRequest{
		Filters: oapi.FiltersVolume{
			VolumeIds: []string{vID},
		},
	}
	var err error
	var vols *oapi.POST_ReadVolumesResponses
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		vols, err = conn.POST_ReadVolumes(request)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if (err != nil) || isElegibleToLink(vols.OK.Volumes, iID) {
		// This handles the situation where the instance is created by
		// a spot request and whilst the request has been fulfilled the
		// instance is not running yet
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending"},
			Target:     []string{"running"},
			Refresh:    VMStateRefreshFunc(conn, iID, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for volume link (%s) to become ready: %s",
				iID, err)
		}

		// not attached
		opts := oapi.LinkVolumeRequest{
			DeviceName: name,
			VmId:       iID,
			VolumeId:   vID,
		}

		log.Printf("[DEBUG] Attaching Volume (%s) to Instance (%s)", vID, iID)

		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			_, err = conn.POST_LinkVolume(opts)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				return fmt.Errorf("[WARN] Error attaching volume (%s) to instance (%s), message: \"%s\", code: \"%s\"",
					vID, iID, awsErr.Message(), awsErr.Code())
			}
			return err
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"attaching"},
		Target:     []string{"attached"},
		Refresh:    volumeOAPIAttachmentStateRefreshFunc(conn, vID, iID),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Volume (%s) to attach to Instance: %s, error: %s",
			vID, iID, err)
	}

	d.SetId(volumeOAPIAttachmentID(name, vID, iID))
	return resourceOAPIVolumeLinkRead(d, meta)
}

func isElegibleToLink(volumes []oapi.Volume, instanceID string) bool {
	elegible := true

	if len(volumes) > 0 {
		for _, link := range volumes[0].LinkedVolumes {
			if instanceID == link.VmId {
				elegible = false
				break
			}
		}
	}

	return elegible
}

func isElegibleToUnLink(volumes []oapi.Volume, instanceID string) bool {
	return !isElegibleToLink(volumes, instanceID)
}

func volumeOAPIAttachmentStateRefreshFunc(conn *oapi.Client, volumeID, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		request := oapi.ReadVolumesRequest{
			Filters: oapi.FiltersVolume{
				VolumeIds: []string{volumeID},
			},
		}

		var err error
		var resp *oapi.POST_ReadVolumesResponses

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.POST_ReadVolumes(request)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				return nil, "failed", fmt.Errorf("code: %s, message: %s", awsErr.Code(), awsErr.Message())
			}
			return nil, "failed", err
		}

		if len(resp.OK.Volumes) > 0 {
			v := resp.OK.Volumes[0]
			for _, a := range v.LinkedVolumes {
				if a.VmId == instanceID {
					return a, a.State, nil
				}
			}
		}
		// assume detached if volume count is 0
		return 42, "detached", nil
	}
}

func resourceOAPIVolumeLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	request := oapi.ReadVolumesRequest{
		Filters: oapi.FiltersVolume{
			VolumeIds: []string{d.Get("volume_id").(string)},
		},
		//Name:   aws.String("attachment.instance-id"),

	}

	var err error
	var vols *oapi.POST_ReadVolumesResponses

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		vols, err = conn.POST_ReadVolumes(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidVolume.NotFound" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Outscale volume %s for instance: %s: %#v", d.Get("volume_id").(string), d.Get("vm_id").(string), err)
	}

	d.Set("request_id", vols.OK.ResponseContext.RequestId)

	if len(vols.OK.Volumes) == 0 || vols.OK.Volumes[0].State == "available" || isElegibleToLink(vols.OK.Volumes, d.Get("vm_id").(string)) {
		log.Printf("[DEBUG] Volume Attachment (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return nil

}

func resourceOAPIVolumeLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OAPI

	if _, ok := d.GetOk("skip_destroy"); ok {
		log.Printf("[INFO] Found skip_destroy to be true, removing attachment %q from state", d.Id())
		d.SetId("")
		return nil
	}

	vID := d.Get("volume_id").(string)
	iID := d.Get("vm_id").(string)

	opts := oapi.UnlinkVolumeRequest{
		//VmId:       iID,
		//ForceUnlink: d.Get("force_unlink").(bool),
		//DeviceName: d.Get("device_name").(string), //Removed due oAPI Bug.
		VolumeId: vID,
	}

	force, forceOk := d.GetOk("force_unlink")
	if forceOk {
		opts.ForceUnlink = force.(bool)

	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		_, err = conn.POST_UnlinkVolume(opts)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
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
		Delay:      10 * time.Second,
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

func volumeOAPIAttachmentID(name, volumeID, instanceID string) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s-", name))
	buf.WriteString(fmt.Sprintf("%s-", instanceID))
	buf.WriteString(fmt.Sprintf("%s-", volumeID))

	return fmt.Sprintf("vai-%d", hashcode.String(buf.String()))
}

func VMStateRefreshFunc(conn *oapi.Client, instanceID, failState string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp *oapi.POST_ReadVmsResponses
		var err error

		err = resource.Retry(30*time.Second, func() *resource.RetryError {
			resp, err = conn.POST_ReadVms(oapi.ReadVmsRequest{
				Filters: oapi.FiltersVm{VmIds: []string{instanceID}},
			})
			return resource.RetryableError(err)
		})

		if err != nil {
			return nil, "", err
		}

		if resp == nil || len(resp.OK.Vms) == 0 {
			return nil, "", nil
		}

		i := resp.OK.Vms[0]
		state := i.State

		if state == failState {
			return i, state, fmt.Errorf("Failed to reach target state. Reason: %v",
				i.State)

		}

		return i, state, nil
	}
}
