package outscale

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
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
	}
}

func resourceOAPIVolumeLinkCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU
	name := d.Get("device_name").(string)
	iID := d.Get("vm_id").(string)
	vID := d.Get("volume_id").(string)

	// Find out if the volume is already attached to the instance, in which case
	// we have nothing to do
	request := &fcu.DescribeVolumesInput{
		VolumeIds: []*string{aws.String(vID)},
		Filters: []*fcu.Filter{
			&fcu.Filter{
				Name:   aws.String("attachment.instance-id"),
				Values: []*string{aws.String(iID)},
			},
			&fcu.Filter{
				Name:   aws.String("attachment.device_name"),
				Values: []*string{aws.String(name)},
			},
		},
	}

	var err error
	var vols *fcu.DescribeVolumesOutput

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		vols, err = conn.VM.DescribeVolumes(request)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if (err != nil) || (len(vols.Volumes) == 0) {
		// This handles the situation where the instance is created by
		// a spot request and whilst the request has been fulfilled the
		// instance is not running yet
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"pending"},
			Target:     []string{"running"},
			Refresh:    InstanceStateRefreshFunc(conn, iID, ""),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf(
				"Error waiting for instance (%s) to become ready: %s",
				iID, err)
		}

		// not attached
		opts := &fcu.AttachVolumeInput{
			Device:     aws.String(name),
			InstanceId: aws.String(iID),
			VolumeId:   aws.String(vID),
		}

		log.Printf("[DEBUG] Attaching Volume (%s) to Instance (%s)", vID, iID)

		var err error

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			_, err = conn.VM.AttachVolume(opts)
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

func volumeOAPIAttachmentStateRefreshFunc(conn *fcu.Client, volumeID, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		request := &fcu.DescribeVolumesInput{
			VolumeIds: []*string{aws.String(volumeID)},
			Filters: []*fcu.Filter{
				&fcu.Filter{
					Name:   aws.String("attachment.instance-id"),
					Values: []*string{aws.String(instanceID)},
				},
			},
		}

		var err error
		var resp *fcu.DescribeVolumesOutput

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			resp, err = conn.VM.DescribeVolumes(request)
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

		if len(resp.Volumes) > 0 {
			v := resp.Volumes[0]
			for _, a := range v.Attachments {
				if a.InstanceId != nil && *a.InstanceId == instanceID {
					return a, *a.State, nil
				}
			}
		}
		// assume detached if volume count is 0
		return 42, "detached", nil
	}
}

func resourceOAPIVolumeLinkRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	request := &fcu.DescribeVolumesInput{
		VolumeIds: []*string{aws.String(d.Get("volume_id").(string))},
		Filters: []*fcu.Filter{
			&fcu.Filter{
				Name:   aws.String("attachment.instance-id"),
				Values: []*string{aws.String(d.Get("vm_id").(string))},
			},
		},
	}

	var err error
	var vols *fcu.DescribeVolumesOutput

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		vols, err = conn.VM.DescribeVolumes(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidVolume.NotFound" {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Outscale volume %s for instance: %s: %#v", d.Get("volume_id").(string), d.Get("vm_id").(string), err)
	}

	if len(vols.Volumes) == 0 || *vols.Volumes[0].State == "available" {
		log.Printf("[DEBUG] Volume Attachment (%s) not found, removing from state", d.Id())
		d.SetId("")
	}

	return nil

}

func resourceOAPIVolumeLinkDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	if _, ok := d.GetOk("skip_destroy"); ok {
		log.Printf("[INFO] Found skip_destroy to be true, removing attachment %q from state", d.Id())
		d.SetId("")
		return nil
	}

	vID := d.Get("volume_id").(string)
	iID := d.Get("vm_id").(string)

	opts := &fcu.DetachVolumeInput{
		Device:     aws.String(d.Get("device_name").(string)),
		InstanceId: aws.String(iID),
		VolumeId:   aws.String(vID),
		//Force:      aws.Bool(d.Get("force_detach").(bool)),
	}

	force, forceOk := d.GetOk("force_detach")
	if forceOk {
		opts.Force = aws.Bool(force.(bool))

	}

	var err error
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {

		_, err = conn.VM.DetachVolume(opts)

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
