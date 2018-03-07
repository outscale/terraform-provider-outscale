package outscale

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-outscale/osc/fcu"
)

func resourceOutscaleVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceVolumeCreate,
		Read:   resourceVolumeRead,
		Delete: resourceVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: getVolumeSchema(),
	}
}

func resourceVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	request := &fcu.CreateVolumeInput{
		AvailabilityZone: aws.String(d.Get("availability_zone").(string)),
	}
	if value, ok := d.GetOk("iops"); ok {
		request.Iops = aws.Int64(int64(value.(int)))
	}
	if value, ok := d.GetOk("size"); ok {
		request.Size = aws.Int64(int64(value.(int)))
	}
	if value, ok := d.GetOk("snapshot_id"); ok {
		request.SnapshotId = aws.String(value.(string))
	}
	if value, ok := d.GetOk("volume_type"); ok {
		request.VolumeType = aws.String(value.(string))
	}

	var t string
	if value, ok := d.GetOk("type"); ok {
		t = value.(string)
		request.VolumeType = aws.String(t)
	}

	iops := d.Get("iops").(int)
	if t != "io1" && iops > 0 {
		fmt.Printf("[WARN] IOPs is only valid for storate type io1 for EBS Volumes")
	} else if t == "io1" {
		request.Iops = aws.Int64(int64(iops))
	}

	fmt.Printf(
		"[DEBUG] Volume create opts: %s", request)

	var err error
	var result *fcu.Volume

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		result, err = conn.VM.CreateVolume(request)

		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		return fmt.Errorf("Error creating volume: %s", err)
	}

	fmt.Println("[DEBUG] Waiting for Volume to become available")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available"},
		Refresh:    volumeStateRefreshFunc(conn, *result.VolumeId),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for Volume (%s) to become available: %s",
			*result.VolumeId, err)
	}

	d.SetId(*result.VolumeId)

	if _, ok := d.GetOk("tags"); ok {
		if err := setTags(conn, d); err != nil {
			return errwrap.Wrapf("Error setting tags for EBS Volume: {{err}}", err)
		}
	}
	fmt.Printf("[DEBUG] Volume ID: %s ", *result.VolumeId)
	//return readVolume(d, *result)
	return resourceVolumeRead(d, meta)
}

func resourceVolumeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	request := &fcu.DescribeVolumesInput{
		VolumeIds: []*string{aws.String(d.Id())},
	}

	var err error
	var response *fcu.DescribeVolumesOutput

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		response, err = conn.VM.DescribeVolumes(request)
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
		return fmt.Errorf("Error reading Outscale volume %s: %s", d.Id(), err)
	}
	fmt.Printf("[DEBUG] Volume Read: #v", *response.Volumes[0])
	return readVolume(d, *response.Volumes[0])
}

