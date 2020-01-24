package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/antihax/optional"
	oscgo "github.com/marinsalinas/osc-sdk-go"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOutscaleOAPIVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceOAPIVolumeCreate,
		Read:   resourceOAPIVolumeRead,
		Update: resourceOAPIVolumeUpdate,
		Delete: resourceOAPIVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// Arguments
			"subregion_name": {
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
			"linked_volumes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_on_vm_termination": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"device": {
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
			"tags": tagsListOAPISchema(),
			"volume_id": {
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

func resourceOAPIVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	request := oscgo.CreateVolumeRequest{
		SubregionName: d.Get("subregion_name").(string),
	}
	if value, ok := d.GetOk("size"); ok {
		request.SetSize(int64(value.(int)))
	}
	if value, ok := d.GetOk("snapshot_id"); ok {
		request.SetSnapshotId(value.(string))
	}

	var t string
	if value, ok := d.GetOk("volume_type"); ok {
		request.SetVolumeType(value.(string))
		t = value.(string)
	}

	iops := d.Get("iops").(int)
	if t != "io1" && iops > 0 {
		log.Printf("[WARN] IOPs is only valid for storate type io1 for EBS Volumes")
	} else if t == "io1" {
		request.SetIops(int64(iops))
	}

	var resp oscgo.CreateVolumeResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.VolumeApi.CreateVolume(context.Background(), &oscgo.CreateVolumeOpts{CreateVolumeRequest: optional.NewInterface(request)})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating Outscale BSU volume: %s", utils.GetErrorResponse(err))
	}

	log.Println("[DEBUG] Waiting for Volume to become available")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available"},
		Refresh:    volumeOAPIStateRefreshFunc(conn, resp.Volume.GetVolumeId()),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Volume (%s) to become available: %s", resp.Volume.GetVolumeId(), err)
	}

	d.SetId(resp.Volume.GetVolumeId())

	if d.IsNewResource() {
		if err := setOSCAPITags(conn, d); err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	return resourceOAPIVolumeRead(d, meta)
}

func resourceOAPIVolumeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	request := oscgo.ReadVolumesRequest{
		Filters: &oscgo.FiltersVolume{VolumeIds: &[]string{d.Id()}},
	}

	var resp oscgo.ReadVolumesResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		r, _, err := conn.VolumeApi.ReadVolumes(context.Background(), &oscgo.ReadVolumesOpts{ReadVolumesRequest: optional.NewInterface(request)})
		if err != nil {
			if strings.Contains(err.Error(), "RequestLimitExceeded:") {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		resp = r
		return nil
	})

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "InvalidVolume.NotFound") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Outscale volume %s: %s", d.Id(), err)
	}
	d.Set("request_id", resp.ResponseContext.GetRequestId())
	return readOAPIVolume(d, resp.GetVolumes()[0])
}

func resourceOAPIVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	d.Partial(true)

	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available"},
		Refresh:    volumeOAPIStateRefreshFunc(conn, d.Id()),
		Timeout:    5 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Volume (%s) to update: %s", d.Id(), err)
	}

	d.Partial(false)
	return resourceOAPIVolumeRead(d, meta)
}

func resourceOAPIVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := oscgo.DeleteVolumeRequest{
			VolumeId: d.Id(),
		}
		_, _, err := conn.VolumeApi.DeleteVolume(context.Background(), &oscgo.DeleteVolumeOpts{DeleteVolumeRequest: optional.NewInterface(request)})
		if err == nil {
			return nil
		}

		if strings.Contains(fmt.Sprint(err), "VolumeInUse") {
			return resource.RetryableError(fmt.Errorf("Outscale VolumeInUse - trying again while it detaches"))
		}
		fmt.Println(err)

		return resource.NonRetryableError(err)
	})

}

func volumeOAPIStateRefreshFunc(conn *oscgo.APIClient, volumeID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, _, err := conn.VolumeApi.ReadVolumes(context.Background(), &oscgo.ReadVolumesOpts{ReadVolumesRequest: optional.NewInterface(oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{
				VolumeIds: &[]string{volumeID},
			},
		})})

		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok {
				log.Printf("Error on Volume State Refresh: message: \"%s\", code:\"%s\"", ec2err.Message(), ec2err.Code())
				//resp = nil
				return nil, "", err
			}
			log.Printf("Error on Volume State Refresh: %s", err)
			return nil, "", err
		}

		v := resp.GetVolumes()[0]
		return v, v.GetState(), nil
	}
}

func readOAPIVolume(d *schema.ResourceData, volume oscgo.Volume) error {
	d.SetId(volume.GetVolumeId())

	d.Set("subregion_name", volume.GetSubregionName())

	//Commented until backend issues is resolved.
	d.Set("size", volume.Size)
	d.Set("snapshot_id", volume.GetSnapshotId())

	if volume.GetVolumeType() != "" {
		d.Set("volume_type", volume.GetVolumeType())
	} else if vType, ok := d.GetOk("volume_type"); ok {
		volume.SetVolumeType(vType.(string))
	} else {
		d.Set("volume_type", "")
	}

	d.Set("iops", volume.GetIops())
	d.Set("state", volume.GetState())
	d.Set("volume_id", volume.GetVolumeId())

	if volume.GetLinkedVolumes() != nil {
		res := make([]map[string]interface{}, len(volume.GetLinkedVolumes()))
		for k, g := range volume.GetLinkedVolumes() {
			r := make(map[string]interface{})
			if g.DeleteOnVmDeletion != nil {
				r["delete_on_vm_termination"] = g.GetDeleteOnVmDeletion()
			}
			if g.GetDeviceName() != "" {
				r["device"] = g.DeviceName
			}
			if g.GetVmId() != "" {
				r["vm_id"] = g.VmId
			}
			if g.GetState() != "" {
				r["state"] = g.State
			}
			if g.GetVolumeId() != "" {
				r["volume_id"] = g.VolumeId
			}

			res[k] = r

		}

		if err := d.Set("linked_volumes", res); err != nil {
			return err
		}
	} else {
		if err := d.Set("linked_volumes", []map[string]interface{}{
			map[string]interface{}{
				"delete_on_vm_termination": false,
				"device":                   "none",
				"vm_id":                    "none",
				"state":                    "none",
				"volume_id":                "none",
			},
		}); err != nil {
			return err
		}
	}
	if volume.GetTags() != nil {
		if err := d.Set("tags", tagsOSCAPIToMap(volume.GetTags())); err != nil {
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
