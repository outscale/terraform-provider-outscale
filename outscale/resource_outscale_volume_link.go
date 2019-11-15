package outscale

import (
	"bytes"
	"context"
	"fmt"
	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
	connOsc := meta.(*OutscaleClient).OSCAPI
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
	var vols oscgo.ReadVolumesResponse
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		vols, _, err = conn.VolumeApi.ReadVolumes(context.Background(), &oscgo.ReadVolumesOpts{ReadVolumesRequest: optional.NewInterface(request)})

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if (err != nil) || isElegibleToLink(vols.GetVolumes(), iID) {
		// This handles the situation where the instance is created by
		// a spot request and whilst the request has been fulfilled the
		// instance is not running yet

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending"},
			Target:     []string{"running"},
			Refresh:    volumeOAPIAttachmentStateRefreshFunc(conn, iID, ""),
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
		opts := oscgo.LinkVolumeRequest{
			DeviceName: name,
			VmId:       iID,
			VolumeId:   vID,
		}

		log.Printf("[DEBUG] Attaching Volume (%s) to Instance (%s)", vID, iID)

		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			_, _, err = conn.VolumeApi.LinkVolume(context.Background(), &oscgo.LinkVolumeOpts{LinkVolumeRequest: optional.NewInterface(opts)})
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

func isElegibleToUnLink(volumes []oscgo.Volume, instanceID string) bool {
	return !isElegibleToLink(volumes, instanceID)
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
			resp, _, err = conn.VolumeApi.ReadVolumes(context.Background(), &oscgo.ReadVolumesOpts{ReadVolumesRequest: optional.NewInterface(request)})
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
			VolumeIds: &[]string{d.Get("volume_id").(string)},
		},
	}

	var err error
	var vols oscgo.ReadVolumesResponse

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		vols, _, err = conn.VolumeApi.ReadVolumes(context.Background(), &oscgo.ReadVolumesOpts{ReadVolumesRequest: optional.NewInterface(request)})
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

	d.Set("request_id", vols.ResponseContext.GetRequestId())

	if len(vols.GetVolumes()) == 0 || vols.GetVolumes()[0].GetState() == "available" || isElegibleToLink(vols.GetVolumes(), d.Get("vm_id").(string)) {
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

	vID := d.Get("volume_id").(string)
	iID := d.Get("vm_id").(string)

	opts := oscgo.UnlinkVolumeRequest{
		//VmId:       iID,
		//ForceUnlink: d.Get("force_unlink").(bool),
		//DeviceName: d.Get("device_name").(string), //Removed due oAPI Bug.
		VolumeId: vID,
	}

	force, forceOk := d.GetOk("force_unlink")
	if forceOk {
		opts.SetForceUnlink(force.(bool))

	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, _, err = conn.VolumeApi.UnlinkVolume(context.Background(), &oscgo.UnlinkVolumeOpts{UnlinkVolumeRequest: optional.NewInterface(opts)})

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
