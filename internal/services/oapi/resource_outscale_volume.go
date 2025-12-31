package oapi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	oscgo "github.com/outscale/osc-sdk-go/v2"
	"github.com/outscale/terraform-provider-outscale/internal/client"
	"github.com/outscale/terraform-provider-outscale/internal/utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceOutscaleVolume() *schema.Resource {
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
			"creation_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"iops": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					iopsVal := val.(int32)
					if iopsVal < MinIops || iopsVal > MaxIops {
						errs = append(errs, fmt.Errorf("%q must be between %d and %d inclusive, got: %d", key, MinIops, MaxIops, iopsVal))
					}
					return
				},
			},
			"size": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					vSize := val.(int32)
					if vSize < 1 || vSize > MaxSize {
						errs = append(errs, fmt.Errorf("%q must be between 1 and %d gibibytes inclusive, got: %d", key, MaxSize, vSize))
					}
					return
				},
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"termination_snapshot_name": {
				Type:     schema.TypeString,
				Optional: true,
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
			"tags": TagsSchemaSDK(),
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
	conn := meta.(*client.OutscaleClient).OSCAPI

	request := oscgo.CreateVolumeRequest{
		SubregionName: d.Get("subregion_name").(string),
	}
	vSize := int32(d.Get("size").(int))
	snapId := d.Get("snapshot_id").(string)
	vType := d.Get("volume_type").(string)

	if snapId == "" && vSize == 0 {
		return fmt.Errorf("Error: 'size' parameter is required if the volume is not created from a snapshot (SnapshotId unspecified)")
	}
	if value, ok := d.GetOk("iops"); ok {
		if vType != "io1" {
			return fmt.Errorf("Error: %s", VolumeIOPSError)
		}
		request.SetIops(int32(value.(int)))
	}
	if snapId != "" {
		request.SetSnapshotId(snapId)
	}
	if vType != "" {
		request.SetVolumeType(vType)
	}
	request.SetSize(vSize)

	var resp oscgo.CreateVolumeResponse
	var err error

	err = retry.Retry(5*time.Minute, func() *retry.RetryError {
		rp, httpResp, err := conn.VolumeApi.CreateVolume(context.Background()).CreateVolumeRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = rp
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error creating Outscale BSU volume: %s", utils.GetErrorResponse(err))
	}

	log.Println("[DEBUG] Waiting for Volume to become available")

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available"},
		Refresh:    volumeOAPIStateRefreshFunc(conn, resp.Volume.GetVolumeId()),
		Timeout:    5 * time.Minute,
		Delay:      4 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Volume (%s) to become available: %s", resp.Volume.GetVolumeId(), err)
	}

	d.SetId(resp.Volume.GetVolumeId())

	if d.IsNewResource() {
		if err := updateOAPITagsSDK(conn, d); err != nil {
			return err
		}
	}
	return resourceOAPIVolumeRead(d, meta)
}

func resourceOAPIVolumeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	request := oscgo.ReadVolumesRequest{
		Filters: &oscgo.FiltersVolume{VolumeIds: &[]string{d.Id()}},
	}

	var resp oscgo.ReadVolumesResponse
	var statusCode int
	err := retry.Retry(5*time.Minute, func() *retry.RetryError {
		r, httpResp, err := conn.VolumeApi.ReadVolumes(context.Background()).ReadVolumesRequest(request).Execute()
		if err != nil {
			return utils.CheckThrottling(httpResp, err)
		}
		resp = r
		statusCode = httpResp.StatusCode
		return nil
	})
	if err != nil {
		if statusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Outscale volume %s: %s", d.Id(), err)
	}

	if utils.IsResponseEmpty(len(resp.GetVolumes()), "Snapshot", d.Id()) {
		d.SetId("")
		return nil
	}
	return readOAPIVolume(d, resp.GetVolumes()[0])
}

func resourceOAPIVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	if err := updateOAPITagsSDK(conn, d); err != nil {
		return err
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available", "in-use"},
		Refresh:    volumeOAPIStateRefreshFunc(conn, d.Id()),
		Timeout:    5 * time.Minute,
		Delay:      2 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for Volume (%s) to update: %s", d.Id(), err)
	}
	return resourceOAPIVolumeRead(d, meta)
}

func resourceOAPIVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.OutscaleClient).OSCAPI

	if snpName, ok := d.GetOk("termination_snapshot_name"); ok {
		volId := d.Get("volume_id").(string)
		description := "Created before volume deletion"
		resp := oscgo.CreateSnapshotResponse{}
		request := oscgo.CreateSnapshotRequest{
			Description: &description,
			VolumeId:    &volId,
		}
		err := retry.Retry(5*time.Minute, func() *retry.RetryError {
			var err error
			r, httpResp, err := conn.SnapshotApi.CreateSnapshot(context.Background()).CreateSnapshotRequest(request).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = r
			return nil
		})
		if err != nil {
			return err
		}
		snapTagsReq := oscgo.CreateTagsRequest{
			ResourceIds: []string{resp.Snapshot.GetSnapshotId()},
			Tags: []oscgo.ResourceTag{
				{
					Key:   "Name",
					Value: snpName.(string),
				},
			},
		}
		err = retry.Retry(60*time.Second, func() *retry.RetryError {
			_, httpResp, err := conn.TagApi.CreateTags(context.Background()).CreateTagsRequest(snapTagsReq).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return retry.Retry(5*time.Minute, func() *retry.RetryError {
		request := oscgo.DeleteVolumeRequest{
			VolumeId: d.Id(),
		}
		_, httpResp, err := conn.VolumeApi.DeleteVolume(context.Background()).DeleteVolumeRequest(request).Execute()
		if err != nil {
			if strings.Contains(fmt.Sprint(err), "VolumeInUse") {
				return retry.RetryableError(fmt.Errorf("Outscale VolumeInUse - trying again while it detaches"))
			}
			return utils.CheckThrottling(httpResp, err)
		}
		return nil
	})
}

func volumeOAPIStateRefreshFunc(conn *oscgo.APIClient, volumeID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp oscgo.ReadVolumesResponse
		var err error
		err = retry.Retry(3*time.Minute, func() *retry.RetryError {
			rp, httpResp, err := conn.VolumeApi.ReadVolumes(context.Background()).ReadVolumesRequest(oscgo.ReadVolumesRequest{
				Filters: &oscgo.FiltersVolume{
					VolumeIds: &[]string{volumeID},
				},
			}).Execute()
			if err != nil {
				return utils.CheckThrottling(httpResp, err)
			}
			resp = rp
			return nil
		})
		if err != nil {
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

	// Commented until backend issues is resolved.
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
	if err := d.Set("creation_date", volume.GetCreationDate()); err != nil {
		return err
	}
	if snapName, ok := d.GetOk("termination_snapshot_name"); ok {
		if err := d.Set("termination_snapshot_name", snapName.(string)); err != nil {
			return err
		}
	}
	if err := d.Set("iops", getIops(volume.GetVolumeType(), volume.GetIops())); err != nil {
		return err
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
		if err := d.Set("tags", FlattenOAPITagsSDK(volume.GetTags())); err != nil {
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
