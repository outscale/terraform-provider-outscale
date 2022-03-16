package outscale

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/terraform-providers/terraform-provider-outscale/utils"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	minIops     = 100
	maxIops     = 13000
	defaultIops = 150
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
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
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
				Computed: true,
			},
			// Attributes
			"linked_volumes": {
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
	var vType string
	var vSize int32
	var vIops int32
	var snapId string

	if value, ok := d.GetOk("size"); ok {
		request.SetSize(int32(value.(int)))
		vSize = int32(value.(int))
	}
	if value, ok := d.GetOk("volume_type"); ok {
		request.SetVolumeType(value.(string))
		vType = value.(string)
	}
	if value, ok := d.GetOk("iops"); ok {
		request.SetIops(int32(value.(int)))
		vIops = int32(value.(int))
	}
	if value, ok := d.GetOk("snapshot_id"); ok {
		request.SetSnapshotId(value.(string))
		snapId = value.(string)
	}
	if vType != "io1" && vIops > 0 {
		return fmt.Errorf("Error IOPs is only valid for storage type 'io1' for EBS Volumes")
	}
	if vType == "io1" && (vIops < minIops || vIops > maxIops) {
		return fmt.Errorf("Cannot create volume type 'io1' without 'iops'. The number of 'iops' allowed for 'io1' volumes is: min: %d and max: %d", minIops, maxIops)
	}
	if vSize == 0 && snapId == "" {
		return fmt.Errorf("The size of the volume is required if the volume is not created from a snapshot")
	}

	var resp oscgo.CreateVolumeResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		resp, _, err = conn.VolumeApi.CreateVolume(context.Background()).CreateVolumeRequest(request).Execute()
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

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		r, _, err := conn.VolumeApi.ReadVolumes(context.Background()).ReadVolumesRequest(request).Execute()
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

	return readOAPIVolume(d, resp.GetVolumes()[0])
}

func resourceOAPIVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	updateReq, err, isAttrChange := setUpdateVolumeAttrRequest(d)
	if err != nil {
		return err
	}
	if isAttrChange {
		readVolReq := oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{VolumeIds: &[]string{d.Id()}},
		}

		var volResp oscgo.ReadVolumesResponse
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			volResp, _, err = conn.VolumeApi.ReadVolumes(context.Background()).ReadVolumesRequest(readVolReq).Execute()
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

		readVolume := volResp.GetVolumes()[0]
		if readVolume.GetVolumeId() == "" {
			return fmt.Errorf("Volume not found")
		}
		if readVolume.GetState() != "available" {
			vm_id := readVolume.GetLinkedVolumes()[0].GetVmId()

			var vmResp oscgo.ReadVmsResponse
			err := resource.Retry(30*time.Second, func() *resource.RetryError {
				vmResp, _, err = conn.VmApi.ReadVms(context.Background()).ReadVmsRequest(oscgo.ReadVmsRequest{
					Filters: &oscgo.FiltersVm{VmIds: &[]string{vm_id}}}).Execute()
				if err != nil {
					if strings.Contains(err.Error(), "RequestLimitExceeded:") {
						return resource.RetryableError(err)
					}
					return resource.NonRetryableError(err)
				}
				return nil
			})
			if err != nil {
				return err
			}
			readVm := vmResp.GetVms()[0]
			if readVm.GetVmId() == "" {
				return fmt.Errorf("Vm not found")
			}
			if readVm.GetState() == "running" && (updateReq.HasVolumeType() || updateReq.HasSize()) {
				log.Println("[WARNING] VM WILL BE STOPPED TO UPDATE VOLUME")
				if err := stopVM(vm_id, conn); err != nil {
					return err
				}
				_, _, err := conn.VolumeApi.UpdateVolume(context.Background()).UpdateVolumeRequest(*updateReq).Execute()
				if err != nil {
					return err
				}
				if err := startVM(vm_id, conn); err != nil {
					return err
				}
			} else {
				_, _, err := conn.VolumeApi.UpdateVolume(context.Background()).UpdateVolumeRequest(*updateReq).Execute()
				if err != nil {
					return err
				}
			}
		} else {
			_, _, err := conn.VolumeApi.UpdateVolume(context.Background()).UpdateVolumeRequest(*updateReq).Execute()
			if err != nil {
				return err
			}
		}

		log.Println("[DEBUG] Waiting for Volume refresh")
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"creating"},
			Target:     []string{"available", "in-use"},
			Refresh:    volumeOAPIStateRefreshFunc(conn, readVolume.GetVolumeId()),
			Timeout:    5 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for Volume (%s) to become available: %s", readVolume.GetVolumeId(), err)
		}
	}
	d.Partial(true)
	if err := setOSCAPITags(conn, d); err != nil {
		return err
	}

	d.SetPartial("tags")
	d.Partial(false)

	return resourceOAPIVolumeRead(d, meta)
}

func resourceOAPIVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*OutscaleClient).OSCAPI

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := oscgo.DeleteVolumeRequest{
			VolumeId: d.Id(),
		}
		_, _, err := conn.VolumeApi.DeleteVolume(context.Background()).DeleteVolumeRequest(request).Execute()
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
		resp, _, err := conn.VolumeApi.ReadVolumes(context.Background()).ReadVolumesRequest(oscgo.ReadVolumesRequest{
			Filters: &oscgo.FiltersVolume{
				VolumeIds: &[]string{volumeID},
			},
		}).Execute()

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

	if err := d.Set("subregion_name", volume.GetSubregionName()); err != nil {
		return err
	}

	//Commented until backend issues is resolved.
	if err := d.Set("size", volume.Size); err != nil {
		return err
	}
	if err := d.Set("snapshot_id", volume.GetSnapshotId()); err != nil {
		return err
	}

	if volume.GetVolumeType() != "" {
		if err := d.Set("volume_type", volume.GetVolumeType()); err != nil {
			return err
		}
	} else if vType, ok := d.GetOk("volume_type"); ok {
		volume.SetVolumeType(vType.(string))
	} else {
		if err := d.Set("volume_type", ""); err != nil {
			return err
		}
	}
	if volume.GetVolumeType() != "standard" {
		if err := d.Set("iops", volume.GetIops()); err != nil {
			return err
		}
	} else {
		if err := d.Set("iops", defaultIops); err != nil {
			return err
		}
	}
	if err := d.Set("state", volume.GetState()); err != nil {
		return err
	}
	if err := d.Set("volume_id", volume.GetVolumeId()); err != nil {
		return err
	}

	if volume.LinkedVolumes != nil {
		res := make([]map[string]interface{}, len(volume.GetLinkedVolumes()))
		for k, g := range volume.GetLinkedVolumes() {
			r := make(map[string]interface{})
			r["delete_on_vm_deletion"] = g.GetDeleteOnVmDeletion()
			if g.GetDeviceName() != "" {
				r["device_name"] = g.GetDeviceName()
			}
			if g.GetVmId() != "" {
				r["vm_id"] = g.GetVmId()
			}
			if g.GetState() != "" {
				r["state"] = g.GetState()
			}
			if g.GetVolumeId() != "" {
				r["volume_id"] = g.GetVolumeId()
			}

			res[k] = r

		}

		if err := d.Set("linked_volumes", res); err != nil {
			return err
		}
	} else {
		if err := d.Set("linked_volumes", []map[string]interface{}{
			{
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
	if volume.GetTags() != nil {
		if err := d.Set("tags", tagsOSCAPIToMap(volume.GetTags())); err != nil {
			return err
		}
	} else {
		if err := d.Set("tags", []map[string]string{
			{
				"key":   "",
				"value": "",
			},
		}); err != nil {
			return err
		}
	}

	return nil
}

func setUpdateVolumeAttrRequest(d *schema.ResourceData) (*oscgo.UpdateVolumeRequest, error, bool) {
	isAttrChange := false
	request := oscgo.UpdateVolumeRequest{
		VolumeId: d.Get("volume_id").(string),
	}
	if d.HasChange("iops") && !d.IsNewResource() {
		_, newIops := d.GetChange("iops")
		oldType, newType := d.GetChange("volume_type")

		if newIops.(int) > 0 && newType.(string) != "io1" {
			return nil, fmt.Errorf("Update IOPS is only valid for volume type 'io1'. Current volume type is '%s'", oldType.(string)), false
		}
		if newIops.(int) < minIops || newIops.(int) > maxIops {
			return nil, fmt.Errorf("The number of IOPS allowed for io1 volumes is: min: %d and max: %d", minIops, maxIops), false
		}
		request.SetIops(int32(newIops.(int)))
		isAttrChange = true
	}

	if d.HasChange("volume_type") && !d.IsNewResource() {
		_, newType := d.GetChange("volume_type")

		if newType.(string) == "io1" && !d.HasChange("iops") {
			return nil, fmt.Errorf("To update volume type 'io1', you must also specify the 'iops' parameter"), false
		}
		request.SetVolumeType(newType.(string))
		isAttrChange = true
	}

	if d.HasChange("size") && !d.IsNewResource() {
		oldSize, newSize := d.GetChange("size")
		if newSize.(int) < oldSize.(int) {
			return nil, fmt.Errorf("The new size of the volume value must be equal to or greater than the current size of the volume"), false
		}
		request.SetSize(int32(newSize.(int)))
		isAttrChange = true
	}
	return &request, nil, isAttrChange
}
