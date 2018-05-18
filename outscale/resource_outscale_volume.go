package outscale

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
		Schema: map[string]*schema.Schema{
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
			"request_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Attributes
			"attachment_set": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_on_termination": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"volume_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag_set": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"tag": tagsSchema(),
			"volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	request := &fcu.CreateVolumeInput{
		AvailabilityZone: aws.String(d.Get("availability_zone").(string)),
	}
	if value, ok := d.GetOk("size"); ok {
		request.Size = aws.Int64(int64(value.(int)))
	}
	if value, ok := d.GetOk("snapshot_id"); ok {
		request.SnapshotId = aws.String(value.(string))
	}

	var t string
	if value, ok := d.GetOk("volume_type"); ok {
		t = value.(string)
		request.VolumeType = aws.String(t)
	}

	iops := d.Get("iops").(int)
	if t != "io1" && iops > 0 {
		log.Printf("[WARN] IOPs is only valid for storate type io1 for EBS Volumes")
	} else if t == "io1" {
		request.Iops = aws.Int64(int64(iops))
	}

	tagsSpec := make([]*fcu.TagSpecification, 0)

	if v, ok := d.GetOk("tag"); ok {
		tag := tagsFromMap(v.(map[string]interface{}))

		spec := &fcu.TagSpecification{
			ResourceType: aws.String("volume"),
			Tags:         tag,
		}

		tagsSpec = append(tagsSpec, spec)
	}

	if len(tagsSpec) > 0 {
		request.TagSpecifications = tagsSpec
	}

	var result *fcu.Volume
	var err error

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
		return fmt.Errorf("Error creating Outscale VM volume: %s", err)
	}

	log.Println("[DEBUG] Waiting for Volume to become available")

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

	if d.IsNewResource() {
		if err := setTags(conn, d); err != nil {
			return err
		}
		d.SetPartial("tag_set")
	}

	return resourceVolumeRead(d, meta)
}

func resourceVolumeRead(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*OutscaleClient).FCU

	request := &fcu.DescribeVolumesInput{
		VolumeIds: []*string{aws.String(d.Id())},
	}

	var response *fcu.DescribeVolumesOutput
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.VM.DescribeVolumes(request)
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidVolume.NotFound") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Outscale volume %s: %s", d.Id(), err)
	}

	d.Set("request_id", response.RequestId)

	return readVolume(d, response.Volumes[0])
}

func resourceVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := &fcu.DeleteVolumeInput{
			VolumeId: aws.String(d.Id()),
		}
		_, err := conn.VM.DeleteVolume(request)
		if err == nil {
			return nil
		}

		if strings.Contains(fmt.Sprint(err), "VolumeInUse") {
			return resource.RetryableError(fmt.Errorf("Outscale VolumeInUse - trying again while it detaches"))
		}

		return resource.NonRetryableError(err)
	})

}

func volumeStateRefreshFunc(conn *fcu.Client, volumeID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := conn.VM.DescribeVolumes(&fcu.DescribeVolumesInput{
			VolumeIds: []*string{aws.String(volumeID)},
		})

		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok {
				log.Printf("Error on Volume State Refresh: message: \"%s\", code:\"%s\"", ec2err.Message(), ec2err.Code())
				resp = nil
				return nil, "", err
			}
			log.Printf("Error on Volume State Refresh: %s", err)
			return nil, "", err
		}

		v := resp.Volumes[0]
		return v, *v.State, nil
	}
}

func readVolume(d *schema.ResourceData, volume *fcu.Volume) error {
	d.SetId(*volume.VolumeId)

	d.Set("availability_zone", *volume.AvailabilityZone)
	if volume.Size != nil {
		d.Set("size", *volume.Size)
	}
	if volume.SnapshotId != nil {
		d.Set("snapshot_id", *volume.SnapshotId)
	}
	if volume.VolumeType != nil {
		d.Set("volume_type", *volume.VolumeType)
	}

	if volume.VolumeType != nil && *volume.VolumeType == "io1" {
		if volume.Iops != nil {
			d.Set("iops", *volume.Iops)
		}
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

		if err := d.Set("attachment_set", res); err != nil {
			return err
		}
	} else {
		if err := d.Set("attachment_set", []map[string]interface{}{
			map[string]interface{}{
				"delete_on_termination": false,
				"device":                "none",
				"instance_id":           "none",
				"status":                "none",
				"volume_id":             "none",
			},
		}); err != nil {
			return err
		}
	}
	if volume.Tags != nil {
		if err := d.Set("tag_set", tagsToMap(volume.Tags)); err != nil {
			return err
		}
	} else {
		if err := d.Set("tag_set", []map[string]string{
			map[string]string{
				"key":   "",
				"value": "",
			},
		}); err != nil {
			return err
		}
	}

	return nil
}