func readVolume(d *schema.ResourceData, volume fcu.Volume) error {
	d.SetId(*volume.VolumeId)

	if volume.AvailabilityZone != nil {
		d.Set("availability_zone", *volume.AvailabilityZone)
	}
	if volume.Encrypted != nil {
		d.Set("encrypted", *volume.Encrypted)
	}
	if volume.KmsKeyId != nil {
		d.Set("kms_key_id", *volume.KmsKeyId)
	}
	if volume.Size != nil {
		d.Set("size", *volume.Size)
	}
	if volume.SnapshotId != nil {
		d.Set("snapshot_id", *volume.SnapshotId)
	}
	if volume.Attachments != nil {
		res := make([]map[string]interface{}, len(volume.Attachments))
		for k, g := range volume.Attachments {
			r := make(map[string]interface{})
			if g.DeleteOnTermination != nil {
				r["delete_on_termination"] = *g.DeleteOnTermination
			}
			if g.Device != nil {
				r["device"] = *g.Device
			}
			if g.InstanceId != nil {
				r["instance_id"] = *g.InstanceId
			}
			if g.State != nil {
				r["status"] = *g.State
			}
			if g.VolumeId != nil {
				r["volume_id"] = *g.VolumeId
			}

			res[k] = r

		}

		if len(res) > 0 {
			if err := d.Set("attachment_set", res); err != nil {
				return err
			}
			fmt.Printf("[DEBUG] Attachment Set %s", res)
		}

	}
	if volume.Iops != nil {
		d.Set("Iops", *volume.Iops)
	}
	if volume.State != nil {
		d.Set("status", *volume.State)
	}
	if volume.VolumeId != nil {
		d.Set("volume_id", *volume.VolumeId)
	}
	if volume.VolumeType != nil {
		d.Set("volume_type", *volume.VolumeType)
	}
	// if volume.Tags != nil {
	// 	d.Set("tag_set", tagsToMap(volume.Tags))
	// 	t := tagsToMap(volume.Tags)
	// 	fmt.Printf("DEBUG TAGS: ", t)
	// }
	if volume.VolumeType != nil && *volume.VolumeType == "io1" {
		// Only set the iops attribute if the volume type is io1. Setting otherwise
		// can trigger a refresh/plan loop based on the computed value that is given
		// from AWS, and prevent us from specifying 0 as a valid iops.
		//   See https://github.com/hashicorp/terraform/pull/4146
		if volume.Iops != nil {
			d.Set("iops", *volume.Iops)
		}
	}

	return nil
}

func resourceVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := &fcu.DeleteVolumeInput{
			VolumeId: aws.String(d.Id()),
		}

		var err error
		var response *fcu.DeleteVolumeOutput

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			response, err = conn.VM.DeleteVolume(request)
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})

		if err == nil {
			return nil
		}

		ebsErr, ok := err.(awserr.Error)
		if ebsErr.Code() == "VolumeInUse" {
			return resource.RetryableError(fmt.Errorf("Outscale VolumeInUse - trying again while it detaches"))
		}

		if !ok {
			return resource.NonRetryableError(err)
		}

		return resource.NonRetryableError(err)
	})
}

func getVolumeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Arguments
		"availability_zone": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"iops": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Computed: true,
		},
		"size": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
			Computed: true,
		},
		"snapshot_id": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Computed: true,
		},
		"volume_type": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Computed: true,
		},
		// Attributes
		"attachment_set": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"delete_on_termination": {
						Type:     schema.TypeBool,
						Computed: true,
						Optional: true,
					},
					"device": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},
					"instance_id": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},
					"status": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},
					"volume_id": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},
				},
			},
			Computed: true,
			Optional: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
		},
		// "tag_set": tagsSchemaComputed(),
		"tags": tagsSchema(),
		"volume_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

// volumeStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a the state of a Volume. Returns successfully when volume is available
func volumeStateRefreshFunc(conn *fcu.Client, volumeID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// resp, err := conn.VM.DescribeVolumes(&fcu.DescribeVolumesInput{
		// 	VolumeIds: []*string{aws.String(volumeID)},
		// })

		var err error
		var response *fcu.DescribeVolumesOutput

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {

			response, err = conn.VM.DescribeVolumes(&fcu.DescribeVolumesInput{
				VolumeIds: []*string{aws.String(volumeID)},
			})
			if err != nil {
				if strings.Contains(err.Error(), "RequestLimitExceeded:") {
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return resource.NonRetryableError(err)
		})

		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok {
				// Set this to nil as if we didn't find anything.
				fmt.Printf("Error on Volume State Refresh: message: \"%s\", code:\"%s\"", ec2err.Message(), ec2err.Code())
				response = nil
				return nil, "", err
			} else {
				fmt.Printf("Error on Volume State Refresh: %s", err)
				return nil, "", err
			}
		}

		v := response.Volumes[0]

		fmt.Printf("[DEBUG] Volume #v", v)
		return v, *v.State, nil
	}
}
