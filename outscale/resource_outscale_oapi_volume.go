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

func resourceOAPIOutscaleVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIVolumeCreate,
		Read:   resourceOAPIVolumeRead,
		Delete: resourceOAPIVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// Arguments
			"sub_region_name": {
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
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			// Attributes
			"linked_volume": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_on_vm_deletion": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
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
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
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

func resourceOAPIVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).FCU

	request := &fcu.CreateVolumeInput{
		AvailabilityZone: aws.String(d.Get("sub_region_name").(string)),
	}
	if value, ok := d.GetOk("size"); ok {
		request.Size = aws.Int64(int64(value.(int)))
	}
	if value, ok := d.GetOk("snapshot_id"); ok {
		request.SnapshotId = aws.String(value.(string))
	}

	var t string
	if value, ok := d.GetOk("type"); ok {
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

	log.Printf(
		"[DEBUG] Outscale Volume create opts: %s", request)

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
		Refresh:    volumeOAPIStateRefreshFunc(conn, *result.VolumeId),
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
		d.SetPartial("tags")
	}

	return readOAPIVolume(d, result)
}

func resourceOAPIVolumeRead(d *schema.ResourceData, meta interface{}) error {

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
		return resource.NonRetryableError(err)
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidVolume.NotFound") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Outscale volume %s: %s", d.Id(), err)
	}

	return readOAPIVolume(d, response.Volumes[0])
}

func resourceOAPIVolumeDelete(d *schema.ResourceData, meta interface{}) error {
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

func volumeOAPIStateRefreshFunc(conn *fcu.Client, volumeID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := conn.VM.DescribeVolumes(&fcu.DescribeVolumesInput{
			VolumeIds: []*string{aws.String(volumeID)},
		})

		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok {
				log.Printf("Error on Volume State Refresh: message: \"%s\", code:\"%s\"", ec2err.Message(), ec2err.Code())
				resp = nil
				return nil, "", err
			} else {
				log.Printf("Error on Volume State Refresh: %s", err)
				return nil, "", err
			}
		}

		v := resp.Volumes[0]
		return v, *v.State, nil
	}
}

func readOAPIVolume(d *schema.ResourceData, volume *fcu.Volume) error {
	d.SetId(*volume.VolumeId)

	d.Set("sub_region_name", *volume.AvailabilityZone)
	if volume.Size != nil {
		d.Set("size", *volume.Size)
	}
	if volume.SnapshotId != nil {
		d.Set("snapshot_id", *volume.SnapshotId)
	}
	if volume.VolumeType != nil {
		d.Set("type", *volume.VolumeType)
	}

	if volume.VolumeType != nil && *volume.VolumeType == "io1" {
		if volume.Iops != nil {
			d.Set("iops", *volume.Iops)
		}
	}
	if volume.State != nil {
		d.Set("state", *volume.State)
	}
	if volume.VolumeId != nil {
		d.Set("volume_id", *volume.VolumeId)
	}
	if volume.VolumeType != nil {
		d.Set("type", *volume.VolumeType)
	}
	if volume.Attachments != nil {
		res := make([]map[string]interface{}, len(volume.Attachments))
		for k, g := range volume.Attachments {
			r := make(map[string]interface{})
			if g.DeleteOnTermination != nil {
				r["delete_on_vm_deletion"] = *g.DeleteOnTermination
			}
			if g.Device != nil {
				r["device_name"] = *g.Device
			}
			if g.InstanceId != nil {
				r["vm_id"] = *g.InstanceId
			}
			if g.State != nil {
				r["state"] = *g.State
			}
			if g.VolumeId != nil {
				r["volume_id"] = *g.VolumeId
			}

			res[k] = r

		}

		if err := d.Set("linked_volume", res); err != nil {
			return err
		}
	} else {
		if err := d.Set("linked_volume", []map[string]interface{}{
			map[string]interface{}{
				"delete_on_vm_deletion": false,
				"device_name":           "none",
				"vm_id":                 "none",
				"state":                 "none",
				"volume_id":             "none",
			},
		}); err != nil {
			return err
		}
	}
	if volume.Tags != nil {
		if err := d.Set("tags", tagsToMap(volume.Tags)); err != nil {
			return err
		}
	} else {
		if err := d.Set("tags", []map[string]string{
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
